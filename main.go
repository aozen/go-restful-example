package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"
)

func main() {

	err := loadEnvVariables()
	if err != nil {
		log.Fatal("Error loading environment variables: ", err)
	}

	db, err := sql.Open("postgres", getDBConnectionString())
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	myApp := App{}
	myApp.Initialize(db)
	myApp.Run()
}

func loadEnvVariables() error {
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
