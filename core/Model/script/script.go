package script

import "helloworldcoin-go/util/Stack"

type VirtualMachine struct {
}
type Script []string
type InputScript Script
type OutputScript Script
type ScriptExecuteResult Stack.Stack
