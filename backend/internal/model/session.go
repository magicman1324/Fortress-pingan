package model

import "time"

type Session struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	AssetID   int64      `db:"asset_id" json:"asset_id"`
	Status    string     `db:"status" json:"status"`
	ClientIP  string     `db:"client_ip" json:"client_ip"`
	StartedAt time.Time  `db:"started_at" json:"started_at"`
	EndedAt   *time.Time `db:"ended_at" json:"ended_at"`

	// Joined fields (not from sessions table)
	Username  string `db:"username" json:"username,omitempty"`
	AssetName string `db:"asset_name" json:"asset_name,omitempty"`
}

const (
	SessionStatusActive = "active"
	SessionStatusClosed = "closed"
)
