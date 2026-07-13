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
		regexp.MustCompile(`^https?://(?:www\.|m\.)?youtube\.com/watch\?.*v=[\w-]+`),
		regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.|m\.)?youtube\.com/shorts/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.|m\.)?youtube\.com/live/[\w-]+`),
	},
})

var TikTok = new(Platform{
	Name: "tiktok",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(www\.)?tiktok\.com/@[\w\.-]+/video/\d+`),
		regexp.MustCompile(`^https?://vt\.tiktok\.com/[\w-]+`),
		regexp.MustCompile(`^https?://vm\.tiktok\.com/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.)?tiktok\.com/t/[\w-]+`),
	},
})

var Instagram = new(Platform{
	Name: "instagram",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/(?:reel|p|tv)/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/reels/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/stories/[\w\.-]+/\d+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagr\.am/(?:reel|p|tv)/[\w-]+/?`),
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
