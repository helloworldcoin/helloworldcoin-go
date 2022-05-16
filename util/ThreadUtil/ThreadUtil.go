package ThreadUtil

/*
 @author x.king xdotking@gmail.com
*/

import (
	"time"
)

func MillisecondSleep(millisecond uint64) {
	time.Sleep(time.Duration(millisecond * 1000000))
}
