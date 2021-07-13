package OperateSystemUtil

import (
	"runtime"
)

func IsWindowsOperateSystem() bool {
	return "windows" == runtime.GOOS
}

func IsMacOperateSystem() bool {
	return "darwin" == runtime.GOOS
}
