package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	Migrate() error
	GetLanguage(context.Context, int64) (string, error)
	SetLanguage(context.Context, int64, string) error
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return new(PostgresRepository{db: db})
}

func (r *PostgresRepository) Migrate() error {
	return r.db.AutoMigrate(new(User))
}

func (r *PostgresRepository) GetLanguage(ctx context.Context, telegramID int64) (string, error) {
	u, err := gorm.G[User](r.db).Where("id = ?", telegramID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUserNotFound
		}

		return "", err
	}

	return u.Language, nil
}

func (r *PostgresRepository) SetLanguage(
	ctx context.Context,
	telegramID int64,
	language string,
) error {
	return gorm.G[User](r.db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"language"}),
	}).Create(ctx, new(User{
		ID:       telegramID,
		Language: language,
	}))
}
