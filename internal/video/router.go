package video

import (
	"fmt"

	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	th "github.com/mymmrac/telego/telegohandler"
)

func RegisterRoutes(
	bh *th.BotHandler,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher MetadataFetcher,
	downloader Downloader,
	formatSelector *FormatSelector,
) {
	handler := NewHandler(
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)

	for _, platform := range Platforms {
		platformHandler := handler.Handle(platform, fmt.Sprintf("%s_help", platform.Name))
		bh.HandleMessage(platformHandler, th.CommandPrefix(platform.Name))
		for _, pattern := range platform.Patterns {
			bh.HandleMessage(platformHandler, th.TextMatches(pattern))
		}
	}

}
