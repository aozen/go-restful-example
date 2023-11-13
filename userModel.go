package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

func (u *User) getUsers(db *sql.DB) ([]User, error) {
	userRows, err := db.Query("SELECT id, username, email, created_at FROM users")

	if err != nil {
		return nil, err
	}

	defer userRows.Close()

	users := []User{}
	for userRows.Next() {
		var u User
		if err := userRows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (u *User) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, email, password, created_at FROM users WHERE id=$1",
		u.ID).Scan(&u.Username, &u.Email, &u.Password, &u.CreatedAt)
}

func (u *User) createUser(db *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = db.QueryRow(
		"INSERT INTO users(username, email, password, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		u.Username, u.Email, hashedPassword, u.CreatedAt,
	).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) updateUser(db *sql.DB) error {
	if u.Password != "" { // TODO: Instead of two query, make one but SMART
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = db.Exec("UPDATE users SET password=$1 WHERE id=$2", hashedPassword, u.ID)
		if err != nil {
			return err
		}
	}

	if u.Email != "" {
		_, err := db.Exec("UPDATE users SET email=$1 WHERE id=$2", u.Email, u.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) removeUser(db *sql.DB) error {
	result, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
