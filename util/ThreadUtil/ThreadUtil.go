package ThreadUtil

/*
 @author king 409060350@qq.com
*/

import (
	"time"
)

func MillisecondSleep(millisecond uint64) {
	time.Sleep(time.Duration(millisecond * 1000000))
}
