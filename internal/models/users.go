package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = um.DB.Exec(`insert into users (name, email, hashed_password, created) 
			values ($1, $2, $3, now())`, name, email, string(hashedPassword))
	if err != nil {
		var pqError *pq.Error

		if errors.As(err, &pqError) {
			if pqError.Code == "23505" && strings.Contains(pqError.Message, "user_uc_email") {
				return ErrDuplicateEmail
			}
		}
	}

	return nil

}

func (um *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	err := um.DB.QueryRow(`select id, hashed_password 
			from users where email = $1`, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {

	var exists bool

	err := um.DB.QueryRow(`select exists(select 1 from users where id=$1)`, exists).Scan(&exists)
	return exists, err
}
