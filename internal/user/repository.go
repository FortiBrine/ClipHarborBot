package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"gorm.io/gorm"
)

type LanguageRepository struct {
	database *database.Database
}

func NewUserLanguageRepository(database *database.Database) *LanguageRepository {
	return &LanguageRepository{database: database}
}

func (r *LanguageRepository) Migrate() error {
	err := r.database.GormDB.AutoMigrate(&User{})

	if err != nil {
		return fmt.Errorf("failed to migrate user language repository: %w", err)
	}

	return nil
}

func (r *LanguageRepository) GetUserLanguage(ctx context.Context, telegramID int64) (string, error) {
	user, err := gorm.G[User](r.database.GormDB).Where("id = ?", telegramID).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user language not found: %w", err)
		}

		return "", fmt.Errorf("failed get user language: %w", err)
	}

	return user.Language, nil
}

func (r *LanguageRepository) SetUserLanguage(ctx context.Context, telegramID int64, language string) error {
	user, err := gorm.G[User](r.database.GormDB).Where("id = ?", telegramID).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = User{
				ID:       telegramID,
				Language: language,
			}

			if err = gorm.G[User](r.database.GormDB).Create(ctx, &user); err != nil {
				return fmt.Errorf("failed to create user with language: %w", err)
			}

			return nil
		}

		return fmt.Errorf("failed set user language: %w", err)
	}

	user.Language = language
	_, err = gorm.G[User](r.database.GormDB).Where("id = ?", telegramID).Update(ctx, "language", language)

	if err != nil {
		return fmt.Errorf("failed to update user language: %w", err)
	}

	return nil

}
