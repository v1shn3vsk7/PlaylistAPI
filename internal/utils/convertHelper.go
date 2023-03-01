package utils

import (
	"strings"
	"time"
)
func ParseDuration(durationStr string) (time.Duration, error) {
	durationStr = strings.Replace(durationStr, ":", "h", 1)
	durationStr = strings.Replace(durationStr, ":", "m", 1)
	durationStr += "s"

	return time.ParseDuration(durationStr)
}
