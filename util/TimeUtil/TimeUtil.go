package TimeUtil

import (
	"time"
)

func CurrentMillisecondTimestamp() uint64 {
	return uint64(time.Now().Unix()) * uint64(1000)
}
