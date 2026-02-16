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
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Token:         os.Getenv("BOT_TOKEN"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
	}

	if cfg.Token == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	return cfg
}
