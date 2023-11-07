package main

import (
	"database/sql"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

func (u *User) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, email, created_at FROM users WHERE id=$1",
		u.ID).Scan(&u.Username, &u.Email, &u.CreatedAt)
}
