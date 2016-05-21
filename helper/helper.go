package helper

import "time"

//Convert time to miliseconds..
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
