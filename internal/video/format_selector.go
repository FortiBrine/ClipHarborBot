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
) (*Format, error) {
	meta, err := fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetching metadata: %w", err)
	}

	format := s.selectBest(meta)
	if format == nil {
		return nil, ErrInvalidFormat
	}

	return format, nil
}

func (s *FormatSelector) selectBest(meta *Meta) (format *Format) {
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

		if f.Height > format.Height {
			format = f
		}
	}
	return
}
