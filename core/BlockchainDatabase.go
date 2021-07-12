package core

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/BlockchainActionEnum"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworldcoin-go/core/tool/EncodeDecodeTool"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
	"sync"
)

const BLOCKCHAIN_DATABASE_NAME = "BlockchainDatabase"

type BlockchainDatabase struct {
	Consensus         *Consensus
	Incentive         *Incentive
	CoreConfiguration *CoreConfiguration
}

func (b *BlockchainDatabase) AddBlockDto(blockDto *dto.BlockDto) bool {
	var lock = sync.Mutex{}
	lock.Lock()
	block := BlockDto2Block(b, blockDto)
	checkBlock := b.CheckBlock(block)
	if !checkBlock {
		return false
	}
	kvWriteBatch := b.createBlockWriteBatch(block, BlockchainActionEnum.ADD_BLOCK)
	KvDbUtil.Write(b.getBlockchainDatabasePath(), kvWriteBatch)
	lock.Unlock()
	return true
}
func (b *BlockchainDatabase) DeleteTailBlock() {

}
func (b *BlockchainDatabase) DeleteBlocks(blockHeight uint64) {
}

func (b *BlockchainDatabase) CheckBlock(block *model.Block) bool {
	return true
}
func (b *BlockchainDatabase) CheckTransaction(block *model.Transaction) bool {
	return true
}

func (b *BlockchainDatabase) QueryBlockchainHeight() uint64 {
	bytesBlockchainHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockchainHeightKey())
	if bytesBlockchainHeight == nil {
		return GenesisBlockSetting.HEIGHT
	}
	return ByteUtil.BytesToUint64(bytesBlockchainHeight)
}
func (b *BlockchainDatabase) QueryBlockchainTransactionHeight() uint64 {
	byteTotalTransactionCount := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockchainTransactionHeightKey())
	if byteTotalTransactionCount == nil {
		return uint64(0)
	}
	return ByteUtil.BytesToUint64(byteTotalTransactionCount)
}
func (b *BlockchainDatabase) QueryBlockchainTransactionOutputHeight() uint64 {
	byteTotalTransactionCount := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockchainTransactionOutputHeightKey())
	if byteTotalTransactionCount == nil {
		return uint64(0)
	}
	return ByteUtil.BytesToUint64(byteTotalTransactionCount)
}

func (b *BlockchainDatabase) QueryTailBlock() *model.Block {
	blockchainHeight := b.QueryBlockchainHeight()
	if blockchainHeight <= GenesisBlockSetting.HEIGHT {
		return nil
	}
	return b.QueryBlockByBlockHeight(blockchainHeight)
}
func (b *BlockchainDatabase) QueryBlockByBlockHeight(blockHeight uint64) *model.Block {
	bytesBlock := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(blockHeight))
	if bytesBlock == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToBlock(bytesBlock)
}
func (b *BlockchainDatabase) QueryBlockByBlockHash(blockHash string) *model.Block {
	bytesBlockHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(blockHash))
	if bytesBlockHeight == nil {
		return nil
	}
	return b.QueryBlockByBlockHeight(ByteUtil.BytesToUint64(bytesBlockHeight))
}

func (b *BlockchainDatabase) QueryTransactionByTransactionHeight(transactionHeight uint64) *model.Transaction {
	byteTransaction := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionHeightToTransactionKey(transactionHeight))
	if byteTransaction == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToTransaction(byteTransaction)
}
func (b *BlockchainDatabase) QueryTransactionByTransactionHash(transactionHash string) *model.Transaction {
	transactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionHashToTransactionHeightKey(transactionHash))
	if transactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(transactionHeight))
}
func (b *BlockchainDatabase) QuerySourceTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *model.Transaction {
	sourceTransactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSourceTransactionHeightKey(transactionHash, transactionOutputIndex))
	if sourceTransactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(sourceTransactionHeight))
}
func (b *BlockchainDatabase) QueryDestinationTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *model.Transaction {
	destinationTransactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToDestinationTransactionHeightKey(transactionHash, transactionOutputIndex))
	if destinationTransactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(destinationTransactionHeight))
}

func (b *BlockchainDatabase) QueryTransactionOutputByTransactionOutputHeight(transactionOutputHeight uint64) *model.TransactionOutput {
	bytesTransactionOutput := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputHeightToTransactionOutputKey(transactionOutputHeight))
	if bytesTransactionOutput == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToTransactionOutput(bytesTransactionOutput)
}
func (b *BlockchainDatabase) QueryTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))

}
func (b *BlockchainDatabase) QueryUnspentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QuerySpentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}

