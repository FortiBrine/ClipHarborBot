package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/downloader"
	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/platform"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/FortiBrine/ClipHarborBot/internal/handler"
)

type ClipHarborBot struct {
	bot            *tgbot.Bot
	botConfig      *config.Config
	database       *database.Database
	messageService *messages.MessageService
}

func New(
	cfg *config.Config,
	database *database.Database,
	messageService *messages.MessageService,
	userLanguageRepository *user.LanguageRepository,
) (*ClipHarborBot, error) {

	languageHandler := handler.NewLanguageHandler(userLanguageRepository)
	startHandler := handler.NewStartHandler(messageService)
	defaultHandler := handler.NewDefaultHandler(messageService)

	downloaderUtil, err := downloader.NewDownloader()
	if err != nil {
		log.Printf("Failed to initialize downloader: %v", err)
	}

	downloaderUtil.StartCleanupWorker(10*time.Minute, 1*time.Hour)

	formatSelector := downloader.NewFormatSelector(49 * 1024 * 1024)

	tiktokHandler := downloader.NewVideoHandler(messageService, downloaderUtil, formatSelector, platform.TikTok, "tiktok_help")
	youtubeHandler := downloader.NewVideoHandler(messageService, downloaderUtil, formatSelector, platform.YouTube, "youtube_help")
	instagramHandler := downloader.NewVideoHandler(messageService, downloaderUtil, formatSelector, platform.Instagram, "instagram_help")

	options := []tgbot.Option{
		tgbot.WithDefaultHandler(defaultHandler.Default),
	}

	if cfg.Mode == config.ModeWebhook {
		options = append(options,
			tgbot.WithWebhookSecretToken(cfg.WebhookSecret),
		)
	}

	bot, err := tgbot.New(cfg.Token, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"tiktok",
		tgbot.MatchTypeCommandStartOnly,
		tiktokHandler.Handle,
	)

	for _, pattern := range platform.TikTok.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			tiktokHandler.Handle,
		)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"youtube",
		tgbot.MatchTypeCommandStartOnly,
		youtubeHandler.Handle,
	)

	for _, pattern := range platform.YouTube.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			youtubeHandler.Handle,
		)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"instagram",
		tgbot.MatchTypeCommandStartOnly,
		instagramHandler.Handle,
	)

	for _, pattern := range platform.Instagram.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			instagramHandler.Handle,
		)
	}

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
		botConfig:      cfg,
		database:       database,
		messageService: messageService,
	}, nil
}

func (b *ClipHarborBot) Start(ctx context.Context) error {

	switch b.botConfig.Mode {
	case config.ModeWebhook:
		log.Printf("Starting bot in webhook mode with URL: %s", b.botConfig.WebhookURL)

		_, err := b.bot.SetWebhook(ctx, &tgbot.SetWebhookParams{
			URL:         b.botConfig.WebhookURL,
			SecretToken: b.botConfig.WebhookSecret,
		})

		if err != nil {
			return fmt.Errorf("failed to set webhook: %w", err)
		}

	case config.ModePolling:
		log.Printf("Starting bot in polling mode")

		_, err := b.bot.DeleteWebhook(ctx, &tgbot.DeleteWebhookParams{
			DropPendingUpdates: true,
		})

		if err != nil {
			fmt.Printf("failed to delete webhook: %v", err)
		}
	default:
		return fmt.Errorf("invalid bot mode: %s", b.botConfig.Mode)
	}

	_, err := b.bot.SetMyCommands(ctx, &tgbot.SetMyCommandsParams{
		Scope: &models.BotCommandScopeDefault{},
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
				Command:     "youtube",
				Description: "Download YouTube videos by providing a link",
			},
			{
				Command:     "instagram",
				Description: "Download Instagram videos by providing a link",
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

	if b.botConfig.Mode == config.ModePolling {
		b.bot.Start(ctx)
		return nil
	}

	b.bot.StartWebhook(ctx)
	return nil
}

func (b *ClipHarborBot) WebhookHandler() http.Handler {
	return b.bot.WebhookHandler()
}
