package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DanyAdhi/learn-golang/internal/config"
)

func Connect() (*sql.DB, error) {
	dbConnStr := fmt.Sprintf(
		`user=%s password=%s dbname=%s port=%s sslmode=%s`,
		config.AppConfig.DB_USER,
		config.AppConfig.DB_PASSWORD,
		config.AppConfig.DB_NAME,
		config.AppConfig.DB_PORT,
		config.AppConfig.DB_SSLMODE,
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
