package handlers

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/services"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var tiktokDownloader *services.TikTokDownloader

func init() {
	var err error
	tiktokDownloader, err = services.NewTikTokDownloader()
	if err != nil {
		log.Fatal("Failed to initialize TikTok downloader: ", err)
	}

	go func() {
		_ = tiktokDownloader.CleanupOldFiles(1 * 60 * 60 * 1000000000)
	}()
}

func Tiktok(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	parts := strings.Fields(text)

	if len(parts) < 2 {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.Messages[messages.UA]["tiktok_help"],
		})
		if err != nil {
			log.Printf("Failed to send help message: %v", err)
		}
		return
	}

	url := parts[1]

	if !tiktokDownloader.IsValidTikTokURL(url) {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.Messages[messages.UA]["tiktok_invalid_url"],
		})
		if err != nil {
			log.Printf("Failed to send invalid URL message: %v", err)
		}
		return
	}

	statusMsg, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Messages[messages.UA]["tiktok_downloading"],
	})
	if err != nil {
		log.Printf("Failed to send downloading message: %v", err)
	}

	filePath, err := tiktokDownloader.DownloadVideo(ctx, url)
	if err != nil {
		log.Printf("Failed to download TikTok video: %v", err)

		errorMsg := messages.Messages[messages.UA]["tiktok_download_error"]
		if strings.Contains(err.Error(), "max-filesize") {
			errorMsg = messages.Messages[messages.UA]["tiktok_size_error"]
		}

		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   errorMsg,
		})
		return
	}

	defer func() {
		if err := tiktokDownloader.CleanupFile(filePath); err != nil {
			log.Printf("Failed to cleanup file: %v", err)
		}
	}()

	if statusMsg != nil {
		_, _ = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: statusMsg.ID,
			Text:      messages.Messages[messages.UA]["tiktok_uploading"],
		})
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open video file: %v", err)
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.Messages[messages.UA]["tiktok_download_error"],
		})
		return
	}
	defer file.Close()

	_, err = b.SendVideo(ctx, &tgbot.SendVideoParams{
		ChatID: update.Message.Chat.ID,
		Video:  &models.InputFileUpload{Filename: filePath, Data: file},
	})

	if err != nil {
		log.Printf("Failed to send video: %v", err)
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   messages.Messages[messages.UA]["tiktok_download_error"],
		})
		return
	}

	if statusMsg != nil {
		_, _ = b.DeleteMessage(ctx, &tgbot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: statusMsg.ID,
		})
	}
}
