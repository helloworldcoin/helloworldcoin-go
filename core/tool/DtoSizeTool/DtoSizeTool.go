package DtoSizeTool

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/setting/BlockSetting"
	"helloworld-blockchain-go/setting/ScriptSetting"
	"helloworld-blockchain-go/setting/TransactionSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"strconv"
	"unicode/utf8"
)

//region 校验大小
/**
 * 校验区块大小。用来限制区块的大小。
 * 注意：校验区块的大小，不仅要校验区块的大小
 * ，还要校验区块内部各个属性(时间戳、前哈希、随机数、交易)的大小。
 */
func CheckBlockSize(blockDto *dto.BlockDto) bool {
	//区块的时间戳的长度不需要校验  假设时间戳长度不正确，则在随后的业务逻辑中走不通

	//区块的前哈希的长度不需要校验  假设前哈希长度不正确，则在随后的业务逻辑中走不通

	//校验区块随机数大小
	nonceSize := sizeOfString(blockDto.Nonce)
	if nonceSize != BlockSetting.NONCE_CHARACTER_COUNT {
		LogUtil.Debug("nonce[" + blockDto.Nonce + "]长度非法。")
		return false
	}

	//校验每一笔交易大小
	transactionDtos := blockDto.Transactions
	if transactionDtos != nil {
		for _, transactionDto := range transactionDtos {
			if !CheckTransactionSize(transactionDto) {
				LogUtil.Debug("交易数据异常，交易大小非法。")
				return false
			}
		}
	}

	//校验区块占用的存储空间
	blockSize := CalculateBlockSize(blockDto)
	if blockSize > BlockSetting.BLOCK_MAX_CHARACTER_COUNT {
		LogUtil.Debug("区块数据的大小是[" + StringUtil.ValueOfUint64(blockSize) + "]超过了限制[" + StringUtil.ValueOfUint64(BlockSetting.BLOCK_MAX_CHARACTER_COUNT) + "]。")
		return false
	}
	return true
}

/**
 * 校验交易的大小：用来限制交易的大小。
 * 注意：校验交易的大小，不仅要校验交易的大小
 * ，还要校验交易内部各个属性(交易输入、交易输出)的大小。
 */
func CheckTransactionSize(transactionDto *dto.TransactionDto) bool {
	//校验交易输入
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			//交易的未花费输出大小不需要校验  假设不正确，则在随后的业务逻辑中走不通

			//校验脚本大小
			inputScriptDto := transactionInputDto.InputScript
			//校验输入脚本的大小
			if !checkInputScriptSize(inputScriptDto) {
				return false
			}
		}
	}

	//校验交易输出
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			//交易输出金额大小不需要校验  假设不正确，则在随后的业务逻辑中走不通

			//校验脚本大小
			outputScriptDto := transactionOutputDto.OutputScript
			//校验输出脚本的大小
			if !checkOutputScriptSize(outputScriptDto) {
				return false
			}

		}
	}

	//校验整笔交易大小十分合法
	transactionSize := CalculateTransactionSize(transactionDto)
	if transactionSize > TransactionSetting.TRANSACTION_MAX_CHARACTER_COUNT {
		LogUtil.Debug("交易的大小是[" + StringUtil.ValueOfUint64(transactionSize) + "]，超过了限制值[" + StringUtil.ValueOfUint64(TransactionSetting.TRANSACTION_MAX_CHARACTER_COUNT) + "]。")
		return false
	}
	return true
}

/**
 * 校验输入脚本的大小
 */
func checkInputScriptSize(inputScriptDto *dto.InputScriptDto) bool {
	//校验脚本大小
	if !checkScriptSize(inputScriptDto) {
		return false
	}
	return true
}

/**
 * 校验输出脚本的大小
 */
func checkOutputScriptSize(outputScriptDto *dto.OutputScriptDto) bool {
	//校验脚本大小
	if !checkScriptSize(outputScriptDto) {
		return false
	}
	return true
}

/**
 * 校验脚本的大小
 */
func checkScriptSize(scriptDto *[]string) bool {
	//脚本内的操作码、操作数大小不需要校验，因为操作码、操作数不合规，在脚本结构上就构不成一个合格的脚本。
	if calculateScriptSize(scriptDto) > ScriptSetting.SCRIPT_MAX_CHARACTER_COUNT {
		LogUtil.Debug("交易校验失败：交易输出脚本大小超出限制。")
		return false
	}
	return true
}

//endregion

//region 计算大小
func CalculateBlockSize(blockDto *dto.BlockDto) uint64 {
	size := uint64(0)
	timestamp := blockDto.Timestamp
	size += sizeOfUint64(timestamp)

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
	size += sizeOfUint64(value)
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
	size += sizeOfUint64(transactionOutputIndex)
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
	return uint64(utf8.RuneCountInString(value))
}

func sizeOfUint64(number uint64) uint64 {
	return uint64(utf8.RuneCountInString(strconv.FormatUint(number, 10)))
}

//endregion
