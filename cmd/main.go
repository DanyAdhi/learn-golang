package main

import (
	"log"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/DanyAdhi/learn-golang/internal/config/db"
	"github.com/DanyAdhi/learn-golang/internal/routes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
	config.LoadConfig()

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
	port := config.AppConfig.APP_PORT
	if port == "" {
		log.Fatal("APP_PORT is not set in the .env file")
	}

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
