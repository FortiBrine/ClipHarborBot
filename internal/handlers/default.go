package handlers

import (
	"context"
	"log"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
)

func Default(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Messages[messages.UA]["unknown_command"],
	})

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
