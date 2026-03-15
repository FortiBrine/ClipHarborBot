package downloader

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type MetadataFetcher interface {
	Fetch(ctx context.Context, url string) (*VideoMeta, error)
}

type YTDLPFetcher struct{}

func (f *YTDLPFetcher) Fetch(ctx context.Context, url string) (*VideoMeta, error) {

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"-J",
		"--no-playlist",
		"--no-warnings",
		"--",
		url,
	)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {

		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, fmt.Errorf("yt-dlp error: %s", string(ee.Stderr))
		}

		return nil, err
	}

	var meta VideoMeta

	decoder := json.NewDecoder(io.LimitReader(&buf, 10<<20))
	if err := decoder.Decode(&meta); err != nil {
		return nil, err
	}

	if meta.Duration <= 0 || len(meta.Formats) == 0 {
		return nil, errors.New("invalid metadata")
	}

	return &meta, nil
}

type VideoMeta struct {
	Duration float64  `json:"duration"`
	Formats  []Format `json:"formats"`
}

type Format struct {
	FormatID       string  `json:"format_id"`
	Ext            string  `json:"ext"`
	Filesize       int64   `json:"filesize"`
	FilesizeApprox int64   `json:"filesize_approx"`
	Tbr            float64 `json:"tbr"`
	Height         int     `json:"height"`
	Vcodec         string  `json:"vcodec"`
	Acodec         string  `json:"acodec"`
}
