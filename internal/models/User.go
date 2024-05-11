package models

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
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

func (users *UserModel) ValidateUser(username, password string) (int, error) {
	var passHash string
	var id int
	err := users.db.
		QueryRow("SELECT passHash, id from users where username like ?", username).
		Scan(&passHash, &id)
	if err != nil {
		return 0, err
	}

	return id, bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
}

func (users *UserModel) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := users.db.
		QueryRow("SELECT name, email, SMEid from users where id like ?", id).
		Scan(&user.Name, &user.Email, &user.SmeID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (users *UserModel) CreateUser(user *User, passHash []byte) error {
	_, err := users.db.
		Query(
			`INSERT into users (username, SMEid, email, passHash) values (?, ?, ?, ?)`,
			user.Name, user.id, user.Email, passHash,
		)
	return err
}
