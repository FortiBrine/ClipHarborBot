package bot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/repository"
	"github.com/FortiBrine/ClipHarborBot/internal/service"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/FortiBrine/ClipHarborBot/internal/handler"
)

type ClipHarborBot struct {
	bot            *tgbot.Bot
	botConfig      *config.Config
	database       *database.Database
	messageService *service.MessageService
}

func New(
	config *config.Config,
	database *database.Database,
	messageService *service.MessageService,
	userLanguageRepository *repository.UserLanguageRepository,
) (*ClipHarborBot, error) {

	defaultHandler := handler.NewDefaultHandler(messageService)
	languageHandler := handler.NewLanguageHandler(userLanguageRepository)
	tiktokHandler := handler.NewTiktokHandler(messageService)
	startHandler := handler.NewStartHandler(messageService)

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
		tiktokHandler.Handle,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"start",
		tgbot.MatchTypeCommandStartOnly,
		startHandler.Handle,
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
		bot:            bot,
		botConfig:      config,
		database:       database,
		messageService: messageService,
	}, nil
}

func (b *ClipHarborBot) Start(ctx context.Context) error {
	if _, err := b.bot.SetWebhook(ctx, &tgbot.SetWebhookParams{
		URL:         b.botConfig.WebhookURL,
		SecretToken: b.botConfig.WebhookSecret,
	}); err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	_, err := b.bot.SetMyCommands(ctx, &tgbot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "start",
				Description: "Start the bot and get a welcome message",
			},
			{
				Command:     "tiktok",
				Description: "Download TikTok videos by providing a link",
			},
			{
				Command:     "lang",
				Description: "Change the bot's language",
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to set bot commands: %w", err)
	}

	b.bot.StartWebhook(ctx)
	return nil
}

func (b *ClipHarborBot) WebhookHandler() http.Handler {
	return b.bot.WebhookHandler()
}
