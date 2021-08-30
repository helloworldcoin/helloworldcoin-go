package SystemUtil

import (
	"helloworld-blockchain-go/util/LogUtil"
	"os"
)

func ErrorExit(message string, exception interface{}) {
	LogUtil.Error("system error occurred, and exited, please check the error！"+message, exception)
	os.Exit(1)
}
