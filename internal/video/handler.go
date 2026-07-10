package video

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/dustin/go-humanize"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	logger         *slog.Logger
	userService    *user.Service
	i18nService    *i18n.Service
	fetcher        MetadataFetcher
	downloader     Downloader
	formatSelector *FormatSelector
}

func NewHandler(
	logger *slog.Logger,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher MetadataFetcher,
	downloader Downloader,
	formatSelector *FormatSelector,
) *Handler {
	return new(Handler{
		logger:         logger,
		userService:    userService,
		i18nService:    i18nService,
		fetcher:        fetcher,
		downloader:     downloader,
		formatSelector: formatSelector,
	})
}

func (h *Handler) Handle(platform *Platform, helpKey string) tgbot.HandlerFunc {
	return func(
		ctx context.Context,
		bot *tgbot.Bot,
		update *models.Update,
	) {
		if update.Message == nil {
			return
		}
		lang, err := h.userService.GetLanguage(ctx, update.Message.From.ID)
		if err != nil {
			h.logger.Error("getting user language", "error", err)
			return
		}

		localizer := h.i18nService.FromLanguage(lang)

		text := strings.TrimSpace(update.Message.Text)
		parts := strings.Fields(text)

		var url string
		if len(parts) == 1 {
			url = parts[0]
		} else {
			url = parts[1]
		}

		if !platform.IsValidURL(url) {
			if _, err = bot.SendMessage(ctx, new(tgbot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   i18n.T(localizer, "invalid_video_url"),
			})); err != nil {
				h.logger.Error("sending message to user", "error", err)
			}

			return
		}

		format, err := h.formatSelector.FetchBest(ctx, h.fetcher, url)
		if err != nil {
			h.logger.Error("fetching best format", "error", err)
			h.sendError(ctx, bot, update, localizer, "video_format_error")
			return
		}

		if format.Filesize <= 0 {
			return
		}

		if _, err = bot.SendMessage(ctx, new(tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: fmt.Sprintf(
				i18n.T(localizer, "video_expected_size"),
				humanize.IBytes(uint64(format.Filesize)),
			),
		})); err != nil {
			h.logger.Error("sending message to user", "error", err)
		}

		statusMsg, err := bot.SendMessage(ctx, new(tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   i18n.T(localizer, "video_downloading"),
		}))

		if err != nil {
			h.logger.Error("sending message to user", "error", err)
		}

		filePath, err := h.downloader.DownloadVideo(ctx, DownloadOptions{
			URL:    url,
			Format: format.FormatID,
			Prefix: platform.Name,
		})
		if err != nil {
			h.logger.Error("downloading video", "error", err)
			h.handleDownloadError(ctx, bot, update, localizer, err)
			return
		}
		defer h.downloader.CleanupFile(filePath)

		if statusMsg != nil {
			if _, err = bot.EditMessageText(ctx, new(tgbot.EditMessageTextParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: statusMsg.ID,
				Text:      i18n.T(localizer, "video_uploading"),
			})); err != nil {
				h.logger.Error("editing message to user", "error", err)
			}
		}

		file, err := os.Open(filePath)
		if err != nil {
			h.logger.Error("opening file", "error", err)
			h.sendError(ctx, bot, update, localizer, "video_download_error")
			return
		}
		defer file.Close()

		if _, err = bot.SendVideo(ctx, new(tgbot.SendVideoParams{
			ChatID: update.Message.Chat.ID,
			Video:  new(models.InputFileUpload{Filename: filePath, Data: file}),
		})); err != nil {
			h.logger.Error("sending video", "error", err)
			h.sendError(ctx, bot, update, localizer, "video_upload_error")
		}
	}
}

func (h *Handler) handleDownloadError(
	ctx context.Context,
	b *tgbot.Bot,
	update *models.Update,
	localizer *i18n.Localizer,
	err error,
) {
	msg := "video_download_error"

	switch {
	case errors.Is(err, ErrFileTooLarge):
		msg = "video_size_error"
	case errors.Is(err, ErrInvalidFormat):
		msg = "video_format_error"
	}

	h.sendError(ctx, b, update, localizer, msg)
}

func (h *Handler) sendError(
	ctx context.Context,
	bot *tgbot.Bot,
	update *models.Update,
	localizer *i18n.Localizer,
	key string,
) {
	if _, err := bot.SendMessage(ctx, new(tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   i18n.T(localizer, key),
	})); err != nil {
		h.logger.Error("sending error message to user", "error", err)
	}
}
