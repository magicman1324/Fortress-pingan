package ws

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pingan/bastion/internal/model"
	sshcli "github.com/pingan/bastion/internal/ssh"

	"github.com/pingan/bastion/internal/repository"
	"github.com/pingan/bastion/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type wsMessage struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	Cols int    `json:"cols,omitempty"`
	Rows int    `json:"rows,omitempty"`
}

type Handler struct {
	assetSvc   *service.AssetService
	sessionSvc *service.SessionService
	auditRepo  *repository.AuditRepo
	jwtSecret  string
}

func NewHandler(
	assetSvc *service.AssetService,
	sessionSvc *service.SessionService,
	auditRepo *repository.AuditRepo,
	jwtSecret string,
) *Handler {
	return &Handler{
		assetSvc:   assetSvc,
		sessionSvc: sessionSvc,
		auditRepo:  auditRepo,
		jwtSecret:  jwtSecret,
	}
}

func (h *Handler) HandleTerminal(c *gin.Context, hub *Hub) {
	assetID, err := strconv.ParseInt(c.Param("assetId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}

	userID := c.GetInt64("userID")
	if userID == 0 {
		userID = 1
	}

	asset, err := h.assetSvc.GetByIDWithCredential(assetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found: " + err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WS upgrade failed: %v", err)
		return
	}

	sshClient, err := sshcli.Dial(asset.Host, asset.Port, asset.Username, asset.AuthType, asset.Credential)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("SSH connection failed: "+err.Error()+"\r\n"))
		conn.Close()
		return
	}

	pty, err := sshcli.NewPTYSession(sshClient, 120, 40)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("PTY session failed: "+err.Error()+"\r\n"))
		conn.Close()
		sshClient.Close()
		return
	}

	session := &model.Session{
		UserID:    userID,
		AssetID:   assetID,
		Status:    model.SessionStatusActive,
		ClientIP:  c.ClientIP(),
		StartedAt: time.Now(),
	}
	if err := h.sessionSvc.Create(session); err != nil {
		log.Printf("Failed to create session record: %v", err)
	}

	// Audit: use io.Pipe + bufio.Scanner to extract command lines
	auditReader, auditWriter := io.Pipe()
	auditCh := make(chan *model.AuditLog, 256)

	// Tee: SSH stdout → both WebSocket and audit pipe
	teeReader := io.TeeReader(pty.Stdout, auditWriter)

	// Audit scanner goroutine
	go func() {
		scanner := bufio.NewScanner(auditReader)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 {
				auditCh <- &model.AuditLog{
					SessionID:  session.ID,
					UserID:     userID,
					AssetID:    assetID,
					Command:    line,
					ExecutedAt: time.Now(),
				}
			}
		}
	}()

	activeSess := &ActiveSession{
		ID:        session.ID,
		UserID:    userID,
		AssetID:   assetID,
		Conn:      conn,
		SSHConn:   sshClient.Client,
		SSHSess:   pty.Session,
		StartedAt: time.Now(),
		closeCh:   make(chan struct{}),
	}
	hub.Register(activeSess)

	// Welcome banner
	welcome := "\x1b[1;32m=== Bastion Host ===\x1b[0m\r\n" +
		"\x1b[36mTarget: " + asset.Name + " (" + asset.Host + ":" + strconv.Itoa(asset.Port) + ")\x1b[0m\r\n" +
		"\x1b[36mUser: " + asset.Username + "\x1b[0m\r\n" +
		"\x1b[33mAll commands are audited.\x1b[0m\r\n\r\n"
	conn.WriteMessage(websocket.TextMessage, []byte(welcome))

	// Goroutine 1: Read from TeeReader (SSH stdout) → write to WebSocket
	go func() {
		defer func() {
			pty.Close()
			sshClient.Close()
			auditWriter.Close()
			hub.Unregister(session.ID)
			h.sessionSvc.Close(session.ID)
			conn.Close()
		}()

		buf := make([]byte, 4096)
		for {
			select {
			case <-activeSess.closeCh:
				return
			default:
			}
			n, err := teeReader.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("SSH read error: %v", err)
				}
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}
	}()

	// Goroutine 2: Read from WebSocket → write to SSH stdin
	go func() {
		for {
			select {
			case <-activeSess.closeCh:
				return
			default:
			}
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var ctrl wsMessage
			if json.Unmarshal(msg, &ctrl) == nil && ctrl.Type != "" {
				switch ctrl.Type {
				case "resize":
					pty.Resize(ctrl.Cols, ctrl.Rows)
				case "ping":
					conn.WriteJSON(wsMessage{Type: "pong"})
				case "input":
					pty.Stdin.Write([]byte(ctrl.Data))
				}
				continue
			}
			pty.Stdin.Write(msg)
		}
	}()

	// Goroutine 3: Audit batch inserter
	go func() {
		batch := make([]*model.AuditLog, 0, 50)
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		flush := func() {
			if len(batch) == 0 {
				return
			}
			if err := h.auditRepo.BatchInsert(batch); err != nil {
				log.Printf("Audit batch insert error: %v", err)
			}
			batch = batch[:0]
		}

		for {
			select {
			case <-activeSess.closeCh:
				flush()
				return
			case logEntry := <-auditCh:
				batch = append(batch, logEntry)
				if len(batch) >= 50 {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()

	// Goroutine 4: heartbeat
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-activeSess.closeCh:
				return
			case <-ticker.C:
				conn.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()
}
