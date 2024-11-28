package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DanyAdhi/learn-golang/internal/config/db"
	"github.com/DanyAdhi/learn-golang/internal/config/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. %v", err)
	}
	// Koneksi ke database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()

	// Setup user module
	routes.SetupAuthRouter(router, database)
	routes.SetupUserRouter(router, database)

	// run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT is not set in the .env file")
	}

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
