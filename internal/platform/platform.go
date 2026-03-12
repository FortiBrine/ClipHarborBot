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
	Format: "best[height<=720][ext=mp4]/best[height<=720]",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/watch\?.*v=[\w-]+`),
		regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/shorts/[\w-]+`),
	},
}

var TikTok = &Platform{
	Name:   "tiktok",
	Format: "best[ext=mp4]",
	Patterns: []*regexp.Regexp{
		regexp.MustCompile(`^https?://(www\.)?tiktok\.com/@[\w\.-]+/video/\d+`),
		regexp.MustCompile(`^https?://vt\.tiktok\.com/[\w-]+`),
		regexp.MustCompile(`^https?://vm\.tiktok\.com/[\w-]+`),
	},
}

var Platforms = []*Platform{
	YouTube,
	TikTok,
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
