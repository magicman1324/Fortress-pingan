package ssh

import (
	"fmt"
	"io"

	gossh "golang.org/x/crypto/ssh"
)

type PTYSession struct {
	Session  *gossh.Session
	Stdin    io.WriteCloser
	Stdout   io.Reader
	cols     int
	rows     int
	client   *gossh.Client
}

func NewPTYSession(client *Client, cols, rows int) (*PTYSession, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	modes := gossh.TerminalModes{
		gossh.ECHO:          1,
		gossh.TTY_OP_ISPEED: 14400,
		gossh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to request pty: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Combine stdout and stderr
	combinedReader := io.MultiReader(stdout, stderr)

	if err := session.Shell(); err != nil {
		session.Close()
		return nil, fmt.Errorf("failed to start shell: %w", err)
	}

	return &PTYSession{
		Session: session,
		Stdin:   stdin,
		Stdout:  combinedReader,
		cols:    cols,
		rows:    rows,
		client:  client.Client,
	}, nil
}

func (p *PTYSession) Resize(cols, rows int) error {
	if cols == p.cols && rows == p.rows {
		return nil
	}
	p.cols = cols
	p.rows = rows
	return p.Session.WindowChange(rows, cols)
}

func (p *PTYSession) Close() {
	p.Session.Close()
}
