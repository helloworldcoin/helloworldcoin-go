package OperateSystemUtil

/*
 @author king 409060350@qq.com
*/

import (
	"runtime"
)

func IsWindowsOperateSystem() bool {
	return "windows" == runtime.GOOS
}

func IsMacOperateSystem() bool {
	return "darwin" == runtime.GOOS
}
