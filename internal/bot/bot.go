package bot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
	tgbot "github.com/go-telegram/bot"
)

type Bot struct {
	bot    *tgbot.Bot
	logger *slog.Logger
}

func New(
	ctx context.Context, cfg config.Config,
	logger *slog.Logger,
) (b *Bot, err error) {
	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}

	i18nService := i18n.NewService()
	if err = i18nService.LoadTranslations(); err != nil {
		return nil, fmt.Errorf("loading translations: %w", err)
	}

	userRepository := user.NewPostgresRepository(db)
	if err = userRepository.Migrate(); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}
	userService := user.NewService(userRepository, cfg.DefaultLang)

	maxSize := int64(49 * 1024 * 1024)
	fetcher := video.NewYTDLPFetcher()
	downloader, err := video.NewYTDLPDownloader(
		logger,
		cfg.DownloadTimeout,
		maxSize,
	)
	if err != nil {
		return nil, fmt.Errorf("creating video downloader: %w", err)
	}
	formatSelector := video.NewFormatSelector(maxSize)

	options := []tgbot.Option{
		tgbot.WithDefaultHandler(NewDefaultHandler(
			logger,
			userService,
			i18nService,
		)),
	}

	if cfg.WebhookSecret != "" {
		options = append(options,
			tgbot.WithWebhookSecretToken(cfg.WebhookSecret),
		)
	}

	bot, err := tgbot.New(cfg.TelegramToken, options...)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	RegisterRoutes(
		logger, bot,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)

	b = new(Bot{
		bot:    bot,
		logger: logger,
	})

	return
}

func (b *Bot) Start(ctx context.Context, cfg config.Config) error {
	if cfg.WebhookURL != "" && cfg.WebhookSecret != "" {
		if _, err := b.bot.SetWebhook(ctx, new(tgbot.SetWebhookParams{
			URL:         cfg.WebhookURL,
			SecretToken: cfg.WebhookSecret,
		})); err != nil {
			return fmt.Errorf("setting webhook: %w", err)
		}

		mux := http.NewServeMux()

		mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		mux.Handle("/webhook", b.bot.WebhookHandler())

		b.bot.StartWebhook(ctx)

		go func(logger *slog.Logger) {
			if err := http.ListenAndServe(cfg.HttpAddress, mux); err != nil {
				logger.Error("error starting HTTP server", "error", err)
			}
		}(b.logger)
		return nil
	}

	if _, err := b.bot.DeleteWebhook(ctx, new(tgbot.DeleteWebhookParams{
		DropPendingUpdates: true,
	})); err != nil {
		return fmt.Errorf("deleting webhook: %w", err)
	}

	b.bot.Start(ctx)

	return nil
}
