package service

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/repository"
)

type MessageService struct {
	userLanguageRepository *repository.UserLanguageRepository
}

func NewMessageService(
	userLanguageRepository *repository.UserLanguageRepository,
) *MessageService {
	return &MessageService{
		userLanguageRepository: userLanguageRepository,
	}
}

func (s *MessageService) GetMessage(ctx context.Context, userID int64, messageKey string) string {
	userLanguage, err := s.userLanguageRepository.GetUserLanguage(ctx, userID)

	if err != nil {
		log.Printf("Failed to get user language for user %d: %v. Defaulting to Ukrainian.", userID, err)

		for lang := range messages.Messages {
			userLanguage = lang
			break
		}
	}

	return messages.Messages[userLanguage][messageKey]
}
