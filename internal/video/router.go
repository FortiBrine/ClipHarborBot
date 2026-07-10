package video

import (
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	tgbot "github.com/go-telegram/bot"
)

func RegisterRoutes(
	logger *slog.Logger,
	bot *tgbot.Bot,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher MetadataFetcher,
	downloader Downloader,
	formatSelector *FormatSelector,
) {
	handler := NewHandler(
		logger,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)
	tiktokHandler := handler.Handle(TikTok, "tiktok_help")
	youtubeHandler := handler.Handle(YouTube, "youtube_help")
	instagramHandler := handler.Handle(Instagram, "instagram_help")

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"tiktok",
		tgbot.MatchTypeCommandStartOnly,
		tiktokHandler,
	)

	for _, pattern := range TikTok.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			tiktokHandler,
		)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"youtube",
		tgbot.MatchTypeCommandStartOnly,
		youtubeHandler,
	)

	for _, pattern := range YouTube.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			youtubeHandler,
		)
	}

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"instagram",
		tgbot.MatchTypeCommandStartOnly,
		instagramHandler,
	)

	for _, pattern := range Instagram.Patterns {
		bot.RegisterHandlerRegexp(
			tgbot.HandlerTypeMessageText,
			pattern,
			instagramHandler,
		)
	}

}
