package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT    string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	DB_SSLMODE  string
	JWT_SECRET  string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. %v", err)
	}

	AppConfig = Config{
		APP_PORT:    os.Getenv("APP_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
	}

}
