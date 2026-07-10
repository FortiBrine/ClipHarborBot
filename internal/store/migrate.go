package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/db/migrations"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("creating goose provider: %w", err)
	}
	defer provider.Close()
	_, err = provider.Up(ctx)
	return err
}
