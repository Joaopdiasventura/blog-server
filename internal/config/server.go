package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN       string
	Port      string
	JWTSecret string
}

func Load() Config {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	if port == "" {
		port = "8080"
	}
	return Config{DSN: dsn, Port: port, JWTSecret: jwtSecret}
}
