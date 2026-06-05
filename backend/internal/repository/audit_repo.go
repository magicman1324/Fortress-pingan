package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pingan/bastion/internal/model"
)

type AuditRepo struct {
	db *sqlx.DB
}

func NewAuditRepo(db *sqlx.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) BatchInsert(logs []*model.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}
	placeholders := make([]string, len(logs))
	args := make([]interface{}, 0, len(logs)*4)
	for i, l := range logs {
		placeholders[i] = "(?, ?, ?, ?, ?)"
		args = append(args, l.SessionID, l.UserID, l.AssetID, l.Command, l.ExecutedAt)
	}
	query := fmt.Sprintf(
		"INSERT INTO audit_logs (session_id, user_id, asset_id, command, executed_at) VALUES %s",
		strings.Join(placeholders, ", "),
	)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AuditRepo) Search(userID, assetID int64, keyword string, limit, offset int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	query := `
		SELECT a.id, a.session_id, a.user_id, a.asset_id, a.command, a.executed_at,
		       u.username, ast.name as asset_name
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		LEFT JOIN assets ast ON a.asset_id = ast.id
		WHERE 1=1`
	args := []interface{}{}

	if userID > 0 {
		query += " AND a.user_id = ?"
		args = append(args, userID)
	}
	if assetID > 0 {
		query += " AND a.asset_id = ?"
		args = append(args, assetID)
	}
	if keyword != "" {
		query += " AND a.command LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	query += " ORDER BY a.executed_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := r.db.Select(&logs, query, args...)
	return logs, err
}

func (r *AuditRepo) Count(userID, assetID int64, keyword string) (int, error) {
	query := "SELECT COUNT(*) FROM audit_logs WHERE 1=1"
	args := []interface{}{}

	if userID > 0 {
		query += " AND user_id = ?"
		args = append(args, userID)
	}
	if assetID > 0 {
		query += " AND asset_id = ?"
		args = append(args, assetID)
	}
	if keyword != "" {
		query += " AND command LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	var count int
	err := r.db.Get(&count, query, args...)
	return count, err
}
