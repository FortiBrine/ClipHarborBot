package user

import (
	"context"
	"errors"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"gorm.io/gorm"
)

type Service struct {
	repository LanguageRepository
	i18n       *i18n.I18n
}

func NewService(
	repository LanguageRepository,
	i18n *i18n.I18n,
) *Service {
	return &Service{
		repository: repository,
		i18n:       i18n,
	}
}

func (s *Service) GetLanguage(ctx context.Context, telegramID int64) string {
	lang, err := s.repository.GetUserLanguage(ctx, telegramID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Failed to get user language: %v", err)
		}

		return s.i18nDefault()
	}

	return lang
}

func (s *Service) SetLanguage(ctx context.Context, telegramID int64, lang string) error {
	return s.repository.SetUserLanguage(ctx, telegramID, lang)
}

func (s *Service) T(ctx context.Context, telegramID int64, key string) string {
	lang := s.GetLanguage(ctx, telegramID)
	return s.i18n.T(lang, key)
}

func (s *Service) i18nDefault() string {
	return s.i18n.DefaultLang
}
