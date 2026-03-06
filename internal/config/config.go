package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl     string
	Port      string
	JWTSecret string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return Config{
		DBUrl:     dbURL,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}
