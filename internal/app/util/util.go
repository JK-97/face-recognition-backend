package util

import "time"

// NowMilli returns current unix time in milliseconds
func NowMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// MilliToSeconds translate milliseconds to seconds
func MilliToSeconds(milliseconds int64) int64 {
    return milliseconds / 1000
}
