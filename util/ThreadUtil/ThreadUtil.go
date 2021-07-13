package ThreadUtil

import (
	"time"
)

func MillisecondSleep(millisecond uint64) {
	time.Sleep(time.Duration(millisecond * 1000000))
}
