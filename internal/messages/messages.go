package messages

type Lang string

const (
	ChangeLanguageMessage   = "Виберіть мову / Select language / Wybierz język"
	SelectedLanguageMessage = "Мова змінена! / Language changed! / Zmieniono język!"
)

var Messages = map[string]map[string]string{
	"uk_UA": {
		"unknown_command":        "❌ Невідома команда",
		"start_command":          "Привіт! Я ClipHarborBot, і я допоможу тобі завантажувати відео з TikTok та YouTube. Просто надішли мені посилання на відео, і я зроблю все інше! 🚀",
		"tiktok_help":            "ℹ️ Використання: /tiktok <посилання на TikTok відео>\n\nПриклад:\n/tiktok https://www.tiktok.com/@user/video/123456789",
		"tiktok_downloading":     "⏳ Завантажую відео...",
		"tiktok_uploading":       "⬆️ Відправляю відео...",
		"tiktok_invalid_url":     "❌ Невірне посилання на TikTok. Будь ласка, надішліть коректне посилання.",
		"tiktok_download_error":  "❌ Помилка при завантаженні відео. Спробуйте ще раз пізніше.",
		"tiktok_size_error":      "❌ Відео занадто велике (більше 100MB). Telegram не підтримує такі великі файли.",
		"tiktok_success":         "✅ Відео успішно завантажено!",
		"youtube_help":           "ℹ️ Використання: /youtube <посилання на YouTube відео>\n\nПриклад:\n/youtube https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"youtube_downloading":    "⏳ Завантажую відео з YouTube...",
		"youtube_uploading":      "⬆️ Відправляю відео...",
		"youtube_invalid_url":    "❌ Невірне посилання на YouTube. Будь ласка, надішліть коректне посилання.",
		"youtube_download_error": "❌ Помилка при завантаженні відео. Спробуйте ще раз пізніше.",
		"youtube_size_error":     "❌ Відео занадто велике (більше 100MB). Telegram не підтримує такі великі файли.",
		"youtube_success":        "✅ Відео успішно завантажено!",
	},
	"en_US": {
		"unknown_command":        "❌ Unknown command",
		"start_command":          "Hello! I'm ClipHarborBot, and I'm here to help you download videos from TikTok and YouTube. Just send me a link to the video, and I'll take care of the rest! 🚀",
		"tiktok_help":            "ℹ️ Usage: /tiktok <TikTok video URL>\n\nExample:\n/tiktok https://www.tiktok.com/@user/video/123456789",
		"tiktok_downloading":     "⏳ Downloading video...",
		"tiktok_uploading":       "⬆️ Uploading video...",
		"tiktok_invalid_url":     "❌ Invalid TikTok URL. Please send a valid link.",
		"tiktok_download_error":  "❌ Error downloading video. Please try again later.",
		"tiktok_size_error":      "❌ Video is too large (over 100MB). Telegram does not support such large files.",
		"tiktok_success":         "✅ Video downloaded successfully!",
		"youtube_help":           "ℹ️ Usage: /youtube <YouTube video URL>\n\nExample:\n/youtube https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"youtube_downloading":    "⏳ Downloading video from YouTube...",
		"youtube_uploading":      "⬆️ Uploading video...",
		"youtube_invalid_url":    "❌ Invalid YouTube URL. Please send a valid link.",
		"youtube_download_error": "❌ Error downloading video. Please try again later.",
		"youtube_size_error":     "❌ Video is too large (over 100MB). Telegram does not support such large files.",
		"youtube_success":        "✅ Video downloaded successfully!",
	},
}
