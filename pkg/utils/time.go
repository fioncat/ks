package utils

import (
	"fmt"
	"time"
)

func FormatTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	since := time.Since(t)

	years := int(since.Hours() / 24 / 365)
	if years > 0 {
		return fmt.Sprintf("%dy", years)
	}

	days := int(since.Hours() / 24)
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}

	hours := int(since.Hours())
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}

	mins := int(since.Minutes())
	if mins > 0 {
		return fmt.Sprintf("%dm", mins)
	}

	return fmt.Sprintf("%ds", int(since.Seconds()))
}
