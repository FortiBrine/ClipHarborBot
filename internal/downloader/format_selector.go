package downloader

import (
	"context"
	"errors"
	"fmt"
)

var ErrInvalidFormat = errors.New("requested format not available")

type FormatSelector struct {
	maxSize int64
}

func NewFormatSelector(maxSize int64) *FormatSelector {
	return &FormatSelector{maxSize: maxSize}
}

type ChooseFormatResult struct {
	FormatID string
	Height   int
	Filesize int64
}

func (s *FormatSelector) ChooseFormat(
	ctx context.Context,
	fetcher MetadataFetcher,
	url string,
) (*ChooseFormatResult, error) {

	meta, err := fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video metadata: %w", err)
	}

	res := s.Select(meta)
	if res == nil {
		return nil, ErrInvalidFormat
	}

	return res, nil
}

func (s *FormatSelector) Select(meta *VideoMeta) *ChooseFormatResult {

	var best ChooseFormatResult
	var found bool

	for _, f := range meta.Formats {

		if f.Ext != "mp4" ||
			f.Vcodec == "none" ||
			f.Acodec == "none" ||
			f.Height == 0 {
			continue
		}

		size := estimateSize(&f, meta.Duration)

		if size == 0 || size > s.maxSize {
			continue
		}

		if !found || f.Height > best.Height {
			best = ChooseFormatResult{
				FormatID: f.FormatID,
				Height:   f.Height,
				Filesize: size,
			}
			found = true
		}
	}

	if !found {
		return nil
	}

	return &best
}

func estimateSize(f *Format, duration float64) int64 {

	if f.Filesize != 0 {
		return f.Filesize
	}

	if f.FilesizeApprox != 0 {
		return f.FilesizeApprox
	}

	if f.Tbr != 0 {
		return int64(f.Tbr * duration * 125)
	}

	return 0
}
