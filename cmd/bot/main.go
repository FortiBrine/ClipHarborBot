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
	botConfig := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	clipHarborBot, err := bot.New(botConfig)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/webhook", clipHarborBot.WebhookHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	go func() {
		if err := http.ListenAndServe(":2000", mux); err != nil {
			log.Fatal(err)
		}
	}()

	if err = clipHarborBot.Start(ctx); err != nil {
		log.Fatal("Failed to start bot: ", err)
	}
}
