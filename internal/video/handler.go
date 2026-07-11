package video

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/dustin/go-humanize"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Handler struct {
	userService    *user.Service
	i18nService    *i18n.Service
	fetcher        MetadataFetcher
	downloader     Downloader
	formatSelector *FormatSelector
}

func NewHandler(
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher MetadataFetcher,
	downloader Downloader,
	formatSelector *FormatSelector,
) *Handler {
	return new(Handler{
		userService:    userService,
		i18nService:    i18nService,
		fetcher:        fetcher,
		downloader:     downloader,
		formatSelector: formatSelector,
	})
}

func (h *Handler) Handle(platform *Platform, helpKey string) th.MessageHandler {
	return func(
		ctx *th.Context,
		message telego.Message,
	) error {
		userID := message.From.ID
		chatID := message.Chat.ChatID()
		lang, err := h.userService.GetLanguage(ctx, userID)
		if err != nil {
			return fmt.Errorf("getting user language: %w", err)
		}

		localizer := h.i18nService.FromLanguage(lang)

		text := strings.TrimSpace(message.Text)
		parts := strings.Fields(text)

		var url string
		if len(parts) == 1 {
			url = parts[0]
		} else {
			url = parts[1]
		}

		if !platform.IsValidURL(url) {
			if _, err = ctx.Bot().SendMessage(ctx, tu.Message(
				chatID,
				i18n.T(localizer, "invalid_video_url"),
			)); err != nil {
				return fmt.Errorf("sending message to user: %w", err)
			}

			return nil
		}

		format, err := h.formatSelector.FetchBest(ctx, h.fetcher, url)
		if err != nil {
			h.sendError(ctx, ctx.Bot(), chatID, localizer, "video_format_error")
			return fmt.Errorf("fetching best format: %w", err)
		}

		if format.Filesize <= 0 {
			h.sendError(ctx, ctx.Bot(), chatID, localizer, "video_format_error")
			return fmt.Errorf("format not available")
		}

		if _, err = ctx.Bot().SendMessage(ctx, tu.Messagef(
			chatID,
			i18n.T(localizer, "video_expected_size"),
			humanize.IBytes(uint64(format.Filesize)),
		)); err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}

		statusMsg, err := ctx.Bot().SendMessage(ctx, tu.Message(
			chatID,
			i18n.T(localizer, "video_downloading"),
		))
		if err != nil {
			return fmt.Errorf("sending message to user: %w", err)
		}

		filePath, err := h.downloader.DownloadVideo(ctx, DownloadOptions{
			URL:    url,
			Format: format.FormatID,
			Prefix: platform.Name,
		})
		if err != nil {
			h.handleDownloadError(ctx, ctx.Bot(), chatID, localizer, err)
			return fmt.Errorf("downloading video: %w", err)
		}
		defer h.downloader.CleanupFile(filePath)

		if _, err = ctx.Bot().EditMessageText(ctx, tu.EditMessageText(
			chatID,
			statusMsg.MessageID,
			i18n.T(localizer, "video_uploading"),
		)); err != nil {
			return fmt.Errorf("editing message to user: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			h.sendError(ctx, ctx.Bot(), chatID, localizer, "video_download_error")
			return fmt.Errorf("opening file: %w", err)
		}
		defer file.Close()

		if _, err = ctx.Bot().SendVideo(ctx, tu.Video(
			chatID,
			tu.File(file),
		)); err != nil {
			h.sendError(ctx, ctx.Bot(), chatID, localizer, "video_upload_error")
			return fmt.Errorf("sending video to user: %w", err)
		}

		return nil
	}
}

func (h *Handler) handleDownloadError(
	ctx context.Context,
	b *telego.Bot,
	chatID telego.ChatID,
	localizer *i18n.Localizer,
	err error,
) error {
	msg := "video_download_error"

	switch {
	case errors.Is(err, ErrFileTooLarge):
		msg = "video_size_error"
	case errors.Is(err, ErrInvalidFormat):
		msg = "video_format_error"
	}

	if sendingErr := h.sendError(ctx, b, chatID, localizer, msg); sendingErr != nil {
		return fmt.Errorf("sending error: %w", sendingErr)
	}
	return nil
}

func (h *Handler) sendError(
	ctx context.Context,
	bot *telego.Bot,
	chatID telego.ChatID,
	localizer *i18n.Localizer,
	key string,
) error {
	if _, err := bot.SendMessage(ctx, tu.Message(
		chatID,
		i18n.T(localizer, key),
	)); err != nil {
		return fmt.Errorf("sending message to user: %w", err)
	}
	return nil
}
