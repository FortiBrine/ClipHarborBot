package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/user"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StartHandler struct {
	userMessagesService *user.Service
}

func NewStartHandler(
	userMessagesService *user.Service,
) *StartHandler {
	return &StartHandler{
		userMessagesService: userMessagesService,
	}
}

func (h *StartHandler) Handle(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.userMessagesService.T(ctx, update.Message.From.ID, "start_command"),
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
