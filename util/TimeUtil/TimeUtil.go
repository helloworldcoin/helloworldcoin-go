package TimeUtil

import (
	"time"
)

func CurrentMillisecondTimestamp() uint64 {
	return uint64(time.Now().Second()) * uint64(1000)
}
