package model

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/util/StringStack"
)

type Script = []string
type InputScript = Script
type OutputScript = Script
type ScriptExecuteResult = StringStack.StringStack
