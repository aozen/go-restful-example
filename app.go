package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() error {
	connStr := getDBConnectionString()
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	a.DB = db

	return nil
}

func (a *App) Run(addr string) {}

func getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("APP_DB_HOST"), os.Getenv("APP_DB_PORT"), os.Getenv("APP_DB_NAME"), os.Getenv("APP_DB_USER"), os.Getenv("APP_DB_PASSWORD"))
}
