package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Mode string

const (
	ModeWebhook Mode = "webhook"
	ModePolling Mode = "polling"
)

type Config struct {
	Token         string
	WebhookURL    string
	WebhookSecret string
	Mode          Mode

	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDb   string
}

func Load() Config {
	_ = godotenv.Load()

	mode := Mode(os.Getenv("BOT_MODE"))

	if mode != ModeWebhook && mode != ModePolling {
		log.Printf("Invalid BOT_MODE '%s', defaulting to 'polling'", mode)
		mode = ModePolling
	}

	config := Config{
		Token:         os.Getenv("BOT_TOKEN"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		Mode:          mode,

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
