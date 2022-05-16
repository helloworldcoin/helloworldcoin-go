package SizeTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/tool/DtoSizeTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
)

//region check Size
/**
 * check Block Size: used to limit the size of the block.
 */
func CheckBlockSize(block *model.Block) bool {
	return DtoSizeTool.CheckBlockSize(Model2DtoTool.Block2BlockDto(block))
}

/**
 * Check transaction size: used to limit the size of the transaction.
 */
func CheckTransactionSize(transaction *model.Transaction) bool {
	return DtoSizeTool.CheckTransactionSize(Model2DtoTool.Transaction2TransactionDto(transaction))
}

//endregion

//region calculate Size
func CalculateBlockSize(block *model.Block) uint64 {
	return DtoSizeTool.CalculateBlockSize(Model2DtoTool.Block2BlockDto(block))
}
func CalculateTransactionSize(transaction *model.Transaction) uint64 {
	return DtoSizeTool.CalculateTransactionSize(Model2DtoTool.Transaction2TransactionDto(transaction))
}

//endregion
