package downloader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/humanize"
	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/platform"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type VideoHandler struct {
	messageService *messages.MessageService
	downloader     *Downloader
	formatSelector *FormatSelector
	platform       *platform.Platform
	helpMessageKey string
}

func NewVideoHandler(
	messageService *messages.MessageService,
	downloader *Downloader,
	formatSelector *FormatSelector,
	platform *platform.Platform,
	helpKey string,
) *VideoHandler {
	return &VideoHandler{
		messageService: messageService,
		downloader:     downloader,
		formatSelector: formatSelector,
		platform:       platform,
		helpMessageKey: helpKey,
	}
}

func (h *VideoHandler) Handle(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := strings.TrimSpace(update.Message.Text)
	parts := strings.Fields(text)

	var url string
	if len(parts) == 1 {
		url = parts[0]
	} else {
		url = parts[1]
	}

	if !h.platform.IsValidURL(url) {
		_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"invalid_video_url",
			),
		})

		if err != nil {
			log.Printf("Failed to send invalid video url: %v", err)
		}

		return
	}

	formatResult, err := h.formatSelector.ChooseFormat(ctx, &YTDLPFetcher{}, url)
	if err != nil {
		log.Printf("Failed to choose format: %v", err)
		h.sendError(ctx, b, update, "video_format_error")
		return
	}

	if formatResult.Filesize > 0 {
		_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: fmt.Sprintf(
				h.messageService.GetMessage(ctx, update.Message.From.ID, "video_expected_size"),
				humanize.FormatBytes(formatResult.Filesize),
			),
		})

		if err != nil {
			log.Printf("Failed to send video expected size: %v", err)
		}
	} else {
		return
	}

	statusMsg, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			"video_downloading",
		),
	})

	if err != nil {
		log.Printf("Failed to send status message: %v", err)
		return
	}

	filePath, err := h.downloader.DownloadVideo(ctx, DownloadOptions{
		URL:    url,
		Format: formatResult.FormatID,
		Prefix: h.platform.Name,
	})
	if err != nil {
		h.handleDownloadError(ctx, b, update, err)
		return
	}

	defer func(downloader *Downloader, filePath string) {
		err := downloader.CleanupFile(filePath)
		if err != nil {
			log.Printf("Failed to cleanup %s: %v", filePath, err)
		}
	}(h.downloader, filePath)

	if statusMsg != nil {
		_, err := b.EditMessageText(ctx, &tgbot.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: statusMsg.ID,
			Text: h.messageService.GetMessage(
				ctx,
				update.Message.From.ID,
				"video_uploading",
			),
		})
		if err != nil {
			log.Printf("Failed to edit message text: %v", err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		h.sendError(ctx, b, update, "video_download_error")
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
		h.sendError(ctx, b, update, "video_upload_error")
	}
}

func (h *VideoHandler) handleDownloadError(
	ctx context.Context,
	b *tgbot.Bot,
	update *models.Update,
	err error,
) {

	log.Printf("Download error: %v", err)

	msg := "video_download_error"

	switch {
	case errors.Is(err, ErrFileTooLarge):
		msg = "video_size_error"
	case errors.Is(err, ErrInvalidFormat):
		msg = "video_format_error"
	}

	h.sendError(ctx, b, update, msg)
}

func (h *VideoHandler) sendError(
	ctx context.Context,
	b *tgbot.Bot,
	update *models.Update,
	key string,
) {
	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: h.messageService.GetMessage(
			ctx,
			update.Message.From.ID,
			key,
		),
	})

	if err != nil {
		log.Printf("Failed to send error message: %v", err)
	}
}
