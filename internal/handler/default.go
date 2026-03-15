package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/user"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DefaultHandler struct {
	userMessagesService *user.Service
}

func NewDefaultHandler(
	userMessagesService *user.Service,
) *DefaultHandler {
	return &DefaultHandler{
		userMessagesService: userMessagesService,
	}
}

func (h *DefaultHandler) Default(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.userMessagesService.T(ctx, update.Message.From.ID, "unknown_command"),
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
