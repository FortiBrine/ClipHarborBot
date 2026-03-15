package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FortiBrine/ClipHarborBot/internal/bot"
	"github.com/FortiBrine/ClipHarborBot/internal/config"
	"github.com/FortiBrine/ClipHarborBot/internal/database"
	"github.com/FortiBrine/ClipHarborBot/internal/downloader"
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
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

	languageRepository := user.NewPostgresLanguageRepository(db)
	err = languageRepository.Migrate()

	if err != nil {
		log.Printf("Failed to migrate user language repository: %v", err)
	}

	downloaderUtil, err := downloader.NewDownloader()
	if err != nil {
		log.Printf("Failed to initialize downloader: %v", err)
		return
	}

	downloaderUtil.StartCleanupWorker(10*time.Minute, 1*time.Hour)
	formatSelector := downloader.NewFormatSelector(49 * 1024 * 1024)

	tr, err := i18n.New("uk_UA")

	if err != nil {
		log.Printf("Failed to initialize i18n: %v", err)
		return
	}

	userMessagesService := user.NewService(languageRepository, tr)

	clipHarborBot, err := bot.New(
		botConfig,
		userMessagesService,
		downloaderUtil,
		formatSelector,
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
