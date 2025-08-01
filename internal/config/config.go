package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr   string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
}

func Load() (*Config, error) {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: .env file not found, reading from environment variables: %s\n", err.Error())
		}
	}

	jwtExpiresIn, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_MINUTES", "60"))
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddr:   getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", "defaultsecret"),
		JWTExpiresIn: time.Duration(jwtExpiresIn) * time.Minute,
	}, nil

}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
