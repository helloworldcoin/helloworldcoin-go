package BlockTool

import (
	"helloworld-blockchain-go/core/Model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockDtoTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/DataStructureUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/TimeUtil"

	"helloworld-blockchain-go/core/Model"
)

func CalculateBlockHash(block *Model.Block) string {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return BlockDtoTool.CalculateBlockHash(blockDto)
}

func CalculateBlockMerkleTreeRoot(block *Model.Block) string {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return BlockDtoTool.CalculateBlockMerkleTreeRoot(blockDto)
}

func GetTransactionCount(block *Model.Block) uint64 {
	transactions := block.Transactions
	return uint64(len(transactions))
}
func GetTransactionOutputCount(block *Model.Block) uint64 {
	transactionOutputCount := uint64(0)
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionOutputCount = transactionOutputCount + TransactionTool.GetTransactionOutputCount(transaction)
		}
	}
	return transactionOutputCount
}
func GetBlockFee(block *Model.Block) uint64 {
	blockFee := uint64(0)
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			if transaction.TransactionType == TransactionType.GENESIS_TRANSACTION {
				continue
			} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
				fee := TransactionTool.GetTransactionFee(transaction)
				blockFee += fee
			} else {
			}
		}
	}
	return blockFee
}
func GetWritedIncentiveValue(block *Model.Block) uint64 {
	return block.Transactions[0].Outputs[0].Value
}
func GetNextBlockHeight(currentBlock *Model.Block) uint64 {
	var nextBlockHeight uint64
	if currentBlock == nil {
		nextBlockHeight = GenesisBlockSetting.HEIGHT + uint64(1)
	} else {
		nextBlockHeight = currentBlock.Height + uint64(1)
	}
	return nextBlockHeight
}
func CheckBlockHeight(previousBlock *Model.Block, currentBlock *Model.Block) bool {
	if previousBlock == nil {
		return (GenesisBlockSetting.HEIGHT + 1) == currentBlock.Height
	} else {
		return (previousBlock.Height + 1) == currentBlock.Height
	}
}
func CheckPreviousBlockHash(previousBlock *Model.Block, currentBlock *Model.Block) bool {
	if previousBlock == nil {
		return StringUtil.IsEquals(GenesisBlockSetting.HASH, currentBlock.PreviousHash)
	} else {
		return StringUtil.IsEquals(previousBlock.Hash, currentBlock.PreviousHash)
	}
}
func CheckBlockTimestamp(previousBlock *Model.Block, currentBlock *Model.Block) bool {
	if currentBlock.Timestamp > TimeUtil.MillisecondTimestamp() {
		return false
	}
	if previousBlock == nil {
		return true
	} else {
		return currentBlock.Timestamp > previousBlock.Timestamp
	}
}

/**
 * 区块新产生的哈希是否存在重复
 */
func IsExistDuplicateNewHash(block *Model.Block) bool {
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
	return DataStructureUtil.IsExistDuplicateElement(&newHashs)
}

/**
 * 区块新产生的地址是否存在重复
 */
func IsExistDuplicateNewAddress(block *Model.Block) bool {
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
	return DataStructureUtil.IsExistDuplicateElement(&newAddresss)
}

/**
 * 区块中是否存在重复的[未花费交易输出]
 */
func IsExistDuplicateUtxo(block *Model.Block) bool {
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
	return DataStructureUtil.IsExistDuplicateElement(&utxoIds)
}

/**
 * 简单的校验两个区块是否相等
 * 注意：这里没有严格校验,例如没有校验区块中的交易是否完全一样
 * ，所以即使这里认为两个区块相等，实际上这两个区块还是有可能不相等的。
 */
func IsBlockEquals(block1 *Model.Block, block2 *Model.Block) bool {
	//如果任一区块为为空，则认为两个区块不相等
	if block1 == nil || block2 == nil {
		return false
	}
	blockDto1 := Model2DtoTool.Block2BlockDto(block1)
	blockDto2 := Model2DtoTool.Block2BlockDto(block2)
	return BlockDtoTool.IsBlockEquals(blockDto1, blockDto2)
}
