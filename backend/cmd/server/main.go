package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/thankly/backend/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	log.Printf("Thankly API starting on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
