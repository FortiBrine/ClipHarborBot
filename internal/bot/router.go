package bot

import (
	"log/slog"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
	tgbot "github.com/go-telegram/bot"
)

func RegisterRoutes(
	logger *slog.Logger,
	bot *tgbot.Bot,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher video.MetadataFetcher,
	downloader video.Downloader,
	formatSelector *video.FormatSelector,
) {
	video.RegisterRoutes(
		logger, bot,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"start",
		tgbot.MatchTypeCommandStartOnly,
		NewStartHandler(logger, userService, i18nService),
	)

	userHandler := user.NewHandler(logger, userService, i18nService)
	bot.RegisterHandler(
		tgbot.HandlerTypeMessageText,
		"lang",
		tgbot.MatchTypeCommandStartOnly,
		userHandler.SetLanguage,
	)

	bot.RegisterHandler(
		tgbot.HandlerTypeCallbackQueryData,
		"lang",
		tgbot.MatchTypePrefix,
		userHandler.CallbackHandler,
	)
}
