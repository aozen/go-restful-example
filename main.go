package main

import "os"

func main() {
	myApp := App{}
	myApp.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	myApp.Run(":8080")
}
