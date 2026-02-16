package bot

import (
	"context"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	tgbot "github.com/go-telegram/bot"

	"github.com/FortiBrine/ClipHarborBot/internal/handlers"
)

type Bot struct {
	bot *tgbot.Bot
	cfg *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(handlers.Default),
		tgbot.WithWebhookSecretToken(cfg.WebhookSecret),
	}

	b, err := tgbot.New(cfg.Token, opts...)
	if err != nil {
		return nil, err
	}

	b.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"tiktok",
		tgbot.MatchTypeCommandStartOnly,
		handlers.Tiktok,
	)

	b.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"start",
		tgbot.MatchTypeCommandStartOnly,
		handlers.StartHandler,
	)

	return &Bot{
		bot: b,
		cfg: cfg,
	}, err
}

func (b *Bot) Start(ctx context.Context) error {
	if _, err := b.bot.SetWebhook(ctx, &tgbot.SetWebhookParams{
		URL:         b.cfg.WebhookURL,
		SecretToken: b.cfg.WebhookSecret,
	}); err != nil {
		return err
	}

	b.bot.StartWebhook(ctx)
	return nil
}

func (b *Bot) WebhookHandler() http.Handler {
	return b.bot.WebhookHandler()
}
