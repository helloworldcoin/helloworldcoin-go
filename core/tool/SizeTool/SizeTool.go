package SizeTool

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/core/tool/DtoSizeTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
)

//region 校验大小
/**
 * 校验区块大小。用来限制区块的大小。
 * 注意：校验区块的大小，不仅要校验区块的大小
 * ，还要校验区块内部各个属性(时间戳、前哈希、随机数、交易)的大小。
 */
func CheckBlockSize(block *Model.Block) bool {
	return DtoSizeTool.CheckBlockSize(Model2DtoTool.Block2BlockDto(block))
}

/**
 * 校验交易的大小：用来限制交易的大小。
 * 注意：校验交易的大小，不仅要校验交易的大小
 * ，还要校验交易内部各个属性(交易输入、交易输出)的大小。
 */
func CheckTransactionSize(transaction *Model.Transaction) bool {
	return DtoSizeTool.CheckTransactionSize(Model2DtoTool.Transaction2TransactionDto(transaction))
}

//endregion

//region 计算大小
func CalculateBlockSize(block *Model.Block) uint64 {
	return DtoSizeTool.CalculateBlockSize(Model2DtoTool.Block2BlockDto(block))
}
func calculateTransactionSize(transaction *Model.Transaction) uint64 {
	return DtoSizeTool.CalculateTransactionSize(Model2DtoTool.Transaction2TransactionDto(transaction))
}

//endregion
