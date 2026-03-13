package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrFileTooLarge  = errors.New("file too large for telegram")
	ErrInvalidFormat = errors.New("requested format not available")
)

type DownloadOptions struct {
	URL    string
	Format string
	Prefix string
}

type Downloader struct {
	tempDir string
	timeout time.Duration
	limiter chan struct{}
}

func NewDownloader() (*Downloader, error) {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return nil, fmt.Errorf("yt-dlp not installed: %w", err)
	}

	tempDir := filepath.Join(os.TempDir(), "clipharbor")

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}

	return &Downloader{
		tempDir: tempDir,
		timeout: 5 * time.Minute,
		limiter: make(chan struct{}, 5),
	}, nil
}

func (d *Downloader) DownloadVideo(
	ctx context.Context,
	opts DownloadOptions,
) (string, error) {
	d.limiter <- struct{}{}
	defer func() { <-d.limiter }()

	output := filepath.Join(
		d.tempDir,
		fmt.Sprintf("%s_%d_%s.mp4", opts.Prefix, time.Now().Unix(), uuid.NewString()),
	)

	cmdCtx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"--no-playlist",
		"--no-warnings",
		"--no-progress",
		"--restrict-filenames",
		"-f", opts.Format,
		"--merge-output-format", "mp4",
		"-o", output,
		"--max-filesize", "49M",
		"--socket-timeout", "30",
		opts.URL,
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
		outStr := string(out)
		if strings.Contains(outStr, "Requested format is not available") {
			return "", ErrInvalidFormat
		}

		return "", fmt.Errorf("download failed: %w (%s)", err, string(out))
	}

	info, err := os.Stat(output)
	if err != nil {
		return "", fmt.Errorf("downloaded file not found: %w", err)
	}

	if info.Size() > 49*1024*1024 {
		err = os.Remove(output)

		if err != nil {
			log.Printf("failed to remove oversized file: %v", err)
		}

		return "", ErrFileTooLarge
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
			err = os.Remove(path)

			if err != nil {
				log.Printf("failed to remove old file: %v", err)
			}
		}
	}

	return nil
}

func (d *Downloader) StartCleanupWorker(interval time.Duration, olderThan time.Duration) {
	go func() {
		err := d.CleanupOldFiles(olderThan)

		if err != nil {
			log.Printf("initial cleanup error: %v", err)
		}

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
