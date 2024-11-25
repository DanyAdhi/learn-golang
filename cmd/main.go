package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DanyAdhi/learn-golang/internal/users"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
	// Koneksi ke database
	db, err := sql.Open("postgres", "user=root password=root1234 dbname=learn port=5439 sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()

	// Setup user module
	usersRepo := users.NewRepository(db)
	usersService := users.NewService(usersRepo)
	userHandler := users.NewHandler(usersService)

	// Routing
	router.HandleFunc("/users", userHandler.GetAllUsersHandler).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", userHandler.GetOneUsersHandler).Methods(http.MethodGet)
	router.HandleFunc("/users", userHandler.CreateUsersHandler).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", userHandler.UpdateUsersHandler).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", userHandler.DeleteusersHandler).Methods(http.MethodDelete)

	// Jalankan server
	log.Println("Server running on port 3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}
