package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/Model/Script"
	"helloworldcoin-go/core/Model/Script/BooleanCodeEnum"
	"helloworldcoin-go/core/Model/Script/OperationCodeEnum"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/util/StringUtil"
)

type VirtualMachine struct {
}

func (this *VirtualMachine) ExecuteScript(transactionEnvironment *Model.Transaction, script Script.Script) *Script.ScriptExecuteResult {

	var stack Script.ScriptExecuteResult
	for i, _ := range script {
		operationCode := script[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.IsEquals(OperationCodeEnum.OP_DUP.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("指令运行异常")
			}
			stack.Push(*stack.Peek())
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_HASH160.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("指令运行异常")
			}
			publicKey := stack.Pop()
			publicKeyHash := AccountUtil.PublicKeyHashFromPublicKey(*publicKey)
			stack.Push(publicKeyHash)
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_EQUALVERIFY.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("指令运行异常")
			}
			if !StringUtil.IsEquals(*stack.Pop(), *stack.Pop()) {
				panic("脚本执行失败")
			}
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_CHECKSIG.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("指令运行异常")
			}
			publicKey := stack.Pop()
			signature := stack.Pop()
			bytesSignature := ByteUtil.HexStringToBytes(*signature)
			verifySignatureSuccess := TransactionTool.VerifySignature(transactionEnvironment, *publicKey, bytesSignature)
			if !verifySignatureSuccess {
				panic("脚本执行失败")
			}
			stack.Push(ByteUtil.BytesToHexString(BooleanCodeEnum.TRUE.Code))
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_PUSHDATA.Code, bytesOperationCode) {
			if len(script) < i+2 {
				panic("指令运行异常")
			}
			i++
			stack.Push(script[i])
		} else {
			panic("不能识别的操作码")
		}
	}
	return &stack
}
