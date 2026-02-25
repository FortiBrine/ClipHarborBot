package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/service"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StartHandler struct {
	messageService *service.MessageService
}

func NewStartHandler(
	messageService *service.MessageService,
) *StartHandler {
	return &StartHandler{
		messageService: messageService,
	}
}

func (h *StartHandler) Handle(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"start_command",
		),
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
