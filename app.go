package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Run() {
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}

func (app *App) Initialize(db *sql.DB) {
	app.DB = db
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("APP_DB_HOST"), os.Getenv("APP_DB_PORT"), os.Getenv("APP_DB_NAME"), os.Getenv("APP_DB_USER"), os.Getenv("APP_DB_PASSWORD"))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/user/{id}", app.getUser).Methods("GET") // TODO: Try like REGEX. Check just digits, otherwise 404
}

func (app *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"]) // TODO: Check better ways to do.
	if err != nil {
		fmt.Println("Invalid User ID") // TODO: Return response to the writer.
		return
	}

	user := User{ID: id}
	if err := user.getUser(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("User Not Found") // TODO: Return response to the writer.
		default:
			fmt.Println(err.Error()) // TODO: Return response to the writer.
		}
		return
	}

	fmt.Println(user) // TODO: Return response to the writer.
}
