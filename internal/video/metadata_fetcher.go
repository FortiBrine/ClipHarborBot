package video

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var (
	ErrInvalidMetadata = errors.New("invalid metadata")
)

type MetadataFetcher interface {
	Fetch(ctx context.Context, url string) (*Meta, error)
}

type YTDLPFetcher struct{}

func NewYTDLPFetcher() *YTDLPFetcher {
	return new(YTDLPFetcher{})
}

func (f *YTDLPFetcher) Fetch(ctx context.Context, url string) (*Meta, error) {
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
		if ee, ok := errors.AsType[*exec.ExitError](err); ok {
			return nil, fmt.Errorf("yt-dlp exited with code %d: %s", ee.ExitCode(), string(ee.Stderr))
		}

		return nil, fmt.Errorf("running yt-dlp: %w", err)
	}

	var meta Meta

	decoder := json.NewDecoder(io.LimitReader(&buf, 10<<20))
	if err := decoder.Decode(&meta); err != nil {
		return nil, fmt.Errorf("decoding yt-dlp error: %w", err)
	}

	if meta.Duration <= 0 || len(meta.Formats) == 0 {
		return nil, ErrInvalidMetadata
	}

	return &meta, nil
}
