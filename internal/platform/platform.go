package platform

import (
	"regexp"
)

type Platform struct {
	Name     string
	Patterns []*regexp.Regexp
	Format   string
}

var YouTube = &Platform{
	Name:   "youtube",
	Format: "best[ext=mp4][acodec!=none][vcodec!=none][height<=720][filesize<49M]/best[ext=mp4][acodec!=none][vcodec!=none][height<=720][filesize_approx<49M]",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/watch\?.*v=[\w-]+`),
		regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/shorts/[\w-]+`),
	},
}

var TikTok = &Platform{
	Name:   "tiktok",
	Format: "\"best[ext=mp4][acodec!=none][vcodec!=none][filesize<49M]/best[ext=mp4][acodec!=none][vcodec!=none][filesize_approx<49M]",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(www\.)?tiktok\.com/@[\w\.-]+/video/\d+`),
		regexp.MustCompile(`^https?://vt\.tiktok\.com/[\w-]+`),
		regexp.MustCompile(`^https?://vm\.tiktok\.com/[\w-]+`),
	},
}

var Instagram = &Platform{
	Name:   "instagram",
	Format: "best[ext=mp4][acodec!=none][vcodec!=none][filesize<49M]/best[ext=mp4][acodec!=none][vcodec!=none][filesize_approx<49M]",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/(?:reel|p)/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/reels/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/stories/[\w\.-]+/\d+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagr\.am/(?:reel|p)/[\w-]+/?`),
		regexp.MustCompile(`^https?://(?:www\.)?instagram\.com/share/(?:reel|p|story)/[\w-]+/?`),
	},
}

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
