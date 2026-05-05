package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	log.Fatal((http.ListenAndServe(":3006", nil)))
}
