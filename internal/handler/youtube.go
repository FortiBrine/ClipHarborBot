package handler

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/FortiBrine/ClipHarborBot/internal/service"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type YouTubeHandler struct {
	messageService *service.MessageService
	downloader     *service.Downloader
}

func NewYouTubeHandler(messageService *service.MessageService) *YouTubeHandler {

	downloader, err := service.NewDownloader(service.PlatformConfig{
		Name:   "youtube",
		Format: "bestvideo[ext=mp4][height<=1080]+bestaudio[ext=m4a]/best",
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/watch\?.*v=[\w-]+`),
			regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
			regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/shorts/[\w-]+`),
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize YouTube downloader: %v", err)
	}

	downloader.StartCleanupWorker(
		10*time.Minute,
		1*time.Hour,
	)

	return &YouTubeHandler{
		messageService: messageService,
		downloader:     downloader,
	}
}

func (h *YouTubeHandler) Handle(ctx context.Context, b *tgbot.Bot, update *models.Update) {
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
				"youtube_help",
			),
		})
		if err != nil {
			log.Printf("Failed to send help message: %v", err)
		}
		return
	}

	url := parts[1]

	if !h.downloader.IsValidURL(url) {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"youtube_invalid_url",
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
			"youtube_downloading",
		),
	})
	if err != nil {
		log.Printf("Failed to send downloading message: %v", err)
	}

	filePath, err := h.downloader.DownloadVideo(ctx, url)
	if err != nil {
		log.Printf("Failed to download YouTube video: %v", err)

		errorMsg := h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"youtube_download_error",
		)
		if strings.Contains(err.Error(), "max-filesize") {
			errorMsg = h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"youtube_size_error",
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
		if err := h.downloader.CleanupFile(filePath); err != nil {
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
				"youtube_uploading",
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
				"youtube_download_error",
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
				"youtube_download_error",
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
