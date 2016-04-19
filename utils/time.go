package utils

import "time"

// GetNowTruncateSecond ... returns the truncated result of rounding time Now down to a multiple of second
func GetNowTruncateSecond() time.Time {
	return time.Now().Truncate(time.Second)
}
