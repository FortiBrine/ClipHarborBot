package service

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

type PlatformConfig struct {
	Name     string
	Patterns []*regexp.Regexp
	Format   string
}

type Downloader struct {
	tempDir string
	config  PlatformConfig
	timeout time.Duration
	limiter chan struct{}
}

func NewDownloader(config PlatformConfig) (*Downloader, error) {
	tempDir := filepath.Join(os.TempDir(), "clipharbor")

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}

	return &Downloader{
		tempDir: tempDir,
		config:  config,
		timeout: 5 * time.Minute,
		limiter: make(chan struct{}, 5),
	}, nil
}

func (d *Downloader) IsValidURL(url string) bool {
	for _, pattern := range d.config.Patterns {
		if pattern.MatchString(url) {
			return true
		}
	}
	return false
}

func (d *Downloader) DownloadVideo(ctx context.Context, url string) (string, error) {
	if !d.IsValidURL(url) {
		return "", fmt.Errorf("invalid %s url", d.config.Name)
	}

	d.limiter <- struct{}{}
	defer func() { <-d.limiter }()

	output := filepath.Join(
		d.tempDir,
		fmt.Sprintf("%s_%d_%s.mp4", d.config.Name, time.Now().Unix(), uuid.NewString()),
	)

	cmdCtx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"--no-playlist",
		"--no-warnings",
		"--no-progress",
		"--restrict-filenames",
		"-f", d.config.Format,
		"--merge-output-format", "mp4",
		"-o", output,
		"--max-filesize", "100M",
		"--socket-timeout", "30",
		url,
	}

	cmd := exec.CommandContext(cmdCtx, "yt-dlp", args...)

	cmd.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"TMPDIR=" + os.Getenv("TMPDIR"),
		"LANG=" + os.Getenv("LANG"),
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("download failed: %w (%s)", err, string(out))
	}

	if _, err := os.Stat(output); err != nil {
		return "", fmt.Errorf("downloaded file not found: %w", err)
	}

	return output, nil
}

func (d *Downloader) CleanupFile(filePath string) error {

	if filePath == "" {
		return nil
	}

	clean, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	temp, err := filepath.Abs(d.tempDir)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(clean, temp+string(os.PathSeparator)) {
		return fmt.Errorf("invalid path")
	}

	return os.Remove(clean)
}

func (d *Downloader) CleanupOldFiles(olderThan time.Duration) error {

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
			path := filepath.Join(d.tempDir, entry.Name())
			_ = os.Remove(path)
		}
	}

	return nil
}

func (d *Downloader) StartCleanupWorker(interval time.Duration, olderThan time.Duration) {
	go func() {
		_ = d.CleanupOldFiles(olderThan)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			err := d.CleanupOldFiles(olderThan)
			if err != nil {
				log.Printf("cleanup error: %v", err)
			}
		}
	}()
}
