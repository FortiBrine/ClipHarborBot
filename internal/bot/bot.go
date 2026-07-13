package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/middleware"
	"github.com/FortiBrine/ClipHarborBot/internal/store"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
	"github.com/dustin/go-humanize"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type Bot struct {
	logger    *slog.Logger
	pool      *pgxpool.Pool
	bh        *th.BotHandler
	transport *transport
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
	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	i18nService := i18n.NewService()
	if err = i18nService.LoadTranslations(); err != nil {
		return nil, fmt.Errorf("loading translations: %w", err)
	}
	if !i18n.IsSupported(cfg.DefaultLang) {
		return nil, fmt.Errorf("unsupported DEFAULT_LANG: %q", cfg.DefaultLang)
	}

	userRepository := user.NewPostgresRepository(pool)
	userService := user.NewService(userRepository, cfg.DefaultLang)

	maxSize := int64(49 * humanize.MiByte)
	fetcher := video.NewYTDLPFetcher()
	downloader, err := video.NewYTDLPDownloader(
		logger,
		cfg.DownloadTimeout,
		maxSize,
	)
	if err != nil {
		return nil, fmt.Errorf("creating video downloader: %w", err)
	}
	if cleanupErr := downloader.CleanupOldFiles(0); cleanupErr != nil {
		logger.Warn("cleaning up old downloads on startup", "error", cleanupErr)
	}
	formatSelector := video.NewFormatSelector(maxSize)

	telegoBot, err := telego.NewBot(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("creating bot: %w", err)
	}

	tr, updates, err := newTransport(ctx, cfg, logger, telegoBot)
	if err != nil {
		return nil, fmt.Errorf("setting up transport: %w", err)
	}
	defer func() {
		if err != nil {
			tr.Stop(ctx)
		}
	}()

	bh, err := th.NewBotHandler(telegoBot, updates, th.WithErrorHandler(
		middleware.NewCustomErrorHandler(logger),
	))
	if err != nil {
		return nil, fmt.Errorf("creating bot handler: %w", err)
	}
	defer func() {
		if err != nil {
			bh.Stop()
		}
	}()
	bh.Use(th.PanicRecoveryHandler(func(recovered any) error {
		return fmt.Errorf("panic in handler: %v", recovered)
	}))
	RegisterRoutes(
		bh,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
		logger,
	)

	b = new(Bot{
		logger:    logger,
		pool:      pool,
		bh:        bh,
		transport: tr,
	})

	return
}

func (b *Bot) Start(ctx context.Context) error {
	b.transport.Start()

	done := make(chan error, 1)
	go func() {
		done <- b.bh.Start()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-done:
		if err != nil {
			return fmt.Errorf("starting bot handler: %w", err)
		}
		return nil
	}
}

func (b *Bot) Close(ctx context.Context) error {
	if err := b.transport.Stop(ctx); err != nil {
		return fmt.Errorf("stopping transport: %w", err)
	}

	if err := b.bh.StopWithContext(ctx); err != nil {
		return fmt.Errorf("stopping bot handler: %w", err)
	}

	b.pool.Close()
	return nil
}
