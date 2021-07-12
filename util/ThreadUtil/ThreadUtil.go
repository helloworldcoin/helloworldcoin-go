package ThreadUtil

import (
	"time"
)

func MillisecondSleep(millis uint64) {
	time.Sleep(time.Duration(millis * 1000000))
}
