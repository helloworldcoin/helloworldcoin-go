package TransactionTool

import (
	"helloworldcoin-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/util/DataStructureUtil"

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

/**
 * 区块新产生的地址是否存在重复
 */
func IsExistDuplicateNewAddress(transaction *model.Transaction) bool {
	var newAddresss []string
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			address := output.Address
			newAddresss = append(newAddresss, address)
		}
	}
	return DataStructureUtil.IsExistDuplicateElement(&newAddresss)
}
func GetTransactionOutputId(transactionOutput *model.TransactionOutput) string {
	return BlockchainDatabaseKeyTool.BuildTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
}

/**
 * 交易中是否存在重复的[未花费交易输出]
 */
func IsExistDuplicateUtxo(transaction *model.Transaction) bool {
	var utxoIds []string
	inputs := transaction.Inputs
	if inputs != nil {
		for _, transactionInput := range inputs {
			unspentTransactionOutput := transactionInput.UnspentTransactionOutput
			utxoId := GetTransactionOutputId(&unspentTransactionOutput)
			utxoIds = append(utxoIds, utxoId)
		}
	}
	return DataStructureUtil.IsExistDuplicateElement(&utxoIds)
}
