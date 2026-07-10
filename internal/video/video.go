package video

type Meta struct {
	Duration float64   `json:"duration"`
	Formats  []*Format `json:"formats"`
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

func (f *Format) estimateSize(duration float64) int64 {
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
