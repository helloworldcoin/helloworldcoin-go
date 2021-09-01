package SystemUtil

import (
	"helloworld-blockchain-go/util/LogUtil"
	"os"
	"path/filepath"
	"runtime"
)

func ErrorExit(message string, exception interface{}) {
	LogUtil.Error("system error occurred, and exited, please check the error！"+message, exception)
	os.Exit(1)
}
func SystemRootDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	systemRootDirectory := filepath.Join(filepath.Dir(b), "../..")
	return systemRootDirectory
}
