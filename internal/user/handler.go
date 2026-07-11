package user

import (
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Handler struct {
	service     *Service
	i18nService *i18n.Service
}

func NewHandler(service *Service, i18nService *i18n.Service) *Handler {
	return new(Handler{
		service:     service,
		i18nService: i18nService,
	})
}

func (h *Handler) SetLanguage(ctx *th.Context, message telego.Message) error {
	userID := message.From.ID
	lang, err := h.service.GetLanguage(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting user language: %w", err)
	}

	localizer := h.i18nService.FromLanguage(lang)

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Українська").WithCallbackData("lang_ukrainian_button"),
			tu.InlineKeyboardButton("English").WithCallbackData("lang_english_button"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Polski").WithCallbackData("lang_polish_button"),
		),
	)
	if _, err = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(message.Chat.ID),
		i18n.T(localizer, "change_language_message"),
	).WithReplyMarkup(keyboard)); err != nil {
		return fmt.Errorf("sending change_language_message: %w", err)
	}

	return nil
}

func (h *Handler) CallbackHandler(ctx *th.Context, query telego.CallbackQuery) error {
	userID := query.From.ID
	lang, err := h.service.GetLanguage(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting user language: %w", err)
	}

	localizer := h.i18nService.FromLanguage(lang)

	if err = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID)); err != nil {
		return fmt.Errorf("answering callback query: %w", err)
	}

	var language string
	switch query.Data {
	case "lang_ukrainian_button":
		language = "ua"
	case "lang_english_button":
		language = "en"
	case "lang_polish_button":
		language = "pl"
	default:
		return fmt.Errorf("unknown callback data: %s", query.Data)
	}

	if err = h.service.SetLanguage(ctx, userID, language); err != nil {
		return fmt.Errorf("setting language: %w", err)
	}

	if err = ctx.Bot().DeleteMessage(ctx, tu.Delete(
		query.Message.GetChat().ChatID(),
		query.Message.GetMessageID(),
	)); err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	if _, err = ctx.Bot().SendMessage(ctx, tu.Message(
		query.Message.GetChat().ChatID(),
		i18n.T(localizer, "selected_language_message"),
	)); err != nil {
		return fmt.Errorf("sending user language message: %w", err)
	}
	return nil
}
