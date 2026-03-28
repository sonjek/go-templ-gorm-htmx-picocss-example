package datetime

import (
	"time"

	"github.com/dustin/go-humanize"
)

func FormatToAgo(datetime time.Time) string {
	return humanize.Time(datetime)
}

func FormatToDateTime(datetime time.Time) string {
	return datetime.Format("2006-01-02 15:04")
}
