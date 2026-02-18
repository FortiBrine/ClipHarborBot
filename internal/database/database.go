package database

import (
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	GormDB *gorm.DB
}

func New(config *config.Config) (*Database, error) {
	db, err := gorm.Open(
		postgres.Open(fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.PostgresHost,
			config.PostgresPort,
			config.PostgresUser,
			config.PostgresPass,
			config.PostgresDb,
		)),
		&gorm.Config{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		GormDB: db,
	}, nil
}
