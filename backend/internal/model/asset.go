package model

import "time"

type Asset struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Host       string    `db:"host" json:"host"`
	Port       int       `db:"port" json:"port"`
	Username   string    `db:"username" json:"username"`
	AuthType   string    `db:"auth_type" json:"auth_type"`
	Credential string    `db:"credential" json:"-"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

const (
	AuthTypePassword = "password"
	AuthTypeKey      = "key"
)
