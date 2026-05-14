package main

import (
	"context"
	"log"
	"net/http"
	"pompidou41/go-shortener/internal/config"
	"pompidou41/go-shortener/internal/handler"
	"pompidou41/go-shortener/internal/service"
	"pompidou41/go-shortener/internal/storage"

	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	conf := config.New()

	ctx := context.Background()

	store := storage.New()

	service := service.New(store)

	h := handler.NewHandler(ctx, service)

	port := ":" + conf.Port

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", h.ShortenHandler)
	mux.HandleFunc("/", h.LengthenHandler)

	srv := http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(srv.ListenAndServe())

}
