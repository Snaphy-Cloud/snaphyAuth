package helper

import (
	"time"
	"github.com/satori/go.uuid"
)

//Convert time to miliseconds..
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}


func CreateUUID() (Id string) {
	return uuid.NewV4().String() //unique number provider
}
