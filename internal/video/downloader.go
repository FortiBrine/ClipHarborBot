package video

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
)

var (
	ErrFileTooLarge = errors.New("file too large for telegram")
)

type Downloader interface {
	DownloadVideo(ctx context.Context, opts DownloadOptions) (output string, err error)
	CleanupFile(filePath string) error
	CleanupOldFiles(olderThan time.Duration) error
}

type YTDLPDownloader struct {
	logger  *slog.Logger
	tempDir string
	timeout time.Duration
	limiter *semaphore.Weighted
	maxSize int64
}

type DownloadOptions struct {
	URL    string
	Format string
	Prefix string
}

func NewYTDLPDownloader(
	logger *slog.Logger,
	timeout time.Duration,
	maxSize int64,
) (*YTDLPDownloader, error) {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return nil, fmt.Errorf("yt-dlp not installed: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "clipharborbot-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}

	return new(YTDLPDownloader{
		logger:  logger,
		tempDir: tempDir,
		timeout: timeout,
		limiter: semaphore.NewWeighted(5),
		maxSize: maxSize,
	}), nil
}

func (d *YTDLPDownloader) DownloadVideo(
	ctx context.Context,
	opts DownloadOptions,
) (output string, err error) {
	if err = d.limiter.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("waiting for download slot: %w", err)
	}
	defer d.limiter.Release(1)

	tempFile, err := os.CreateTemp(d.tempDir, opts.Prefix+"-*.mp4")
	if err != nil {
		return "", fmt.Errorf("reserving temp file: %w", err)
	}
	output = tempFile.Name()
	if err = tempFile.Close(); err != nil {
		return "", fmt.Errorf("closing temp file: %w", err)
	}

	defer func() {
		if err != nil {
			if rmErr := os.Remove(output); rmErr != nil && !os.IsNotExist(rmErr) {
				d.logger.Error("removing incomplete download", "file", output, "error", rmErr)
			}
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
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

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

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

		return "", fmt.Errorf("downloading video: %w (%s)", err, string(out))
	}

	info, err := os.Stat(output)
	if err != nil {
		return "", fmt.Errorf("statting output file: %w", err)
	}

	if info.Size() > d.maxSize {
		return "", ErrFileTooLarge
	}

	return output, nil
}

func (d *YTDLPDownloader) CleanupFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	return os.Remove(filePath)
}

func (d *YTDLPDownloader) CleanupOldFiles(olderThan time.Duration) error {
	entries, err := os.ReadDir(d.tempDir)
	if err != nil {
		return fmt.Errorf("reading temp dir: %w", err)
	}

	now := time.Now()

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, err := e.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > olderThan {
			path := filepath.Join(d.tempDir, e.Name())

			if err = os.Remove(path); err != nil {
				return fmt.Errorf("removing temp file: %w", err)
			}
		}
	}

	return nil
}
