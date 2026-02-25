package handler

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/service"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TiktokHandler struct {
	messageService   *service.MessageService
	tiktokDownloader *service.TikTokDownloader
}

func NewTiktokHandler(messageService *service.MessageService) *TiktokHandler {
	tiktokDownloader, err := service.NewTikTokDownloader()
	if err != nil {
		log.Fatalf("Failed to initialize TikTok downloader: %v", err)
	}

	go func() {
		err = tiktokDownloader.CleanupOldFiles(1 * 60 * 60 * 1000000000)
		if err != nil {
			log.Printf("Failed to cleanup old files: %v", err)
		}
	}()

	return &TiktokHandler{
		messageService:   messageService,
		tiktokDownloader: tiktokDownloader,
	}
}

func (h *TiktokHandler) Handle(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	parts := strings.Fields(text)

	if len(parts) < 2 {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_help",
			),
		})
		if err != nil {
			log.Printf("Failed to send help message: %v", err)
		}
		return
	}

	url := parts[1]

	if !h.tiktokDownloader.IsValidTikTokURL(url) {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_invalid_url",
			),
		})
		if err != nil {
			log.Printf("Failed to send invalid URL message: %v", err)
		}
		return
	}

	statusMsg, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"tiktok_downloading",
		),
	})
	if err != nil {
		log.Printf("Failed to send downloading message: %v", err)
	}

	filePath, err := h.tiktokDownloader.DownloadVideo(ctx, url)
	if err != nil {
		log.Printf("Failed to download TikTok video: %v", err)

		errorMsg := h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"tiktok_download_error",
		)
		if strings.Contains(err.Error(), "max-filesize") {
			errorMsg = h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_size_error",
			)
		}

		_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   errorMsg,
		})

		if err != nil {
			log.Printf("Failed to send error message: %v", err)
		}

		return
	}

	defer func() {
		if err := h.tiktokDownloader.CleanupFile(filePath); err != nil {
			log.Printf("Failed to cleanup file: %v", err)
		}
	}()

	if statusMsg != nil {
		_, err = b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: statusMsg.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_uploading",
			),
		})

		if err != nil {
			log.Printf("Failed to edit status message: %v", err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open video file: %v", err)
		_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_download_error",
			),
		})

		if err != nil {
			log.Printf("Failed to send error message: %v", err)
		}
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(file)

	_, err = b.SendVideo(ctx, &tgbot.SendVideoParams{
		ChatID: update.Message.Chat.ID,
		Video:  &models.InputFileUpload{Filename: filePath, Data: file},
	})

	if err != nil {
		log.Printf("Failed to send video: %v", err)
		_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"tiktok_download_error",
			),
		})

		if err != nil {
			log.Printf("Failed to send error message: %v", err)
		}

		return
	}

	if statusMsg != nil {
		_, _ = b.DeleteMessage(ctx, &tgbot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: statusMsg.ID,
		})
	}
}
