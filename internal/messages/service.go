package messages

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/user"
)

type MessageService struct {
	userLanguageRepository *user.LanguageRepository
}

func NewMessageService(
	userLanguageRepository *user.LanguageRepository,
) *MessageService {
	return &MessageService{
		userLanguageRepository: userLanguageRepository,
	}
}

func (s *MessageService) GetMessage(ctx context.Context, userID int64, messageKey string) string {
	userLanguage, err := s.userLanguageRepository.GetUserLanguage(ctx, userID)

	if err != nil {
		log.Printf("Failed to get user language for user %d: %v. Defaulting to Ukrainian.", userID, err)

		for lang := range Messages {
			userLanguage = lang
			break
		}
	}

	return Messages[userLanguage][messageKey]
}
