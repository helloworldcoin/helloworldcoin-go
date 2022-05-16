package TransactionTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/SizeTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/StringsUtil"
)

/**
 * Get Total Input Value Of Transaction
 */
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

/**
 * Get Total Output Value Of Transaction
 */
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
 * Get Total Fees Of Transaction
 */
func GetTransactionFee(transaction *model.Transaction) uint64 {
	if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		transactionFee := GetInputValue(transaction) - GetOutputValue(transaction)
		return transactionFee
	} else if transaction.TransactionType == TransactionType.COINBASE_TRANSACTION {
		return 0
	} else {
		panic(nil)
	}
}

/**
 * Get Fee Rate Of Transaction
 */
func GetTransactionFeeRate(transaction *model.Transaction) uint64 {
	if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		return GetTransactionFee(transaction) / SizeTool.CalculateTransactionSize(transaction)
	} else if transaction.TransactionType == TransactionType.COINBASE_TRANSACTION {
		return 0
	} else {
		panic("")
	}
}

func getSignatureHashAllRawMaterial(transaction *model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.GetSignatureHashAllRawMaterial(transactionDto)
}

func signature(privateKey string, transaction *model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.Signature(privateKey, transactionDto)
}

func VerifySignature(transaction *model.Transaction, publicKey string, bytesSignature []byte) bool {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.VerifySignature(transactionDto, publicKey, bytesSignature)
}

func CalculateTransactionHash(transaction *model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.CalculateTransactionHash(transactionDto)
}

func GetTransactionInputCount(transaction *model.Transaction) uint64 {
	inputs := transaction.Inputs
	if inputs == nil {
		return uint64(0)
	}
	return uint64(len(inputs))
}
func GetTransactionOutputCount(transaction *model.Transaction) uint64 {
	outputs := transaction.Outputs
	if outputs == nil {
		return uint64(0)
	}
	return uint64(len(outputs))
}

func GetTransactionOutputId(transactionOutput *model.TransactionOutput) string {
	return BlockchainDatabaseKeyTool.BuildTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
}

/**
 * Check Transaction Value
 */
func CheckTransactionValue(transaction *model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil {
		//Check Transaction Input Value
		for _, input := range inputs {
			if !CheckValue(input.UnspentTransactionOutput.Value) {
				LogUtil.Debug("Transaction value is illegal.")
				return false
			}
		}
	}

	outputs := transaction.Outputs
	if outputs != nil {
		//Check Transaction Output Value
		for _, output := range outputs {
			if !CheckValue(output.Value) {
				LogUtil.Debug("Transaction value is illegal.")
				return false
			}
		}
	}

	//further check by transaction type
	if transaction.TransactionType == TransactionType.COINBASE_TRANSACTION {
		//There is no need to check, skip.
	} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		//The transaction input value must be greater than or equal to the transaction output value
		inputsValue := GetInputValue(transaction)
		outputsValue := GetOutputValue(transaction)
		if inputsValue < outputsValue {
			LogUtil.Debug("Transaction value is illegal.")
			return false
		}
		return true
	} else {
		panic(nil)
	}
	return true
}

/**
 * Check whether the transaction value is legal: this is used to limit the maximum value, minimum value, decimal places, etc. of the transaction value
 */
func CheckValue(transactionAmount uint64) bool {
	//The transaction value cannot be less than or equal to 0
	if transactionAmount <= 0 {
		return false
	}
	//The maximum value is 2^64
	//The reserved decimal place is 0
	return true
}

/**
 * Check if the address in the transaction is a P2PKH address
 */
func CheckPayToPublicKeyHashAddress(transaction *model.Transaction) bool {
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			if !AccountUtil.IsPayToPublicKeyHashAddress(output.Address) {
				LogUtil.Debug("Transaction address is illegal.")
				return false
			}
		}
	}
	return true
}

/**
 * Check if the script in the transaction is a P2PKH script
 */
func CheckPayToPublicKeyHashScript(transaction *model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil {
		for _, input := range inputs {
			if !ScriptTool.IsPayToPublicKeyHashInputScript(input.InputScript) {
				return false
			}
		}
	}
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			if !ScriptTool.IsPayToPublicKeyHashOutputScript(output.OutputScript) {
				return false
			}
		}
	}
	return true
}

/**
 * Is there a duplicate [unspent transaction output] in the transaction
 */
func IsExistDuplicateUtxo(transaction *model.Transaction) bool {
	var utxoIds []string
	inputs := transaction.Inputs
	if inputs != nil {
		for _, transactionInput := range inputs {
			unspentTransactionOutput := transactionInput.UnspentTransactionOutput
			utxoId := GetTransactionOutputId(unspentTransactionOutput)
			utxoIds = append(utxoIds, utxoId)
		}
	}
	return StringsUtil.HasDuplicateElement(&utxoIds)
}

/**
 * Whether the newly generated address of the block is duplicated
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
	return StringsUtil.HasDuplicateElement(&newAddresss)
}
