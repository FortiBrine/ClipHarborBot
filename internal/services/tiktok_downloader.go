package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TikTokDownloader struct {
	tempDir string
}

func NewTikTokDownloader() (*TikTokDownloader, error) {
	tempDir := filepath.Join(os.TempDir(), "clipharbor")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &TikTokDownloader{
		tempDir: tempDir,
	}, nil
}

func (d *TikTokDownloader) IsValidTikTokURL(url string) bool {
	patterns := []string{
		`^https?://(?:www\.)?tiktok\.com/@[\w.-]+/video/\d+`,
		`^https?://(?:vm|vt)\.tiktok\.com/[\w]+`,
		`^https?://(?:www\.)?tiktok\.com/t/[\w]+`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}

	return false
}

func (d *TikTokDownloader) DownloadVideo(ctx context.Context, url string) (string, error) {
	if !d.IsValidTikTokURL(url) {
		return "", fmt.Errorf("invalid TikTok URL")
	}

	outputTemplate := filepath.Join(
		d.tempDir,
		fmt.Sprintf("tiktok_%s.mp4", uuid.NewString()),
	)

	cmdCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	args := []string{
		"--no-playlist",
		"--no-warnings",
		"--no-progress",
		"-f", "best[ext=mp4]",
		"-o", outputTemplate,
		"--max-filesize", "100M",
		"--socket-timeout", "30",
		url,
	}

	cmd := exec.CommandContext(cmdCtx, "yt-dlp", args...)

	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w (output: %s)", err, string(output))
	}

	if _, err := os.Stat(outputTemplate); os.IsNotExist(err) {
		return "", fmt.Errorf("downloaded file not found")
	}

	return outputTemplate, nil
}

func (d *TikTokDownloader) CleanupFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	cleanPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	tempDir, _ := filepath.Abs(d.tempDir)
	if !strings.HasPrefix(cleanPath, tempDir+string(os.PathSeparator)) {
		return fmt.Errorf("security error: file path outside temp directory")
	}

	return os.Remove(filePath)
}

func (d *TikTokDownloader) CleanupOldFiles(olderThan time.Duration) error {
	entries, err := os.ReadDir(d.tempDir)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > olderThan {
			filePath := filepath.Join(d.tempDir, entry.Name())

			if err := os.Remove(filePath); err != nil {
				log.Printf("cleanup failed: %s: %v", filePath, err)
			}
		}
	}

	return nil
}
