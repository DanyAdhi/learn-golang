package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DanyAdhi/learn-golang/internal/auth"
	"github.com/DanyAdhi/learn-golang/internal/users"
	"github.com/DanyAdhi/learn-golang/internal/utils"
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
	database, err := sql.Open("postgres", "user=root password=root1234 dbname=learn port=5439 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()

	// Setup user module
	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	router.HandleFunc("/sign-in", authHandler.LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/refresh-token", authHandler.RefreshTokenHandler).Methods(http.MethodPost)

	// Setup user module
	usersRepo := users.NewRepository(database)
	usersService := users.NewService(usersRepo)
	userHandler := users.NewHandler(usersService)

	// Routing
	router.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(userHandler.GetAllUsersHandler))).Methods(http.MethodGet)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.GetOneUsersHandler))).Methods(http.MethodGet)
	router.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(userHandler.CreateUsersHandler))).Methods(http.MethodPost)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.UpdateUsersHandler))).Methods(http.MethodPut)
	router.Handle("/users/{id}", utils.AuthMiddleware(http.HandlerFunc(userHandler.DeleteusersHandler))).Methods(http.MethodDelete)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT is not set in the .env file")
	}

	log.Println("Server running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
