package bot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	tgbot "github.com/go-telegram/bot"

	"github.com/FortiBrine/ClipHarborBot/internal/handler"
)

type ClipHarborBot struct {
	bot       *tgbot.Bot
	botConfig *config.Config
}

func New(config *config.Config) (*ClipHarborBot, error) {
	options := []tgbot.Option{
		tgbot.WithDefaultHandler(handler.Default),
		tgbot.WithWebhookSecretToken(config.WebhookSecret),
	}

	bot, err := tgbot.New(config.Token, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"tiktok",
		tgbot.MatchTypeCommandStartOnly,
		handler.Tiktok,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"start",
		tgbot.MatchTypeCommandStartOnly,
		handler.StartHandler,
	)

	return &ClipHarborBot{
		bot:       bot,
		botConfig: config,
	}, nil
}

func (clipHarborBot *ClipHarborBot) Start(ctx context.Context) error {
	if _, err := clipHarborBot.bot.SetWebhook(ctx, &tgbot.SetWebhookParams{
		URL:         clipHarborBot.botConfig.WebhookURL,
		SecretToken: clipHarborBot.botConfig.WebhookSecret,
	}); err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	clipHarborBot.bot.StartWebhook(ctx)
	return nil
}

func (clipHarborBot *ClipHarborBot) WebhookHandler() http.Handler {
	return clipHarborBot.bot.WebhookHandler()
}
