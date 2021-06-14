package TransactionTool

import (
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"

	"helloworldcoin-go/core/model"
)

func GetTransactionOutputCount(transaction *model.Transaction) uint64 {
	outputs := transaction.Outputs
	if outputs == nil {
		return uint64(0)
	}
	return uint64(len(outputs))
}
func CalculateTransactionHash(transaction model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(&transaction)
	return TransactionDtoTool.CalculateTransactionHash(&transactionDto)
}
func GetTransactionFee(transaction *model.Transaction) uint64 {
	transactionFee := GetInputValue(transaction) - GetOutputValue(transaction)
	return transactionFee
}
func GetInputValue(transaction *model.Transaction) uint64 {
	inputs := transaction.Inputs
	total := uint64(0)
	if inputs != nil {
		for _, input := range inputs {
			total += input.UnspentTransactionOutput.Value
		}
	}
	return total
}
func GetOutputValue(transaction *model.Transaction) uint64 {
	outputs := transaction.Outputs
	total := uint64(0)
	if outputs != nil {
		for _, output := range outputs {
			total += output.Value
		}
	}
	return total
}
