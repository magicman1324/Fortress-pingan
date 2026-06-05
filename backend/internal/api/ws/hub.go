package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	gossh "golang.org/x/crypto/ssh"
)

type ActiveSession struct {
	ID        int64
	UserID    int64
	AssetID   int64
	Conn      *websocket.Conn
	SSHConn   *gossh.Client
	SSHSess   *gossh.Session
	StartedAt time.Time
	closeCh   chan struct{}
}

type Hub struct {
	mu       sync.RWMutex
	sessions map[int64]*ActiveSession
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[int64]*ActiveSession),
	}
}

func (h *Hub) Run() {
	// Background cleanup of stale sessions
}

func (h *Hub) Register(sess *ActiveSession) {
	h.mu.Lock()
	h.sessions[sess.ID] = sess
	h.mu.Unlock()
}

func (h *Hub) Unregister(sessionID int64) {
	h.mu.Lock()
	if s, ok := h.sessions[sessionID]; ok {
		close(s.closeCh)
		delete(h.sessions, sessionID)
	}
	h.mu.Unlock()
}

func (h *Hub) Terminate(sessionID int64) error {
	h.mu.Lock()
	s, ok := h.sessions[sessionID]
	h.mu.Unlock()
	if !ok {
		return nil
	}
	s.Conn.Close()
	s.SSHSess.Close()
	s.SSHConn.Close()
	h.Unregister(sessionID)
	return nil
}

func (h *Hub) ActiveCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.sessions)
}
