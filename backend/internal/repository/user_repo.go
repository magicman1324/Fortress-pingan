package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pingan/bastion/internal/model"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByID(id int64) (*model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u *model.User) error {
	result, err := r.db.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		u.Username, u.PasswordHash, u.Role,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	u.ID = id
	return nil
}
