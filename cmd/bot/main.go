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
	"github.com/FortiBrine/ClipHarborBot/internal/messages"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
)

func main() {
	botConfig := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	db, err := database.New(botConfig)

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	userLanguageRepository := user.NewUserLanguageRepository(db)
	err = userLanguageRepository.Migrate()

	if err != nil {
		log.Printf("Failed to migrate user language repository: %v", err)
	}

	messageService := messages.NewMessageService(userLanguageRepository)

	clipHarborBot, err := bot.New(
		botConfig,
		db,
		messageService,
		userLanguageRepository,
	)

	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	if botConfig.Mode == config.ModeWebhook {
		mux.Handle("/webhook", clipHarborBot.WebhookHandler())
	}

	go func() {
		log.Println("HTTP server listening on :2000")

		if err := http.ListenAndServe(":2000", mux); err != nil {
			log.Printf("Failed to start HTTP server: %v", err)
		}
	}()

	if err = clipHarborBot.Start(ctx); err != nil {
		log.Printf("Failed to start bot: %v", err)
	}
}
