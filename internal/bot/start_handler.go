package bot

import (
	"context"
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func NewStartHandler(
	logger *slog.Logger,
	userService *user.Service,
	i18nService *i18n.Service,
) tgbot.HandlerFunc {
	return func(ctx context.Context, bot *tgbot.Bot, update *models.Update) {
		lang, err := userService.GetLanguage(ctx, update.Message.From.ID)
		if err != nil {
			logger.Error("getting user language", "error", err)
			return
		}

		localizer := i18nService.FromLanguage(lang)
		if _, err = bot.SendMessage(ctx, new(tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   i18n.T(localizer, "start_command"),
		})); err != nil {
			logger.Error("sending start command", "error", err)
		}
	}
}
