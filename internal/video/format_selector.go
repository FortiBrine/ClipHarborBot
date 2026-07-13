package video

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
	return new(FormatSelector{
		maxSize: maxSize,
	})
}

func (s *FormatSelector) FetchBest(
	ctx context.Context,
	fetcher MetadataFetcher,
	url string,
) (*Format, int64, error) {
	meta, err := fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, 0, fmt.Errorf("fetching metadata: %w", err)
	}

	format, size := s.selectBest(meta)
	if format == nil {
		return nil, 0, ErrInvalidFormat
	}

	return format, size, nil
}

func (s *FormatSelector) selectBest(meta *Meta) (best *Format, bestSize int64) {
	for _, f := range meta.Formats {
		if f.Ext != "mp4" ||
			f.Vcodec == "none" ||
			f.Acodec == "none" ||
			f.Height == 0 {
			continue
		}

		size := f.estimateSize(meta.Duration)

		if size == 0 || size > s.maxSize {
			continue
		}

		if best == nil || f.Height > best.Height {
			best = f
			bestSize = size
		}
	}
	return
}
