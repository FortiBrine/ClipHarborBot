package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/FortiBrine/ClipHarborBot/internal/bot"
	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/repository"
)

func main() {
	botConfig := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.New(botConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}

	userLanguageRepository := repository.NewUserLanguageRepository(db)
	err = userLanguageRepository.Migrate()

	if err != nil {
		log.Fatalf("Failed to migrate user language repository: %v", err)
	}

	clipHarborBot, err := bot.New(botConfig, db, userLanguageRepository)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
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
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	if err = clipHarborBot.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
