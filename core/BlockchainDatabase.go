package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/Model/BlockchainActionEnum"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworldcoin-go/core/tool/EncodeDecodeTool"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
	"sync"
)

const BLOCKCHAIN_DATABASE_NAME = "BlockchainDatabase"

type BlockchainDatabase struct {
	Consensus         Consensus
	Incentive         Incentive
	CoreConfiguration CoreConfiguration
}

func (blockchainDatabase *BlockchainDatabase) AddBlock(block *Model.Block) bool {
	var lock = sync.Mutex{}
	lock.Lock()
	checkBlock := blockchainDatabase.CheckBlock(block)
	if !checkBlock {
		return false
	}
	kvWriteBatch := blockchainDatabase.createBlockWriteBatch(block, BlockchainActionEnum.ADD_BLOCK)
	KvDbUtil.Write(blockchainDatabase.getBlockchainDatabasePath(), kvWriteBatch)
	lock.Unlock()
	return true
}
func (blockchainDatabase *BlockchainDatabase) DeleteTailBlock() {

}
func (blockchainDatabase *BlockchainDatabase) DeleteBlocks(blockHeight uint64) {
}

func (blockchainDatabase *BlockchainDatabase) CheckBlock(block *Model.Block) bool {
	return true
}
func (blockchainDatabase *BlockchainDatabase) CheckTransaction(block *Model.Transaction) bool {
	return true
}

func (blockchainDatabase *BlockchainDatabase) QueryBlockchainHeight() uint64 {
	return 1
}
func (blockchainDatabase *BlockchainDatabase) QueryBlockchainTransactionHeight() uint64 {
	return uint64(1)
}
func (blockchainDatabase *BlockchainDatabase) QueryBlockchainTransactionOutputHeight() uint64 {
	return uint64(1)
}

