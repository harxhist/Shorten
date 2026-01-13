package helper

import (
	"time"
)

func CalculateLatency(startTime time.Time) time.Duration {
    return time.Since(startTime)
}