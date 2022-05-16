package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/script/BooleanCode"
	"helloworldcoin-go/core/model/script/OperationCode"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/StringStack"
	"helloworldcoin-go/util/StringUtil"
)

type VirtualMachine struct {
}

func NewVirtualMachine() *VirtualMachine {
	var virtualMachine VirtualMachine
	return &virtualMachine
}

func (this *VirtualMachine) Execute(transactionEnvironment *model.Transaction, script *model.Script) *model.Result {

	stack := StringStack.NewStringStack()
	for i := 0; i < len(*script); i++ {
		operationCode := (*script)[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCode.OP_DUP.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("Virtual Machine Execute Error.")
			}
			stack.Push(*stack.Peek())
		} else if ByteUtil.Equals(OperationCode.OP_HASH160.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("Virtual Machine Execute Error.")
			}
			publicKey := stack.Pop()
			publicKeyHash := AccountUtil.PublicKeyHashFromPublicKey(*publicKey)
			stack.Push(publicKeyHash)
		} else if ByteUtil.Equals(OperationCode.OP_EQUALVERIFY.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("Virtual Machine Execute Error.")
			}
			if !StringUtil.Equals(*stack.Pop(), *stack.Pop()) {
				panic("Virtual Machine Execute Error.")
			}
		} else if ByteUtil.Equals(OperationCode.OP_CHECKSIG.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("Virtual Machine Execute Error.")
			}
			publicKey := stack.Pop()
			signature := stack.Pop()
			bytesSignature := ByteUtil.HexStringToBytes(*signature)
			verifySignatureSuccess := TransactionTool.VerifySignature(transactionEnvironment, *publicKey, bytesSignature)
			if !verifySignatureSuccess {
				panic("Virtual Machine Execute Error.")
			}
			stack.Push(ByteUtil.BytesToHexString(BooleanCode.TRUE.Code))
		} else if ByteUtil.Equals(OperationCode.OP_PUSHDATA.Code, bytesOperationCode) {
			if len(*script) < i+2 {
				panic("Virtual Machine Execute Error.")
			}
			i++
			stack.Push((*script)[i])
		} else {
			panic("Virtual Machine Execute Error.")
		}
	}
	return stack
}
