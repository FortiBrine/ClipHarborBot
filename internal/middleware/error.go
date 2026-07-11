package middleware

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func NewCustomErrorHandler(logger *slog.Logger) th.ErrorHandler {
	return func(
		ctx *th.Context,
		update telego.Update,
		err error,
	) {
		if err == nil {
			return
		}

		logger.Error("error handling update",
			"chat_id", update.Message.Chat.ID,
			"from_id", update.Message.From.ID,
			"message_id", update.Message.MessageID,
			"error", err,
		)
	}
}
