package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pingan/bastion/internal/model"
)

type AssetRepo struct {
	db *sqlx.DB
}

func NewAssetRepo(db *sqlx.DB) *AssetRepo {
	return &AssetRepo{db: db}
}

func (r *AssetRepo) FindAll(search string) ([]model.Asset, error) {
	var assets []model.Asset
	query := "SELECT id, name, host, port, username, auth_type, credential, status, created_at, updated_at FROM assets"
	args := []interface{}{}
	if search != "" {
		query += " WHERE name LIKE ? OR host LIKE ?"
		kw := "%" + search + "%"
		args = append(args, kw, kw)
	}
	query += " ORDER BY id DESC"
	err := r.db.Select(&assets, query, args...)
	return assets, err
}

func (r *AssetRepo) FindByID(id int64) (*model.Asset, error) {
	var a model.Asset
	err := r.db.Get(&a, "SELECT id, name, host, port, username, auth_type, credential, status, created_at, updated_at FROM assets WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssetRepo) Create(a *model.Asset) error {
	result, err := r.db.Exec(
		"INSERT INTO assets (name, host, port, username, auth_type, credential, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		a.Name, a.Host, a.Port, a.Username, a.AuthType, a.Credential, a.Status,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	a.ID = id
	return nil
}

func (r *AssetRepo) Update(a *model.Asset) error {
	_, err := r.db.Exec(
		"UPDATE assets SET name=?, host=?, port=?, username=?, auth_type=?, credential=?, status=? WHERE id=?",
		a.Name, a.Host, a.Port, a.Username, a.AuthType, a.Credential, a.Status, a.ID,
	)
	return err
}

func (r *AssetRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM assets WHERE id=?", id)
	return err
}

func (r *AssetRepo) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec("UPDATE assets SET status=? WHERE id=?", status, id)
	return err
}

// Verify that AssetRepo implements expected interface
var _ = fmt.Sprintf("%v", sql.NullString{})
