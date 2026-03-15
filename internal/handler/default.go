package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DefaultHandler struct {
	messageService *messages.MessageService
}

func NewDefaultHandler(
	messageService *messages.MessageService,
) *DefaultHandler {
	return &DefaultHandler{
		messageService: messageService,
	}
}

func (h *DefaultHandler) Default(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"unknown_command",
		),
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
