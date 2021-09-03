package model

import (
	"helloworld-blockchain-go/util/StringStack"
)

type Script = []string
type InputScript = Script
type OutputScript = Script
type ScriptExecuteResult = StringStack.StringStack
