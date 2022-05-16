package BlockTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockDtoTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/StringUtil"
	"helloworldcoin-go/util/StringsUtil"
	"helloworldcoin-go/util/TimeUtil"

	"helloworldcoin-go/core/model"
)

func CalculateBlockHash(block *model.Block) string {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return BlockDtoTool.CalculateBlockHash(blockDto)
}

func CalculateBlockMerkleTreeRoot(block *model.Block) string {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return BlockDtoTool.CalculateBlockMerkleTreeRoot(blockDto)
}

func GetTransactionCount(block *model.Block) uint64 {
	transactions := block.Transactions
	return uint64(len(transactions))
}

func GetTransactionOutputCount(block *model.Block) uint64 {
	transactionOutputCount := uint64(0)
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionOutputCount = transactionOutputCount + TransactionTool.GetTransactionOutputCount(transaction)
		}
	}
	return transactionOutputCount
}

func GetBlockFee(block *model.Block) uint64 {
	blockFee := uint64(0)
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == TransactionType.COINBASE_TRANSACTION {
				continue
			} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
				fee := TransactionTool.GetTransactionFee(transaction)
				blockFee += fee
			} else {
				panic(nil)
			}
		}
	}
	return blockFee
}

func GetWritedIncentiveValue(block *model.Block) uint64 {
	return block.Transactions[0].Outputs[0].Value
}

func GetNextBlockHeight(currentBlock *model.Block) uint64 {
	var nextBlockHeight uint64
	if currentBlock == nil {
		nextBlockHeight = GenesisBlockSetting.HEIGHT + uint64(1)
	} else {
		nextBlockHeight = currentBlock.Height + uint64(1)
	}
	return nextBlockHeight
}

func CheckBlockHeight(previousBlock *model.Block, currentBlock *model.Block) bool {
	if previousBlock == nil {
		return (GenesisBlockSetting.HEIGHT + 1) == currentBlock.Height
	} else {
		return (previousBlock.Height + 1) == currentBlock.Height
	}
}

func CheckPreviousBlockHash(previousBlock *model.Block, currentBlock *model.Block) bool {
	if previousBlock == nil {
		return StringUtil.Equals(GenesisBlockSetting.HASH, currentBlock.PreviousHash)
	} else {
		return StringUtil.Equals(previousBlock.Hash, currentBlock.PreviousHash)
	}
}

func CheckBlockTimestamp(previousBlock *model.Block, currentBlock *model.Block) bool {
	if currentBlock.Timestamp > TimeUtil.MillisecondTimestamp() {
		return false
	}
	if previousBlock == nil {
		return true
	} else {
		return currentBlock.Timestamp > previousBlock.Timestamp
	}
}

func IsExistDuplicateNewHash(block *model.Block) bool {
	var newHashs []string
	blockHash := block.Hash
	newHashs = append(newHashs, blockHash)
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHash := transaction.TransactionHash
			newHashs = append(newHashs, transactionHash)
		}
	}
	return StringsUtil.HasDuplicateElement(&newHashs)
}

func IsExistDuplicateNewAddress(block *model.Block) bool {
	var newAddresss []string
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					address := output.Address
					newAddresss = append(newAddresss, address)
				}
			}
		}
	}
	return StringsUtil.HasDuplicateElement(&newAddresss)
}

func IsExistDuplicateUtxo(block *model.Block) bool {
	var utxoIds []string
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					utxoId := TransactionTool.GetTransactionOutputId(unspentTransactionOutput)
					utxoIds = append(utxoIds, utxoId)
				}
			}
		}
	}
	return StringsUtil.HasDuplicateElement(&utxoIds)
}

func IsBlockEquals(block1 *model.Block, block2 *model.Block) bool {
	return StringUtil.Equals(block1.Hash, block2.Hash)
}

func FormatDifficulty(difficulty string) string {
	return StringUtil.PrefixPadding(difficulty, 64, "0")
}
