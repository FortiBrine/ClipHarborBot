package bot

import (
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func NewStartHandler(
	userService *user.Service,
	i18nService *i18n.Service,
) th.MessageHandler {
	return func(ctx *th.Context, message telego.Message) error {
		userID := message.From.ID
		lang, err := userService.GetLanguage(ctx, userID)
		if err != nil {
			return fmt.Errorf("getting user language: %w", err)
		}

		localizer := i18nService.FromLanguage(lang)
		if _, err = ctx.Bot().SendMessage(ctx, tu.Message(
			message.Chat.ChatID(),
			i18n.T(localizer, msgStartCommand),
		)); err != nil {
			return fmt.Errorf("sending start_command: %w", err)
		}
		return nil
	}
}
