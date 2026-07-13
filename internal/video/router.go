package video

import (
	"log/slog"
	"time"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	th "github.com/mymmrac/telego/telegohandler"
)

const userRateLimitInterval = 5 * time.Second

func RegisterRoutes(
	bh *th.BotHandler,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher MetadataFetcher,
	downloader Downloader,
	formatSelector *FormatSelector,
	logger *slog.Logger,
) {
	handler := NewHandler(
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
		NewRateLimiter(userRateLimitInterval),
		logger,
	)

	for _, platform := range Platforms {
		for _, pattern := range platform.Patterns {
			bh.HandleMessage(handler.Handle(platform, pattern), th.TextMatches(pattern))
		}
	}
}
