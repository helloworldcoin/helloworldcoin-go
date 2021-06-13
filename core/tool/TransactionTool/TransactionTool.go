package TransactionTool

import (
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"

	"helloworldcoin-go/core/Model"
)

func GetTransactionOutputCount(transaction *Model.Transaction) uint64 {
	outputs := transaction.Outputs
	if outputs == nil {
		return uint64(0)
	}
	return uint64(len(outputs))
}
func CalculateTransactionHash(transaction Model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(&transaction)
	return TransactionDtoTool.CalculateTransactionHash(&transactionDto)
}
