package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func StartHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Messages[messages.UA]["start_command"],
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
