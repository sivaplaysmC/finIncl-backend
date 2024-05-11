package models

import (
	"database/sql"
	"log"
)

type UserModel struct {
	db  *sql.DB
	Log *log.Logger
}

type User struct {
	Name, Email string
	SmeID       int
	id          int
}

func NewUserModel(db *sql.DB, logger *log.Logger) *UserModel {
	return &UserModel{
		db:  db,
		Log: logger,
	}
}
