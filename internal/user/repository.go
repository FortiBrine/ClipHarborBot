package user

import (
	"context"
	"errors"

	"github.com/FortiBrine/ClipHarborBot/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	GetLanguage(context.Context, int64) (string, error)
	SetLanguage(context.Context, int64, string) error
}

type PostgresRepository struct {
	q *store.Queries
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return new(PostgresRepository{q: store.New(pool)})
}

func (r *PostgresRepository) GetLanguage(ctx context.Context, telegramID int64) (string, error) {
	language, err := r.q.GetUserLanguage(ctx, telegramID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", err
	}

	return language, nil
}

func (r *PostgresRepository) SetLanguage(
	ctx context.Context,
	telegramID int64,
	language string,
) error {
	return r.q.UpsertUserLanguage(ctx, store.UpsertUserLanguageParams{
		ID:       telegramID,
		Language: language,
	})
}
