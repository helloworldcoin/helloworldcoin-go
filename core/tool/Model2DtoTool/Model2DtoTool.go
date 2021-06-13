package Model2DtoTool

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/dto"
)

func Block2BlockDto(block *Model.Block) *dto.BlockDto {
	var transactionDtoList []dto.TransactionDto
	transactionList := block.Transactions
	if transactionList != nil {
		for _, transaction := range transactionList {
			transactionDto := Transaction2TransactionDto(&transaction)
			transactionDtoList = append(transactionDtoList, transactionDto)
		}
	}

	var blockDto dto.BlockDto
	blockDto.Timestamp = block.Timestamp
	blockDto.PreviousHash = block.PreviousBlockHash
	blockDto.Transactions = transactionDtoList
	blockDto.Nonce = block.Nonce
	return &blockDto
}

func Transaction2TransactionDto(transaction *Model.Transaction) dto.TransactionDto {
	var inputs []dto.TransactionInputDto
	transactionInputList := transaction.Inputs
	if transactionInputList != nil {
		for _, transactionInput := range transactionInputList {
			var transactionInputDto dto.TransactionInputDto
			transactionInputDto.TransactionHash = transactionInput.UnspentTransactionOutput.TransactionHash
			transactionInputDto.TransactionOutputIndex = transactionInput.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputDto.InputScript = transactionInput.InputScript
			inputs = append(inputs, transactionInputDto)
		}
	}

	var outputs []dto.TransactionOutputDto
	transactionOutputList := transaction.Outputs
	if transactionOutputList != nil {
		for _, transactionOutput := range transactionOutputList {
			transactionOutputDto := TransactionOutput2TransactionOutputDto(&transactionOutput)
			outputs = append(outputs, transactionOutputDto)
		}
	}

	var transactionDto dto.TransactionDto
	transactionDto.Inputs = inputs
	transactionDto.Outputs = outputs
	return transactionDto
}
func TransactionOutput2TransactionOutputDto(transactionOutput *Model.TransactionOutput) dto.TransactionOutputDto {
	var transactionOutputDto dto.TransactionOutputDto
	transactionOutputDto.Value = transactionOutput.Value
	transactionOutputDto.OutputScript = transactionOutput.OutputScript
	return transactionOutputDto
}
