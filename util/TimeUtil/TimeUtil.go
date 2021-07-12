package TimeUtil

import (
	"time"
)

func MillisecondTimestamp() uint64 {
	return uint64(time.Now().Unix()) * uint64(1000)
}
