package user

import (
	"context"

	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LanguageRepository interface {
	Migrate() error
	GetUserLanguage(context.Context, int64) (string, error)
	SetUserLanguage(context.Context, int64, string) error
}

type PostgresLanguageRepository struct {
	database *database.Database
}

func NewPostgresLanguageRepository(database *database.Database) LanguageRepository {
	return &PostgresLanguageRepository{database: database}
}

func (r *PostgresLanguageRepository) Migrate() error {
	return r.database.GormDB.AutoMigrate(&User{})
}

func (r *PostgresLanguageRepository) GetUserLanguage(ctx context.Context, telegramID int64) (string, error) {
	userLanguage, err := gorm.G[User](r.database.GormDB).Where("id = ?", telegramID).First(ctx)
	if err != nil {
		return "", err
	}

	return userLanguage.Language, nil
}

func (r *PostgresLanguageRepository) SetUserLanguage(ctx context.Context, telegramID int64, language string) error {
	return gorm.G[User](r.database.GormDB, clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"language"}),
	}).Create(ctx, &User{
		ID:       telegramID,
		Language: language,
	})
}
