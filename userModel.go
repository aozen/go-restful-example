package main

import (
	"database/sql"
	"fmt"
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

func (u *User) getUsers(db *sql.DB) ([]User, error) {
	userRows, err := db.Query("SELECT username, email, created_at FROM users")

	if err != nil {
		return nil, err
	}

	defer userRows.Close()

	users := []User{}
	fmt.Println(userRows)
	for userRows.Next() {
		var u User
		if err := userRows.Scan(&u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
