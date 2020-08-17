package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/walln/flurry2/flurry"
)

func main() {

	//DEV

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := flurry.Initialize()
	server.ListenAndServe()

}
