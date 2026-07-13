package user

import (
	"fmt"
	"strings"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

const CallbackDataPrefix = "lang:"

const (
	msgChangeLanguage   = "user.change_language_message"
	msgSelectedLanguage = "user.selected_language_message"
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

	buttons := make([]telego.InlineKeyboardButton, 0, len(i18n.SupportedLanguages))
	for _, l := range i18n.SupportedLanguages {
		buttons = append(buttons, tu.InlineKeyboardButton(l.Name).WithCallbackData(CallbackDataPrefix+l.Code))
	}
	keyboard := tu.InlineKeyboard(tu.InlineKeyboardRow(buttons...))
	if _, err = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(message.Chat.ID),
		i18n.T(localizer, msgChangeLanguage),
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

	language := strings.TrimPrefix(query.Data, CallbackDataPrefix)
	if !i18n.IsSupported(language) {
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
		i18n.T(localizer, msgSelectedLanguage),
	)); err != nil {
		return fmt.Errorf("sending user language message: %w", err)
	}
	return nil
}
