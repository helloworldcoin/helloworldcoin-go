package TransactionTool

import (
	"helloworld-blockchain-go/core/Model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/setting/TransactionSettingTool"
	"helloworld-blockchain-go/util/DataStructureUtil"
	"helloworld-blockchain-go/util/LogUtil"

	"helloworld-blockchain-go/core/Model"
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
func GetTransactionFee(transaction *Model.Transaction) uint64 {
	transactionFee := GetInputValue(transaction) - GetOutputValue(transaction)
	return transactionFee
}
func GetInputValue(transaction *Model.Transaction) uint64 {
	inputs := transaction.Inputs
	total := uint64(0)
	if inputs != nil {
		for _, input := range inputs {
			total += input.UnspentTransactionOutput.Value
		}
	}
	return total
}
func GetOutputValue(transaction *Model.Transaction) uint64 {
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
func IsExistDuplicateNewAddress(transaction *Model.Transaction) bool {
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
func GetTransactionOutputId(transactionOutput *Model.TransactionOutput) string {
	return BlockchainDatabaseKeyTool.BuildTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
}

/**
 * 交易中是否存在重复的[未花费交易输出]
 */
func IsExistDuplicateUtxo(transaction *Model.Transaction) bool {
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

/**
 * 交易中的金额是否符合系统的约束
 */
func CheckTransactionValue(transaction *Model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil {
		//校验交易输入的金额
		for _, input := range inputs {
			if !TransactionSettingTool.CheckTransactionValue(input.UnspentTransactionOutput.Value) {
				LogUtil.Debug("交易金额不合法")
				return false
			}
		}
	}

	outputs := transaction.Outputs
	if outputs != nil {
		//校验交易输出的金额
		for _, output := range outputs {
			if !TransactionSettingTool.CheckTransactionValue(output.Value) {
				LogUtil.Debug("交易金额不合法")
				return false
			}
		}
	}

	//根据交易类型，做进一步的校验
	if transaction.TransactionType == TransactionType.GENESIS_TRANSACTION {
		//没有需要校验的，跳过。
	} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		//交易输入必须要大于等于交易输出
		inputsValue := GetInputValue(transaction)
		outputsValue := GetOutputValue(transaction)
		if inputsValue < outputsValue {
			LogUtil.Debug("交易校验失败：交易的输入必须大于等于交易的输出。不合法的交易。")
			return false
		}
		return true
	} else {
		LogUtil.Debug("区块数据异常，不能识别的交易类型。")
		return false
	}
	return true
}

/**
 * 校验交易中的地址是否是P2PKH地址
 */
func CheckPayToPublicKeyHashAddress(transaction *Model.Transaction) bool {
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			if !AccountUtil.IsPayToPublicKeyHashAddress(output.Address) {
				LogUtil.Debug("交易地址不合法")
				return false
			}
		}
	}
	return true
}

/**
 * 校验交易中的脚本是否是P2PKH脚本
 */
func CheckPayToPublicKeyHashScript(transaction *Model.Transaction) bool {
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
 * 获取待签名数据
 */
func SignatureHashAll(transaction *Model.Transaction) string {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.SignatureHashAll(&transactionDto)
}

/**
 * 验证签名
 */
func VerifySignature(transaction *Model.Transaction, publicKey string, bytesSignature []byte) bool {
	transactionDto := Model2DtoTool.Transaction2TransactionDto(transaction)
	return TransactionDtoTool.VerifySignature(&transactionDto, publicKey, bytesSignature)
}
