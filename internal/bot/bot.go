package bot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/repository"
	tgbot "github.com/go-telegram/bot"

	"github.com/FortiBrine/ClipHarborBot/internal/handler"
)

type ClipHarborBot struct {
	bot                    *tgbot.Bot
	botConfig              *config.Config
	database               *database.Database
	userLanguageRepository *repository.UserLanguageRepository
}

func New(
	config *config.Config,
	database *database.Database,
	userLanguageRepository *repository.UserLanguageRepository,
) (*ClipHarborBot, error) {

	defaultHandler := handler.NewDefaultHandler(userLanguageRepository)
	languageHandler := handler.NewLanguageHandler(userLanguageRepository)

	options := []tgbot.Option{
		tgbot.WithDefaultHandler(defaultHandler.Default),
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

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"lang",
		tgbot.MatchTypeCommandStartOnly,
		languageHandler.LanguageCommand,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		"lang",
		tgbot.MatchTypePrefix,
		languageHandler.CallbackHandler,
	)

	return &ClipHarborBot{
		bot:                    bot,
		botConfig:              config,
		database:               database,
		userLanguageRepository: userLanguageRepository,
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
