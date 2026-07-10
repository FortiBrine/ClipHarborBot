package video

import (
	"regexp"
)

type Platform struct {
	Name     string
	Patterns []*regexp.Regexp
}

var YouTube = new(Platform{
	Name: "youtube",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/watch\?.*v=[\w-]+`),
		regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/shorts/[\w-]+`),
	},
})

var TikTok = new(Platform{
	Name: "tiktok",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(www\.)?tiktok\.com/@[\w\.-]+/video/\d+`),
		regexp.MustCompile(`^https?://vt\.tiktok\.com/[\w-]+`),
		regexp.MustCompile(`^https?://vm\.tiktok\.com/[\w-]+`),
	},
})

var Instagram = new(Platform{
	Name: "instagram",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/(?:reel|p)/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/reels/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/stories/[\w\.-]+/\d+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagr\.am/(?:reel|p)/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/share/(?:reel|p|story)/[\w-]+/?`),
	},
})

var Platforms = []*Platform{
	YouTube,
	TikTok,
	Instagram,
}

func (p *Platform) IsValidURL(url string) bool {
	for _, pattern := range p.Patterns {
		if pattern.MatchString(url) {
			return true
		}
	}
	return false
}

func DetectPlatform(url string) *Platform {
	for _, platform := range Platforms {
		if platform.IsValidURL(url) {
			return platform
		}
	}
	return nil
}
