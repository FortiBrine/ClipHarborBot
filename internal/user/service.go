package user

import (
	"context"
	"errors"
)

type Service struct {
	repository  Repository
	defaultLang string
}

func NewService(
	repository Repository,
	defaultLang string,
) *Service {
	return new(Service{
		repository:  repository,
		defaultLang: defaultLang,
	})
}

func (s *Service) GetLanguage(ctx context.Context, telegramID int64) (string, error) {
	lang, err := s.repository.GetLanguage(ctx, telegramID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return s.defaultLang, nil
		}

		return "", err
	}

	return lang, nil
}

func (s *Service) SetLanguage(ctx context.Context, telegramID int64, lang string) error {
	return s.repository.SetLanguage(ctx, telegramID, lang)
}
