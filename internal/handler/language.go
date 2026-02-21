package handler

import (
	"context"
	"log"

	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/repository"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type LanguageHandler struct {
	userLanguageRepository *repository.UserLanguageRepository
}

func NewLanguageHandler(userLanguageRepository *repository.UserLanguageRepository) *LanguageHandler {
	return &LanguageHandler{
		userLanguageRepository: userLanguageRepository,
	}
}

func (h *LanguageHandler) LanguageCommand(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	keyboardButtons := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Українська", CallbackData: "lang_ukrainian_button"},
				{Text: "English", CallbackData: "lang_english_button"},
			}, {
				{Text: "Polish", CallbackData: "lang_polish_button"},
			},
		},
	}

	_, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        messages.ChangeLanguageMessage,
		ReplyMarkup: keyboardButtons,
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (h *LanguageHandler) CallbackHandler(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	callbackQuery := update.CallbackQuery
	if callbackQuery == nil {
		return
	}

	_, err := b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
		ShowAlert:       false,
	})

	if err != nil {
		log.Printf("Failed to answer callback query: %v", err)
		return
	}

	data := update.CallbackQuery.Data
	var language string
	switch data {
	case "lang_ukrainian_button":
		language = "uk_UA"
	case "lang_english_button":
		language = "en_US"
	case "lang_polish_button":
		language = "pl_PL"
	default:
		log.Printf("Unknown callback data: %s", data)
		return
	}

	err = h.userLanguageRepository.SetUserLanguage(ctx, callbackQuery.From.ID, language)
	if err != nil {
		log.Printf("Failed to set user language: %v", err)
		return
	}

	_, err = b.DeleteMessage(ctx, &tgbot.DeleteMessageParams{
		ChatID:    callbackQuery.Message.Message.Chat.ID,
		MessageID: callbackQuery.Message.Message.ID,
	})
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	}

	_, err = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: callbackQuery.Message.Message.Chat.ID,
		Text:   messages.SelectedLanguageMessage,
	})

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}
