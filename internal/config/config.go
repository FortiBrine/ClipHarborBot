package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment string

const (
	EnvDev  Environment = "dev"
	EnvProd Environment = "prod"
)

func (e Environment) IsDev() bool { return e == EnvDev }

type Config struct {
	Environment         Environment   `env:"ENVIRONMENT" envDefault:"dev"`
	DefaultLang         string        `env:"DEFAULT_LANG" envDefault:"en"`
	DownloadTimeout     time.Duration `env:"DOWNLOAD_TIMEOUT" envDefault:"5m"`
	HttpShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`

	TelegramToken string `env:"TELEGRAM_TOKEN,required"`
	WebhookURL    string `env:"WEBHOOK_URL"`
	WebhookSecret string `env:"WEBHOOK_SECRET"`
	HttpAddress   string `env:"HTTP_ADDRESS" envDefault:":2000"`

	PostgresHost string `env:"POSTGRES_HOST,required"`
	PostgresPort string `env:"POSTGRES_PORT,required"`
	PostgresUser string `env:"POSTGRES_USER,required"`
	PostgresPass string `env:"POSTGRES_PASS,required"`
	PostgresDb   string `env:"POSTGRES_DB,required"`
}

func Load() (Config, error) {
	godotenv.Load()
	return env.ParseAs[Config]()
}
