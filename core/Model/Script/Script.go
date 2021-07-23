package Script

import (
	"helloworldcoin-go/util/StringStack"
)

type VirtualMachine struct {
}
type Script = []string
type InputScript = Script
type OutputScript = Script
type ScriptExecuteResult = StringStack.StringStack
