package DtoSizeTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/setting/BlockSetting"
	"helloworldcoin-go/setting/ScriptSetting"
	"helloworldcoin-go/setting/TransactionSetting"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/StringUtil"
)

//region check Size
/**
 * check Block Size: used to limit the size of the block.
 */
func CheckBlockSize(blockDto *dto.BlockDto) bool {
	//The length of the timestamp of the block does not need to be verified. If the timestamp length is incorrect, it will not work in the business logic.

	//The length of the previous hash of the block does not need to be verified. If the previous hash length is incorrect, it will not work in the business logic.

	//Check block nonce size
	nonceSize := sizeOfString(blockDto.Nonce)
	if nonceSize != BlockSetting.NONCE_CHARACTER_COUNT {
		LogUtil.Debug("Illegal nonce length.")
		return false
	}

	//Check the size of each transaction
	transactionDtos := blockDto.Transactions
	if transactionDtos != nil {
		for _, transactionDto := range transactionDtos {
			if !CheckTransactionSize(transactionDto) {
				LogUtil.Debug("Illegal transaction size.")
				return false
			}
		}
	}

	//Check Block size
	blockSize := CalculateBlockSize(blockDto)
	if blockSize > BlockSetting.BLOCK_MAX_CHARACTER_COUNT {
		LogUtil.Debug("Block size exceeds limit.")
		return false
	}
	return true
}

/**
 * Check transaction size: used to limit the size of the transaction.
 */
func CheckTransactionSize(transactionDto *dto.TransactionDto) bool {
	//Check transaction input size
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			//The unspent output size of the transaction does not need to be verified. If the assumption is incorrect, it will not work in the business logic.

			//Check script size
			inputScriptDto := transactionInputDto.InputScript
			//Check the size of the input script
			if !checkInputScriptSize(inputScriptDto) {
				return false
			}
		}
	}

	//Check transaction output size
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			//The size of the transaction output amount does not need to be verified. If the assumption is incorrect, it will not work in the business logic.

			//Check script size
			outputScriptDto := transactionOutputDto.OutputScript
			//Check the size of the output script
			if !checkOutputScriptSize(outputScriptDto) {
				return false
			}

		}
	}

	//Check transaction size
	transactionSize := CalculateTransactionSize(transactionDto)
	if transactionSize > TransactionSetting.TRANSACTION_MAX_CHARACTER_COUNT {
		LogUtil.Debug("Transaction size exceeds limit.")
		return false
	}
	return true
}

/**
 * check Input Script Size
 */
func checkInputScriptSize(inputScriptDto *dto.InputScriptDto) bool {
	if !checkScriptSize(inputScriptDto) {
		return false
	}
	return true
}

/**
 * check Output Script Size
 */
func checkOutputScriptSize(outputScriptDto *dto.OutputScriptDto) bool {
	if !checkScriptSize(outputScriptDto) {
		return false
	}
	return true
}

/**
 * check Script Size
 */
func checkScriptSize(scriptDto *[]string) bool {
	//There is no need to check the size of opcodes and operands in the script.
	//Illegal opcodes, illegal operands cannot constitute a legal script.
	//Illegal script will not work in the business logic.
	if calculateScriptSize(scriptDto) > ScriptSetting.SCRIPT_MAX_CHARACTER_COUNT {
		LogUtil.Debug("Script size exceeds limit.")
		return false
	}
	return true
}

//endregion

//region calculate Size
func CalculateBlockSize(blockDto *dto.BlockDto) uint64 {
	size := uint64(0)
	timestamp := blockDto.Timestamp
	size += sizeOfNumber(timestamp)

	previousBlockHash := blockDto.PreviousHash
	size += sizeOfString(previousBlockHash)

	nonce := blockDto.Nonce
	size += sizeOfString(nonce)
	transactionDtos := blockDto.Transactions
	for _, transactionDto := range transactionDtos {
		size += CalculateTransactionSize(transactionDto)
	}
	return size
}
func CalculateTransactionSize(transactionDto *dto.TransactionDto) uint64 {
	size := uint64(0)
	transactionInputDtos := transactionDto.Inputs
	size += calculateTransactionInputsSize(transactionInputDtos)
	transactionOutputDtos := transactionDto.Outputs
	size += calculateTransactionOutputsSize(transactionOutputDtos)
	return size
}
func calculateTransactionOutputsSize(transactionOutputDtos []*dto.TransactionOutputDto) uint64 {
	size := uint64(0)
	if transactionOutputDtos == nil || len(transactionOutputDtos) == 0 {
		return size
	}
	for _, transactionOutputDto := range transactionOutputDtos {
		size += calculateTransactionOutputSize(transactionOutputDto)
	}
	return size
}
func calculateTransactionOutputSize(transactionOutputDto *dto.TransactionOutputDto) uint64 {
	size := uint64(0)
	outputScriptDto := transactionOutputDto.OutputScript
	size += calculateScriptSize(outputScriptDto)
	value := transactionOutputDto.Value
	size += sizeOfNumber(value)
	return size
}
func calculateTransactionInputsSize(inputs []*dto.TransactionInputDto) uint64 {
	size := uint64(0)
	if inputs == nil || len(inputs) == 0 {
		return size
	}
	for _, transactionInputDto := range inputs {
		size += calculateTransactionInputSize(transactionInputDto)
	}
	return size
}
func calculateTransactionInputSize(input *dto.TransactionInputDto) uint64 {
	size := uint64(0)
	transactionHash := input.TransactionHash
	size += sizeOfString(transactionHash)
	transactionOutputIndex := input.TransactionOutputIndex
	size += sizeOfNumber(transactionOutputIndex)
	inputScriptDto := input.InputScript
	size += calculateScriptSize(inputScriptDto)
	return size
}
func calculateScriptSize(script *[]string) uint64 {
	size := uint64(0)
	if script == nil || len(*script) == 0 {
		return size
	}
	for _, scriptCode := range *script {
		size += sizeOfString(scriptCode)
	}
	return size
}
func sizeOfString(value string) uint64 {
	return StringUtil.Length(value)
}
func sizeOfNumber(number uint64) uint64 {
	return StringUtil.Length(StringUtil.ValueOfUint64(number))
}

//endregion
