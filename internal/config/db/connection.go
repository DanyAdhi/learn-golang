package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Connect() (*sql.DB, error) {
	dbConnStr := fmt.Sprintf(
		`user=%s password=%s dbname=%s port=%s sslmode=%s`,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Pastikan koneksi berhasil
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return db, nil
}
