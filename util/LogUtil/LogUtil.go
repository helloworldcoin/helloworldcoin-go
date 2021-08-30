package LogUtil

import (
	"fmt"
	"helloworld-blockchain-go/util/JsonUtil"
)

func Debug(message string) {
	fmt.Println(message)
}
func Error(message string, exception interface{}) {
	fmt.Println(message + JsonUtil.ToString(exception))
}
