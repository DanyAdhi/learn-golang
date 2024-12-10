package main

import (
	"log"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/config"
	"github.com/DanyAdhi/learn-golang/internal/config/db"
	"github.com/DanyAdhi/learn-golang/internal/config/validator"
	"github.com/DanyAdhi/learn-golang/internal/routes"
	"github.com/DanyAdhi/learn-golang/internal/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Driver PostgreSQL
	"github.com/rs/cors"
)

func main() {
	validator.InitValidator()

	config.LoadConfig()

	// Koneksi ke database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()
	router.Use(utils.LoggingMiddleware)

	// Setup user module
	routes.SetupAuthRouter(router, database)
	routes.SetupUserRouter(router, database)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Ubah sesuai domain front-end Anda
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsMiddleware.Handler(router)

	// run server
	port := config.AppConfig.APP_PORT
	if port == "" {
		log.Fatal("APP_PORT is not set in the .env file")
	}

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
