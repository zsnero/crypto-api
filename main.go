package main

import (
	"crypto-api/config"
	"crypto-api/db"
	"crypto-api/handlers"
	"crypto-api/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db.Connect(cfg)
	services.FetchAndStoreCryptoData(cfg)

	r := chi.NewRouter()
	r.Get("/api/cryptocurrencies", handlers.GetCryptocurrencies)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
