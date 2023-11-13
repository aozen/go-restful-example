package main

import (
	"database/sql"
	"encoding/json"
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
	app.Router.HandleFunc("/users", app.getUsers).Methods("GET")
	app.Router.HandleFunc("/users", app.createUser).Methods("POST")
	app.Router.HandleFunc("/users/{id:[0-9]+}", app.getUser).Methods("GET")
	app.Router.HandleFunc("/users/{id:[0-9]+}", app.updateUser).Methods("PUT")
	app.Router.HandleFunc("/users/{id:[0-9]+}", app.removeUser).Methods("DELETE")
}

/*
Example content:
[{0 asdasd@gmail.com asf  2023-09-10T12:00:00Z} {0 asd@asd.sad 3335555  2023-11-09T13:25:06.713213Z}]

Example response:
[91 123 34 105 100 34 58 48 44 34 117 ... 93]
*/
func responseJson(w http.ResponseWriter, code int, content interface{}) { // Check <var>interface{} is for dynamic params?
	response, _ := json.Marshal(content)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"]) // TODO: Check better ways to do.
	if err != nil {
		responseJson(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := User{ID: id}
	if err := user.getUser(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responseJson(w, http.StatusNotFound, "User Not Found")
		default:
			responseJson(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responseJson(w, http.StatusOK, user)
}

func (app *App) getUsers(w http.ResponseWriter, r *http.Request) {
	u := User{}
	users, err := u.getUsers(app.DB)
	if err != nil {
		responseJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(w, http.StatusOK, users)
}

func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&user); err != nil {
		responseJson(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" { //TODO: Separate the errors
		responseJson(w, http.StatusBadRequest, "Username, Email, and Password are required fields")
		return
	}

	err := user.createUser(app.DB)
	if err != nil {
		responseJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson(w, http.StatusCreated, user)
}

func (app *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"]) // TODO: Check better ways to do.
	if err != nil {
		responseJson(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	user := User{ID: id}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err = decoder.Decode(&user); err != nil {
		responseJson(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := user.updateUser(app.DB); err != nil {
		responseJson(w, http.StatusAccepted, err.Error())
		return
	}

	responseJson(w, http.StatusAccepted, "User updated successfully")
}

func (app *App) removeUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"]) // TODO: Check better ways to do.
	if err != nil {
		responseJson(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user := User{ID: id}
	if err := user.removeUser(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			responseJson(w, http.StatusNotFound, "User Not Found")
		default:
			responseJson(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responseJson(w, http.StatusNoContent, "User removed successfully")
}
