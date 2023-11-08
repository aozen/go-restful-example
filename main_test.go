package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
	"strconv"
)

var TestApp App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`

func TestMain(m *testing.M) {
	loadTestEnvVariables()

	db, _ := sql.Open(
		"postgres",
		getTestDBConnectionString(),
	)

	TestApp.Initialize(db)

	checkUserTableExists()

	// Run the tests
	code := m.Run()

	clearUserTable()

	// If all tests are success return 0, otherwise 1
	os.Exit(code)
}

func getTestDBConnectionString() string {
	//host=localhost port=5432 dbname=test_pq_restful_api user=admin password=admin123 sslmode=disable
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"))
}

func loadTestEnvVariables() error { // Convert this to service or common place to run.
	file, err := os.Open(".env")

	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func checkUserTableExists() {
	_, err := TestApp.DB.Exec(tableCreationQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearUserTable() {
	TestApp.DB.Exec("DELETE FROM users")
	TestApp.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		TestApp.DB.Exec("INSERT INTO users(username, email, password, created_at) VALUES($1, $2, $3, $4)", "Username_"+strconv.Itoa(i), "email"+strconv.Itoa((i+1.0)*10)+"@gmail.com", "*******", "09-10-2023 12:00:00")
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetUser(t *testing.T) {
	clearUserTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()
	TestApp.Router.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusOK, rr.Code)
}
