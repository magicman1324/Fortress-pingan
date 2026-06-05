package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pingan/bastion/internal/model"
)

type SessionRepo struct {
	db *sqlx.DB
}

func NewSessionRepo(db *sqlx.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(s *model.Session) error {
	result, err := r.db.Exec(
		"INSERT INTO sessions (user_id, asset_id, status, client_ip, started_at) VALUES (?, ?, ?, ?, ?)",
		s.UserID, s.AssetID, s.Status, s.ClientIP, s.StartedAt,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	s.ID = id
	return nil
}

func (r *SessionRepo) Close(id int64) error {
	_, err := r.db.Exec(
		"UPDATE sessions SET status=?, ended_at=? WHERE id=?",
		model.SessionStatusClosed, time.Now(), id,
	)
	return err
}

func (r *SessionRepo) FindAll(limit, offset int) ([]model.Session, error) {
	var sessions []model.Session
	query := `
		SELECT s.id, s.user_id, s.asset_id, s.status, s.client_ip, s.started_at, s.ended_at,
		       u.username, a.name as asset_name
		FROM sessions s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN assets a ON s.asset_id = a.id
		ORDER BY s.started_at DESC
		LIMIT ? OFFSET ?`
	err := r.db.Select(&sessions, query, limit, offset)
	return sessions, err
}

func (r *SessionRepo) FindActive() ([]model.Session, error) {
	var sessions []model.Session
	query := `
		SELECT s.id, s.user_id, s.asset_id, s.status, s.client_ip, s.started_at, s.ended_at,
		       u.username, a.name as asset_name
		FROM sessions s
		LEFT JOIN users u ON s.user_id = u.id
		LEFT JOIN assets a ON s.asset_id = a.id
		WHERE s.status = ?
		ORDER BY s.started_at DESC`
	err := r.db.Select(&sessions, query, model.SessionStatusActive)
	return sessions, err
}

func (r *SessionRepo) Count() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM sessions")
	return count, err
}

func (r *SessionRepo) CountActive() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM sessions WHERE status = ?", model.SessionStatusActive)
	return count, err
}
