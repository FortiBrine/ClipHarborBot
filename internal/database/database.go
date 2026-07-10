package database

import (
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(config config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.PostgresHost,
			config.PostgresPort,
			config.PostgresUser,
			config.PostgresPass,
			config.PostgresDb,
		)),
		new(gorm.Config),
	)

	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	return db, nil
}
