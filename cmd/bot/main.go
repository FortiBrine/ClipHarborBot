package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/FortiBrine/ClipHarborBot/internal/bot"
	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}

	l := logger.New(cfg.Environment)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(ctx, cfg, l)
	if err != nil {
		l.Error("error creating bot", "error", err)
		os.Exit(1)
	}
	defer func(b *bot.Bot) {

	}(b)

	if err = b.Start(ctx, cfg); err != nil {
		l.Error("error starting bot", "error", err)
	}
}
