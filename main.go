package main

import (
	"anipics/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Простые эндпоинты
	mux.HandleFunc("GET /api/random/{tag...}", handlers.RandomPicProxy)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
