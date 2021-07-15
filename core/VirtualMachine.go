package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/Model/script"
)

type VirtualMachine struct {
}

func (this *VirtualMachine) ExecuteScript(transactionEnvironment *Model.Transaction, script1 *[]string) *script.ScriptExecuteResult {

	var ret11 script.ScriptExecuteResult
	return &ret11
}
