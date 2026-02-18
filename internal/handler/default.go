package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/repository"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
)

type DefaultHandler struct {
	userLanguageRepository *repository.UserLanguageRepository
}

func NewDefaultHandler(userLanguageRepository *repository.UserLanguageRepository) *DefaultHandler {
	return &DefaultHandler{
		userLanguageRepository: userLanguageRepository,
	}
}

func (handler *DefaultHandler) Default(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	lang, err := handler.userLanguageRepository.GetUserLanguage(ctx, update.Message.From.ID)

	if err != nil {
		log.Printf("Failed to get user language: %v", err)
		lang = "ua"
	}

	_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   lang,
	})

	_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Messages[messages.UA]["unknown_command"],
	})

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
