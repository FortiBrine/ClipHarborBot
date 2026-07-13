package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"golang.org/x/sync/semaphore"
)

var (
	ErrInvalidMetadata = errors.New("invalid metadata")
)

const metadataFetchTimeout = 60 * time.Second

type MetadataFetcher interface {
	Fetch(ctx context.Context, url string) (*Meta, error)
}

type YTDLPFetcher struct {
	limiter *semaphore.Weighted
}

func NewYTDLPFetcher() *YTDLPFetcher {
	return new(YTDLPFetcher{
		limiter: semaphore.NewWeighted(5),
	})
}

func (f *YTDLPFetcher) Fetch(ctx context.Context, url string) (*Meta, error) {
	if err := f.limiter.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("waiting for metadata fetch slot: %w", err)
	}
	defer f.limiter.Release(1)

	ctx, cancel := context.WithTimeout(ctx, metadataFetchTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	c := gocmd.NewCmdOptions(gocmd.Options{
		BeforeExec: []func(cmd *exec.Cmd){
			func(cmd *exec.Cmd) {
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
			},
		},
	}, "yt-dlp",
		"-J",
		"--no-playlist",
		"--no-warnings",
		"--ignore-config",
		"--no-cache-dir",
		"--retries", "3",
		"--fragment-retries", "3",
		"--socket-timeout", "30",
		"--",
		url,
	)
	c.Env = []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"TMPDIR=" + os.Getenv("TMPDIR"),
		"LANG=" + os.Getenv("LANG"),
	}

	statusChan := c.Start()
	var status gocmd.Status
	select {
	case status = <-statusChan:
	case <-ctx.Done():
		c.Stop()
		status = <-statusChan
	}

	if status.Error != nil {
		return nil, fmt.Errorf("running yt-dlp: %w", status.Error)
	}
	if status.Exit != 0 {
		return nil, fmt.Errorf("yt-dlp exited with code %d: %s", status.Exit, stderr.String())
	}

	meta := new(Meta)
	if err := json.Unmarshal(stdout.Bytes(), meta); err != nil {
		return nil, fmt.Errorf("decoding yt-dlp output: %w", err)
	}

	if meta.Duration <= 0 || len(meta.Formats) == 0 {
		return nil, ErrInvalidMetadata
	}

	return meta, nil
}
