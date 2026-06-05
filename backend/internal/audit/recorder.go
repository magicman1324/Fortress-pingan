package audit

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/pingan/bastion/internal/model"
)

type Recorder struct {
	sessionID  int64
	userID     int64
	assetID    int64
	output     chan *model.AuditLog
	done       chan struct{}
	buffer     []*model.AuditLog
}

func NewRecorder(sessionID, userID, assetID int64, output chan *model.AuditLog) *Recorder {
	return &Recorder{
		sessionID: sessionID,
		userID:    userID,
		assetID:   assetID,
		output:    output,
		done:      make(chan struct{}),
	}
}

// Wrap returns an io.Writer that tees data to both the original writer and the audit scanner
func (r *Recorder) Wrap(w io.Writer) io.Writer {
	reader, writer := io.Pipe()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	go func() {
		defer close(r.done)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			r.output <- &model.AuditLog{
				SessionID:  r.sessionID,
				UserID:     r.userID,
				AssetID:    r.assetID,
				Command:    line,
				ExecutedAt: time.Now(),
			}
		}
	}()

	return io.MultiWriter(w, writer)
}

func (r *Recorder) Close() {
	r.Flush()
}

func (r *Recorder) Flush() {
	// The scanner goroutine will close `done` when the pipe is closed
}
