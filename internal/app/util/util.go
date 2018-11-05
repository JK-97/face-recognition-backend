package util

import "time"

// NowMilli returns current unix time in milliseconds
func NowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
