package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/store"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
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

	bh, err := th.NewBotHandler(telegoBot, updates)
	if err != nil {
		return nil, fmt.Errorf("creating bot handler: %w", err)
	}
	defer func() {
		if err != nil {
			bh.Stop()
		}
	}()
	RegisterRoutes(
		bh,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)

	b = new(Bot{
		logger:    logger,
		pool:      pool,
		bh:        bh,
		transport: tr,
	})

	return
}

func (b *Bot) Start() error {
	b.bh.Start()
	b.transport.Start()
	return nil
}

func (b *Bot) Close(cfg config.Config) {
	b.pool.Close()
	b.bh.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HttpShutdownTimeout)
	defer cancel()

	if err := b.transport.Stop(shutdownCtx); err != nil {
		b.logger.Error("error shutting down HTTP server", "error", err)
	}
}
