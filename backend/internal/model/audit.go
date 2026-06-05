package model

import "time"

type AuditLog struct {
	ID         int64     `db:"id" json:"id"`
	SessionID  int64     `db:"session_id" json:"session_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	AssetID    int64     `db:"asset_id" json:"asset_id"`
	Command    string    `db:"command" json:"command"`
	ExecutedAt time.Time `db:"executed_at" json:"executed_at"`

	// Joined fields
	Username  string `db:"username" json:"username,omitempty"`
	AssetName string `db:"asset_name" json:"asset_name,omitempty"`
}
