package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/model"
	"gorm.io/gorm"
)

type UserLanguageRepository struct {
	database *database.Database
}

func NewUserLanguageRepository(database *database.Database) *UserLanguageRepository {
	return &UserLanguageRepository{database: database}
}

func (r *UserLanguageRepository) GetUserLanguage(ctx context.Context, telegramID int64) (string, error) {
	user, err := gorm.G[model.User](r.database.GormDB).Where("id = ?", telegramID).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user language not found: %w", err)
		}

		return "", fmt.Errorf("failed get user language: %w", err)
	}

	return user.Language, nil
}

func (r *UserLanguageRepository) SetUserLanguage(ctx context.Context, telegramID int64, language string) error {
	user, err := gorm.G[model.User](r.database.GormDB).Where("id = ?", telegramID).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = model.User{
				ID:       telegramID,
				Language: language,
			}

			if err = gorm.G[model.User](r.database.GormDB).Create(ctx, &user); err != nil {
				return fmt.Errorf("failed to create user with language: %w", err)
			}
		}

		return fmt.Errorf("failed set user language: %w", err)
	}

	user.Language = language
	_, err = gorm.G[model.User](r.database.GormDB).Where("id = ?", telegramID).Update(ctx, "language", language)

	if err != nil {
		return fmt.Errorf("failed to update user language: %w", err)
	}

	return nil

}
