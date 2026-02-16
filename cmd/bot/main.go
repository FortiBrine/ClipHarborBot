package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/FortiBrine/ClipHarborBot/internal/bot"
	"github.com/FortiBrine/ClipHarborBot/internal/config"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
		return
	}

	go func() {
		if err := http.ListenAndServe(":2000", b.WebhookHandler()); err != nil {
			log.Fatal(err)
		}
	}()

	err = b.Start(ctx)
	if err != nil {
		log.Fatal("Failed to start bot: ", err)
	}
}
