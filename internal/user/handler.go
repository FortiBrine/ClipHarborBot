package user

import (
	"context"
	"log"
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	logger      *slog.Logger
	service     *Service
	i18nService *i18n.Service
}

func NewHandler(logger *slog.Logger, service *Service, i18nService *i18n.Service) *Handler {
	return new(Handler{
		logger:      logger,
		service:     service,
		i18nService: i18nService,
	})
}

func (h *Handler) SetLanguage(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	lang, err := h.service.GetLanguage(ctx, update.Message.From.ID)
	if err != nil {
		h.logger.Error("getting user language", "error", err)
		return
	}

	localizer := h.i18nService.FromLanguage(lang)

	keyboardButtons := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Українська", CallbackData: "lang_ukrainian_button"},
				{Text: "English", CallbackData: "lang_english_button"},
			}, {
				{Text: "Polski", CallbackData: "lang_polish_button"},
			},
		},
	}

	if _, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        i18n.T(localizer, "change_language_message"),
		ReplyMarkup: keyboardButtons,
	}); err != nil {
		h.logger.Error("sending user language message", "error", err)
	}
}

func (h *Handler) CallbackHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	callbackQuery := update.CallbackQuery
	if callbackQuery == nil {
		return
	}
	lang, err := h.service.GetLanguage(ctx, update.Message.From.ID)
	if err != nil {
		h.logger.Error("getting user language", "error", err)
		return
	}

	localizer := h.i18nService.FromLanguage(lang)

	if _, err := b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	}); err != nil {
		h.logger.Error("answering callback query", "error", err)
		return
	}

	data := update.CallbackQuery.Data
	var language string
	switch data {
	case "lang_ukrainian_button":
		language = "ua"
	case "lang_english_button":
		language = "en"
	case "lang_polish_button":
		language = "pl"
	default:
		log.Printf("Unknown callback data: %s", data)
		return
	}

	if err := h.service.SetLanguage(ctx, callbackQuery.From.ID, language); err != nil {
		h.logger.Error("setting language", "error", err)
		return
	}

	if _, err := b.DeleteMessage(ctx, &tgbot.DeleteMessageParams{
		ChatID:    callbackQuery.Message.Message.Chat.ID,
		MessageID: callbackQuery.Message.Message.ID,
	}); err != nil {
		h.logger.Error("deleting message", "error", err)
	}

	if _, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: callbackQuery.Message.Message.Chat.ID,
		Text:   i18n.T(localizer, "selected_language_message"),
	}); err != nil {
		h.logger.Error("sending user language message", "error", err)
	}
}
