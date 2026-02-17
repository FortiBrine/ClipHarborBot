package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token         string
	WebhookURL    string
	WebhookSecret string

	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDb   string
}

func Load() *Config {
	_ = godotenv.Load()

	config := &Config{
		Token:         os.Getenv("BOT_TOKEN"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),

		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPass: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDb:   os.Getenv("POSTGRES_DB"),
	}

	if config.Token == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	return config
}