func (blockchainDatabase *BlockchainDatabase) QueryTailBlock() *Model.Block {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryBlockByBlockHeight(blockHeight int) *Model.Block {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryBlockByBlockHash(blockHash string) *Model.Block {
	return nil
}

func (blockchainDatabase *BlockchainDatabase) QueryTransactionByTransactionHeight(transactionHeight int) *Model.Transaction {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryTransactionByTransactionHash(transactionHash string) *Model.Transaction {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QuerySourceTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.Transaction {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryDestinationTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.Transaction {
	return nil
}

func (blockchainDatabase *BlockchainDatabase) QueryTransactionOutputByTransactionOutputHeight(transactionOutputHeight uint64) *Model.TransactionOutput {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QueryUnspentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QuerySpentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	return nil
}

func (blockchainDatabase *BlockchainDatabase) QueryTransactionOutputByAddress(address string) *Model.TransactionOutput {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QueryUnspentTransactionOutputByAddress(address string) *Model.TransactionOutput {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) QuerySpentTransactionOutputByAddress(address string) *Model.TransactionOutput {
	return nil
}

func (blockchainDatabase *BlockchainDatabase) GetIncentive() *Incentive {
	return nil
}
func (blockchainDatabase *BlockchainDatabase) GetConsensus() *Consensus {
	return nil
}

func (blockchainDatabase *BlockchainDatabase) getBlockchainDatabasePath() string {
	return FileUtil.NewPath(blockchainDatabase.CoreConfiguration.getCorePath(), BLOCKCHAIN_DATABASE_NAME)
}
func (blockchainDatabase *BlockchainDatabase) createBlockWriteBatch(block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) *KvDbUtil.KvWriteBatch {
	blockchainDatabase.fillBlockProperty(block)
	kvWriteBatch := new(KvDbUtil.KvWriteBatch)

	blockchainDatabase.storeHash(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeAddress(kvWriteBatch, block, blockchainActionEnum)

	blockchainDatabase.storeBlockchainHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeBlockchainTransactionHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeBlockchainTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)

	blockchainDatabase.storeBlockHeightToBlock(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeBlockHashToBlockHeight(kvWriteBatch, block, blockchainActionEnum)

	blockchainDatabase.storeTransactionHeightToTransaction(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionHashToTransactionHeight(kvWriteBatch, block, blockchainActionEnum)

	blockchainDatabase.storeTransactionOutputHeightToTransactionOutput(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch, block, blockchainActionEnum)

	blockchainDatabase.storeAddressToTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeAddressToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	blockchainDatabase.storeAddressToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	return kvWriteBatch
}
func (blockchainDatabase *BlockchainDatabase) fillBlockProperty(block *Model.Block) {
	transactionIndex := uint64(0)
	transactionHeight := blockchainDatabase.QueryBlockchainTransactionHeight()
	transactionOutputHeight := blockchainDatabase.QueryBlockchainTransactionOutputHeight()
	blockHeight := block.Height
	blockHash := block.Hash
	transactions := block.Transactions
	transactionCount := BlockTool.GetTransactionCount(block)
	block.TransactionCount = transactionCount
	block.PreviousTransactionHeight = transactionHeight
	for _, transaction := range transactions {
		transactionIndex = transactionIndex + 1
		transactionHeight = transactionHeight + 1
		transaction.BlockHeight = blockHeight

		transaction.TransactionIndex = transactionIndex
		transaction.TransactionHeight = transactionHeight

		outputs := transaction.Outputs
		for index, output := range outputs {
			transactionOutputHeight = transactionOutputHeight + 1

			output.BlockHeight = blockHeight
			output.BlockHash = blockHash
			output.TransactionHeight = transactionHeight
			output.TransactionHash = transaction.TransactionHash
			output.TransactionOutputIndex = uint64(index + 1)
			output.TransactionIndex = transaction.TransactionIndex
			output.TransactionOutputHeight = transactionOutputHeight
		}
	}
}

func (blockchainDatabase *BlockchainDatabase) storeHash(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHashKey := BlockchainDatabaseKeyTool.BuildHashKey(block.Hash)

	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHashKey, blockHashKey)
	} else {
		kvWriteBatch.Delete(blockHashKey)
	}
	transactions := block.Transactions
	for _, transaction := range transactions {
		transactionHashKey := BlockchainDatabaseKeyTool.BuildHashKey(transaction.TransactionHash)
		if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
			kvWriteBatch.Put(transactionHashKey, transactionHashKey)
		} else {
			kvWriteBatch.Delete(transactionHashKey)
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeAddress(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	for _, transaction := range transactions {
		outputs := transaction.Outputs
		for _, output := range outputs {
			addressKey := BlockchainDatabaseKeyTool.BuildAddressKey(output.Address)
			if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
				kvWriteBatch.Put(addressKey, addressKey)
			} else {
				kvWriteBatch.Delete(addressKey)
			}
		}

	}
}
func (blockchainDatabase *BlockchainDatabase) storeBlockchainHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockchainHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Long8ToByte8(block.Height))
	} else {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Long8ToByte8(block.Height-1))
	}
}
func (blockchainDatabase *BlockchainDatabase) storeBlockchainTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionCount := blockchainDatabase.QueryBlockchainTransactionHeight()
	bytesBlockchainTransactionCountKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Long8ToByte8(transactionCount+BlockTool.GetTransactionCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Long8ToByte8(transactionCount-BlockTool.GetTransactionCount(block)))
	}
}
func (blockchainDatabase *BlockchainDatabase) storeBlockchainTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionOutputCount := blockchainDatabase.QueryBlockchainTransactionOutputHeight()
	bytesBlockchainTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionOutputHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Long8ToByte8(transactionOutputCount+BlockTool.GetTransactionOutputCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Long8ToByte8(transactionOutputCount-BlockTool.GetTransactionOutputCount(block)))
	}
}
func (blockchainDatabase *BlockchainDatabase) storeBlockHeightToBlock(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(block.Height)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHeightKey, EncodeDecodeTool.EncodeBlock(block))
	} else {
		kvWriteBatch.Delete(blockHeightKey)
	}
}
func (blockchainDatabase *BlockchainDatabase) storeBlockHashToBlockHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHashBlockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(block.Hash)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHashBlockHeightKey, ByteUtil.Long8ToByte8(block.Height))
	} else {
		kvWriteBatch.Delete(blockHashBlockHeightKey)
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionHeightToTransaction(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHeightToTransactionKey := BlockchainDatabaseKeyTool.BuildTransactionHeightToTransactionKey(transaction.TransactionHeight)
			if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
				kvWriteBatch.Put(transactionHeightToTransactionKey, EncodeDecodeTool.EncodeTransaction(&transaction))
			} else {
				kvWriteBatch.Delete(transactionHeightToTransactionKey)
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionHashToTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHashToTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionHashToTransactionHeightKey(transaction.TransactionHash)
			if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
				kvWriteBatch.Put(transactionHashToTransactionHeightKey, ByteUtil.Long8ToByte8(transaction.TransactionHeight))
			} else {
				kvWriteBatch.Delete(transactionHashToTransactionHeightKey)
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputHeightToTransactionOutput(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputHeightToTransactionOutputKey := BlockchainDatabaseKeyTool.BuildTransactionOutputHeightToTransactionOutputKey(output.TransactionOutputHeight)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputHeightToTransactionOutputKey, EncodeDecodeTool.EncodeTransactionOutput(&output))
					} else {
						kvWriteBatch.Delete(transactionOutputHeightToTransactionOutputKey)
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToTransactionOutputHeightKey, ByteUtil.Long8ToByte8(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Delete(transactionOutputIdToUnspentTransactionOutputHeightKey)
					} else {
						kvWriteBatch.Put(transactionOutputIdToUnspentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(unspentTransactionOutput.TransactionOutputHeight))
					}
				}
			}
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToUnspentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToUnspentTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToSpentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToSpentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(unspentTransactionOutput.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToSpentTransactionOutputHeightKey)
					}
				}
			}
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToSpentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Delete(transactionOutputIdToSpentTransactionOutputHeightKey)
					} else {
						kvWriteBatch.Put(transactionOutputIdToSpentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(output.TransactionOutputHeight))
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, transactionOutput := range outputs {
					transactionOutputIdToToSourceTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSourceTransactionHeightKey(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToToSourceTransactionHeightKey, ByteUtil.Long8ToByte8(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToSourceTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToToDestinationTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToDestinationTransactionHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToToDestinationTransactionHeightKey, ByteUtil.Long8ToByte8(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToDestinationTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (blockchainDatabase *BlockchainDatabase) storeAddressToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions == nil {
		return
	}
	for _, transaction := range transactions {
		outputs := transaction.Outputs
		if outputs != nil {
			for _, transactionOutput := range outputs {
				addressToTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToTransactionOutputHeightKey(transactionOutput.Address)
				if blockchainActionEnum == BlockchainActionEnum.ADD_BLOCK {
					kvWriteBatch.Put(addressToTransactionOutputHeightKey, ByteUtil.Long8ToByte8(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (blockchainDatabase *BlockchainDatabase) storeAddressToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions == nil {
		return
	}
	for _, transaction := range transactions {
		inputs := transaction.Inputs
		if inputs != nil {
			for _, transactionInput := range inputs {
				utxo := transactionInput.UnspentTransactionOutput
				addressToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToUnspentTransactionOutputHeightKey(utxo.Address)
				if blockchainActionEnum == BlockchainActionEnum.ADD_BLOCK {
					kvWriteBatch.Delete(addressToUnspentTransactionOutputHeightKey)
				} else {
					kvWriteBatch.Put(addressToUnspentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(utxo.TransactionOutputHeight))
				}
			}
		}
		outputs := transaction.Outputs
		if outputs != nil {
			for _, transactionOutput := range outputs {
				addressToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToUnspentTransactionOutputHeightKey(transactionOutput.Address)
				if blockchainActionEnum == BlockchainActionEnum.ADD_BLOCK {
					kvWriteBatch.Put(addressToUnspentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToUnspentTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (blockchainDatabase *BlockchainDatabase) storeAddressToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions == nil {
		return
	}
	for _, transaction := range transactions {

		inputs := transaction.Inputs
		if inputs != nil {
			for _, transactionInput := range inputs {
				utxo := transactionInput.UnspentTransactionOutput
				addressToSpentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToSpentTransactionOutputHeightKey(utxo.Address)
				if blockchainActionEnum == BlockchainActionEnum.ADD_BLOCK {
					kvWriteBatch.Put(addressToSpentTransactionOutputHeightKey, ByteUtil.Long8ToByte8(utxo.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToSpentTransactionOutputHeightKey)
				}
			}
		}
	}
}
