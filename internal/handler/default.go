package handler

import (
	"context"
	"log"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/service"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type DefaultHandler struct {
	messageService *service.MessageService
	tiktokHandler  *TiktokHandler
	youtubeHandler *YouTubeHandler
}

func NewDefaultHandler(
	messageService *service.MessageService,
	tiktokHandler *TiktokHandler,
	youtubeHandler *YouTubeHandler,
) *DefaultHandler {
	return &DefaultHandler{
		messageService: messageService,
		tiktokHandler:  tiktokHandler,
		youtubeHandler: youtubeHandler,
	}
}

func (h *DefaultHandler) Default(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)

	if text != "" {
		if h.tiktokHandler.downloader.IsValidURL(text) {
			update.Message.Text = "/tiktok " + text
			h.tiktokHandler.Handle(ctx, b, update)
			return
		}

		if h.youtubeHandler.downloader.IsValidURL(text) {
			update.Message.Text = "/youtube " + text
			h.youtubeHandler.Handle(ctx, b, update)
			return
		}
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
		log.Fatalf("Failed to send message: %v", err)
	}
}