func (b *BlockchainDatabase) QueryTransactionOutputByAddress(address string) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressToTransactionOutputHeightKey(address))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QueryUnspentTransactionOutputByAddress(address string) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressToUnspentTransactionOutputHeightKey(address))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QuerySpentTransactionOutputByAddress(address string) *model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressToSpentTransactionOutputHeightKey(address))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}

func (b *BlockchainDatabase) GetIncentive() *Incentive {
	return b.Incentive
}
func (b *BlockchainDatabase) GetConsensus() *Consensus {
	return b.Consensus
}

func (b *BlockchainDatabase) getBlockchainDatabasePath() string {
	return FileUtil.NewPath(b.CoreConfiguration.getCorePath(), BLOCKCHAIN_DATABASE_NAME)
}
func (b *BlockchainDatabase) createBlockWriteBatch(block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) *KvDbUtil.KvWriteBatch {
	b.fillBlockProperty(block)
	kvWriteBatch := new(KvDbUtil.KvWriteBatch)

	b.storeHash(kvWriteBatch, block, blockchainActionEnum)
	b.storeAddress(kvWriteBatch, block, blockchainActionEnum)

	b.storeBlockchainHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeBlockchainTransactionHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeBlockchainTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)

	b.storeBlockHeightToBlock(kvWriteBatch, block, blockchainActionEnum)
	b.storeBlockHashToBlockHeight(kvWriteBatch, block, blockchainActionEnum)

	b.storeTransactionHeightToTransaction(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionHashToTransactionHeight(kvWriteBatch, block, blockchainActionEnum)

	b.storeTransactionOutputHeightToTransactionOutput(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch, block, blockchainActionEnum)

	b.storeAddressToTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeAddressToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	b.storeAddressToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainActionEnum)
	return kvWriteBatch
}
func (b *BlockchainDatabase) fillBlockProperty(block *model.Block) {
	transactionIndex := uint64(0)
	transactionHeight := b.QueryBlockchainTransactionHeight()
	transactionOutputHeight := b.QueryBlockchainTransactionOutputHeight()
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

func (b *BlockchainDatabase) storeHash(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeAddress(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeBlockchainHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockchainHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height-1))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionCount := b.QueryBlockchainTransactionHeight()
	bytesBlockchainTransactionCountKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount+BlockTool.GetTransactionCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount-BlockTool.GetTransactionCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionOutputCount := b.QueryBlockchainTransactionOutputHeight()
	bytesBlockchainTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionOutputHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount+BlockTool.GetTransactionOutputCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount-BlockTool.GetTransactionOutputCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockHeightToBlock(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(block.Height)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHeightKey, EncodeDecodeTool.EncodeBlock(block))
	} else {
		kvWriteBatch.Delete(blockHeightKey)
	}
}
func (b *BlockchainDatabase) storeBlockHashToBlockHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHashBlockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(block.Hash)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHashBlockHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Delete(blockHashBlockHeightKey)
	}
}
func (b *BlockchainDatabase) storeTransactionHeightToTransaction(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionHashToTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHashToTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionHashToTransactionHeightKey(transaction.TransactionHash)
			if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
				kvWriteBatch.Put(transactionHashToTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
			} else {
				kvWriteBatch.Delete(transactionHashToTransactionHeightKey)
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputHeightToTransactionOutput(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
						kvWriteBatch.Put(transactionOutputIdToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(unspentTransactionOutput.TransactionOutputHeight))
					}
				}
			}
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToUnspentTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToSpentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToSpentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(unspentTransactionOutput.TransactionOutputHeight))
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
						kvWriteBatch.Put(transactionOutputIdToSpentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, transactionOutput := range outputs {
					transactionOutputIdToToSourceTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSourceTransactionHeightKey(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToToSourceTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToSourceTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToToDestinationTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToDestinationTransactionHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputIdToToDestinationTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToDestinationTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeAddressToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
					kvWriteBatch.Put(addressToTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (b *BlockchainDatabase) storeAddressToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
					kvWriteBatch.Put(addressToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(utxo.TransactionOutputHeight))
				}
			}
		}
		outputs := transaction.Outputs
		if outputs != nil {
			for _, transactionOutput := range outputs {
				addressToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToUnspentTransactionOutputHeightKey(transactionOutput.Address)
				if blockchainActionEnum == BlockchainActionEnum.ADD_BLOCK {
					kvWriteBatch.Put(addressToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToUnspentTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (b *BlockchainDatabase) storeAddressToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
					kvWriteBatch.Put(addressToSpentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(utxo.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToSpentTransactionOutputHeightKey)
				}
			}
		}
	}
}
