package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT                 string
	DB_USER                  string
	DB_PASSWORD              string
	DB_NAME                  string
	DB_PORT                  string
	DB_SSLMODE               string
	JWT_SECRET_ACCESS_TOKEN  string
	JWT_SECRET_REFRESH_TOKEN string
	REDIS_HOST               string
	REDIS_PASSWORD           string
	REDIS_PORT               string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file. %v", err)
	}

	AppConfig = Config{
		APP_PORT:                 os.Getenv("APP_PORT"),
		DB_USER:                  os.Getenv("DB_USER"),
		DB_PASSWORD:              os.Getenv("DB_PASSWORD"),
		DB_NAME:                  os.Getenv("DB_NAME"),
		DB_PORT:                  os.Getenv("DB_PORT"),
		DB_SSLMODE:               os.Getenv("DB_SSLMODE"),
		JWT_SECRET_ACCESS_TOKEN:  os.Getenv("JWT_SECRET_ACCESS_TOKEN"),
		JWT_SECRET_REFRESH_TOKEN: os.Getenv("JWT_SECRET_REFRESH_TOKEN"),
		REDIS_HOST:               os.Getenv("REDIS_HOST"),
		REDIS_PASSWORD:           os.Getenv("REDIS_PASSWORD"),
		REDIS_PORT:               os.Getenv("REDIS_PORT"),
	}

}
