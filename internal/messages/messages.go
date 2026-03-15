package messages

type Lang string

const (
	ChangeLanguageMessage   = "Виберіть мову / Select language / Wybierz język"
	SelectedLanguageMessage = "Мова змінена! / Language changed! / Zmieniono język!"
)

var Messages = map[string]map[string]string{
	"uk_UA": {
		"unknown_command":        "❌ Невідома команда",
		"start_command":          "Привіт! Я ClipHarborBot, і я допоможу тобі завантажувати відео з TikTok, YouTube та Instagram. Просто надішли мені посилання на відео, і я зроблю все інше! 🚀",
		"tiktok_help":            "ℹ️ Використання: /tiktok <посилання на TikTok відео>\n\nПриклад:\n/tiktok https://www.tiktok.com/@user/video/123456789",
		"youtube_help":           "ℹ️ Використання: /youtube <посилання на YouTube відео>\n\nПриклад:\n/youtube https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"instagram_help":         "ℹ️ Використання: /instagram <посилання на Instagram відео>\n\nПриклад:\n/instagram https://www.instagram.com/reel/CrqQfC0p8kK/",
		"video_expected_size":    "📦 Розмір відео буде: %s",
		"video_downloading":      "⏳ Завантажую відео...",
		"video_uploading":        "⬆️ Відправляю відео...",
		"invalid_video_url":      "❌ Невірне посилання. Будь ласка, надішліть коректне посилання.",
		"video_download_error":   "❌ Помилка при завантаженні відео. Спробуйте ще раз пізніше.",
		"video_upload_error":     "❌ Помилка при відправці відео. Спробуйте ще раз пізніше.",
		"video_size_error":       "❌ Відео занадто велике (більше 50MB). Telegram не підтримує такі великі файли.",
		"video_format_error":     "❌ Непідтримуваний формат відео. Спробуйте інший формат або інше посилання.",
		"video_download_success": "✅ Відео успішно завантажено!",
	},
	"en_US": {
		"unknown_command":        "❌ Unknown command",
		"start_command":          "Hello! I'm ClipHarborBot, and I'm here to help you download videos from TikTok, YouTube, and Instagram. Just send me a link to the video, and I'll take care of the rest! 🚀",
		"tiktok_help":            "ℹ️ Usage: /tiktok <TikTok video URL>\n\nExample:\n/tiktok https://www.tiktok.com/@user/video/123456789",
		"youtube_help":           "ℹ️ Usage: /youtube <YouTube video URL>\n\nExample:\n/youtube https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"instagram_help":         "ℹ️ Usage: /instagram <Instagram video URL>\n\nExample:\n/instagram https://www.instagram.com/reel/CrqQfC0p8kK/",
		"video_expected_size":    "📦 Expected video size: %s",
		"video_downloading":      "⏳ Downloading video...",
		"video_uploading":        "⬆️ Uploading video...",
		"invalid_video_url":      "❌ Invalid URL. Please send a valid link.",
		"video_download_error":   "❌ Error downloading video. Please try again later.",
		"video_upload_error":     "❌ Error uploading video. Please try again later.",
		"video_size_error":       "❌ Video is too large (over 50MB). Telegram does not support such large files.",
		"video_format_error":     "❌ Unsupported video format. Try a different format or a different link.",
		"video_download_success": "✅ Video downloaded successfully!",
	},
}
