package bot

import (
	"github.com/FortiBrine/ClipHarborBot/internal/i18n"
	"github.com/FortiBrine/ClipHarborBot/internal/user"
	"github.com/FortiBrine/ClipHarborBot/internal/video"
	th "github.com/mymmrac/telego/telegohandler"
)

func RegisterRoutes(
	bh *th.BotHandler,
	userService *user.Service,
	i18nService *i18n.Service,
	fetcher video.MetadataFetcher,
	downloader video.Downloader,
	formatSelector *video.FormatSelector,
) {
	video.RegisterRoutes(
		bh,
		userService,
		i18nService,
		fetcher,
		downloader,
		formatSelector,
	)

	userHandler := user.NewHandler(userService, i18nService)
	bh.HandleMessage(NewStartHandler(userService, i18nService), th.CommandEqual("start"))
	bh.HandleMessage(userHandler.SetLanguage, th.CommandEqual("lang"))
	bh.HandleCallbackQuery(userHandler.CallbackHandler, th.CallbackDataPrefix("lang"))
	bh.HandleMessage(NewDefaultHandler(userService, i18nService), th.AnyCommand())
}
