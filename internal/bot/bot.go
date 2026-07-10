package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/store"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
	tgbot "github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Bot struct {
	logger *slog.Logger
	bot    *tgbot.Bot
	pool   *pgxpool.Pool
}

func New(
	ctx context.Context, cfg config.Config,
	logger *slog.Logger,
) (b *Bot, err error) {
	pool, err := store.NewPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating store: %w", err)
	}
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	if err = store.Migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("migrating store: %w", err)
	}
	defer func(p *pgxpool.Pool) {
		if err != nil {
			p.Close()
		}
	}(pool)

	i18nService := i18n.NewService()
	if err = i18nService.LoadTranslations(); err != nil {
		return nil, fmt.Errorf("loading translations: %w", err)
	}

	userRepository := user.NewPostgresRepository(pool)
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
		logger: logger,
		bot:    bot,
		pool:   pool,
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

		srv := new(http.Server{
			Addr:    cfg.HttpAddress,
			Handler: mux,
		})

		go func(logger *slog.Logger) {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("error starting HTTP server", "error", err)
			}
		}(b.logger)

		b.bot.StartWebhook(ctx)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HttpShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			b.logger.Error("error shutting down HTTP server", "error", err)
		}

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

func (b *Bot) Close() {
	b.pool.Close()
}
