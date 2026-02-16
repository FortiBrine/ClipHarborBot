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

	config := &Config{
		Token:         os.Getenv("BOT_TOKEN"),
		WebhookURL:    os.Getenv("WEBHOOK_URL"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
	}

	if config.Token == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	return config
}
