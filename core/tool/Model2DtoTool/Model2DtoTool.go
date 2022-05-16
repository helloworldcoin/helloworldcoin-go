package Model2DtoTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/netcore-dto/dto"
)

func Block2BlockDto(block *model.Block) *dto.BlockDto {
	var transactionDtos []*dto.TransactionDto
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionDto := Transaction2TransactionDto(transaction)
			transactionDtos = append(transactionDtos, transactionDto)
		}
	}

	var blockDto dto.BlockDto
	blockDto.Timestamp = block.Timestamp
	blockDto.PreviousHash = block.PreviousHash
	blockDto.Transactions = transactionDtos
	blockDto.Nonce = block.Nonce
	return &blockDto
}

func Transaction2TransactionDto(transaction *model.Transaction) *dto.TransactionDto {
	var inputs []*dto.TransactionInputDto
	transactionInputs := transaction.Inputs
	if transactionInputs != nil {
		for _, transactionInput := range transactionInputs {
			transactionInputDto := &dto.TransactionInputDto{}
			transactionInputDto.TransactionHash = transactionInput.UnspentTransactionOutput.TransactionHash
			transactionInputDto.TransactionOutputIndex = transactionInput.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputDto.InputScript = InputScript2InputScriptDto(transactionInput.InputScript)
			inputs = append(inputs, transactionInputDto)
		}
	}

	var outputs []*dto.TransactionOutputDto
	transactionOutputs := transaction.Outputs
	if transactionOutputs != nil {
		for _, transactionOutput := range transactionOutputs {
			transactionOutputDto := TransactionOutput2TransactionOutputDto(transactionOutput)
			outputs = append(outputs, transactionOutputDto)
		}
	}

	transactionDto := &dto.TransactionDto{}
	transactionDto.Inputs = inputs
	transactionDto.Outputs = outputs
	return transactionDto
}
func TransactionOutput2TransactionOutputDto(transactionOutput *model.TransactionOutput) *dto.TransactionOutputDto {
	transactionOutputDto := &dto.TransactionOutputDto{}
	transactionOutputDto.Value = transactionOutput.Value
	transactionOutputDto.OutputScript = OutputScript2OutputScriptDto(transactionOutput.OutputScript)
	return transactionOutputDto
}
func InputScript2InputScriptDto(inputScript *model.InputScript) *dto.InputScriptDto {
	var inputScriptDto dto.InputScriptDto
	inputScriptDto = append(inputScriptDto, *inputScript...)
	return &inputScriptDto
}
func OutputScript2OutputScriptDto(outputScript *model.OutputScript) *dto.OutputScriptDto {
	var outputScriptDto dto.OutputScriptDto
	outputScriptDto = append(outputScriptDto, *outputScript...)
	return &outputScriptDto
}
