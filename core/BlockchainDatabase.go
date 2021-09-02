package core

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/core/Model/BlockchainActionEnum"
	"helloworld-blockchain-go/core/Model/Script/BooleanCodeEnum"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworld-blockchain-go/core/tool/EncodeDecodeTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/SizeTool"
	"helloworld-blockchain-go/core/tool/StructureTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/setting/SystemVersionSettingTool"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"sync"
)

const BLOCKCHAIN_DATABASE_NAME = "BlockchainDatabase"

type BlockchainDatabase struct {
	Consensus         *Consensus
	Incentive         *Incentive
	VirtualMachine    *VirtualMachine
	CoreConfiguration *CoreConfiguration
}

var Mutex = sync.Mutex{}

func (b *BlockchainDatabase) AddBlockDto(blockDto *dto.BlockDto) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	block := BlockDto2Block(b, blockDto)
	checkBlock := b.CheckBlock(block)
	if !checkBlock {
		return false
	}
	kvWriteBatch := b.createBlockWriteBatch(block, BlockchainActionEnum.ADD_BLOCK)
	KvDbUtil.Write(b.getBlockchainDatabasePath(), kvWriteBatch)
	return true
}
func (b *BlockchainDatabase) DeleteTailBlock() {
	Mutex.Lock()
	defer Mutex.Unlock()
	tailBlock := b.QueryTailBlock()
	if tailBlock == nil {
		return
	}
	kvWriteBatch := b.createBlockWriteBatch(tailBlock, BlockchainActionEnum.DELETE_BLOCK)
	KvDbUtil.Write(b.getBlockchainDatabasePath(), kvWriteBatch)
}
func (b *BlockchainDatabase) DeleteBlocks(blockHeight uint64) {
	Mutex.Lock()
	defer Mutex.Unlock()
	for {
		tailBlock := b.QueryTailBlock()
		if tailBlock == nil {
			return
		}
		if tailBlock.Height < blockHeight {
			return
		}
		kvWriteBatch := b.createBlockWriteBatch(tailBlock, BlockchainActionEnum.DELETE_BLOCK)
		KvDbUtil.Write(b.getBlockchainDatabasePath(), kvWriteBatch)
	}
}

func (b *BlockchainDatabase) CheckBlock(block *Model.Block) bool {
	if !SystemVersionSettingTool.CheckSystemVersion(block.Height) {
		LogUtil.Debug("系统版本过低，不支持校验区块，请尽快升级系统。")
		return false
	}

	if !StructureTool.CheckBlockStructure(block) {
		LogUtil.Debug("区块数据异常，请校验区块的结构。")
		return false
	}
	//校验区块的大小
	if !SizeTool.CheckBlockSize(block) {
		LogUtil.Debug("区块数据异常，请校验区块的大小。")
		return false
	}

	//校验业务
	previousBlock := b.QueryTailBlock()
	//校验区块高度的连贯性
	if !BlockTool.CheckBlockHeight(previousBlock, block) {
		LogUtil.Debug("区块写入的区块高度出错。")
		return false
	}
	//校验区块的前区块哈希
	if !BlockTool.CheckPreviousBlockHash(previousBlock, block) {
		LogUtil.Debug("区块写入的前区块哈希出错。")
		return false
	}
	//校验区块时间
	if !BlockTool.CheckBlockTimestamp(previousBlock, block) {
		LogUtil.Debug("区块生成的时间太滞后。")
		return false
	}

	//校验新产生的哈希
	if !b.checkBlockNewHash(block) {
		LogUtil.Debug("区块数据异常，区块中新产生的哈希异常。")
		return false
	}
	//校验新产生的地址
	if !b.checkBlockNewAddress(block) {
		LogUtil.Debug("区块数据异常，区块中新产生的哈希异常。")
		return false
	}

	//校验双花
	if !b.checkBlockDoubleSpend(block) {
		LogUtil.Debug("区块数据异常，检测到双花攻击。")
		return false
	}
	//校验共识
	if !b.Consensus.CheckConsensus(b, block) {
		LogUtil.Debug("区块数据异常，未满足共识规则。")
		return false
	}
	//校验激励
	if !b.Incentive.CheckIncentive(b, block) {
		LogUtil.Debug("区块数据异常，激励校验失败。")
		return false
	}
	//从交易角度校验每一笔交易
	for _, transaction := range block.Transactions {
		transactionCanAddToNextBlock := b.CheckTransaction(transaction)
		if !transactionCanAddToNextBlock {
			LogUtil.Debug("区块数据异常，交易异常。")
			return false
		}
	}
	return true
}
func (b *BlockchainDatabase) CheckTransaction(transaction *Model.Transaction) bool {
	//校验交易的结构
	if !StructureTool.CheckTransactionStructure(transaction) {
		LogUtil.Debug("交易数据异常，请校验交易的结构。")
		return false
	}
	//校验交易的大小
	if !SizeTool.CheckTransactionSize(transaction) {
		LogUtil.Debug("交易数据异常，请校验交易的大小。")
		return false
	}

	//校验交易中的地址是否是P2PKH地址
	if !TransactionTool.CheckPayToPublicKeyHashAddress(transaction) {
		return false
	}
	//校验交易中的脚本是否是P2PKH脚本
	if !TransactionTool.CheckPayToPublicKeyHashScript(transaction) {
		return false
	}

	//业务校验
	//校验新产生的哈希
	if !b.checkTransactionNewHash(transaction) {
		LogUtil.Debug("区块数据异常，区块中新产生的哈希异常。")
		return false
	}
	//校验新产生的地址
	if !b.checkTransactionNewAddress(transaction) {
		LogUtil.Debug("区块数据异常，区块中新产生的哈希异常。")
		return false
	}
	//校验交易金额
	if !TransactionTool.CheckTransactionValue(transaction) {
		LogUtil.Debug("交易金额不合法")
		return false
	}
	//校验双花
	if !b.checkTransactionDoubleSpend(transaction) {
		LogUtil.Debug("交易数据异常，检测到双花攻击。")
		return false
	}
	//校验脚本
	if !b.checkTransactionScript(transaction) {
		LogUtil.Debug("交易校验失败：交易[输入脚本]解锁交易[输出脚本]异常。")
		return false
	}
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

func (b *BlockchainDatabase) QueryTailBlock() *Model.Block {
	blockchainHeight := b.QueryBlockchainHeight()
	if blockchainHeight <= GenesisBlockSetting.HEIGHT {
		return nil
	}
	return b.QueryBlockByBlockHeight(blockchainHeight)
}
func (b *BlockchainDatabase) QueryBlockByBlockHeight(blockHeight uint64) *Model.Block {
	bytesBlock := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(blockHeight))
	if bytesBlock == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToBlock(bytesBlock)
}
func (b *BlockchainDatabase) QueryBlockByBlockHash(blockHash string) *Model.Block {
	bytesBlockHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(blockHash))
	if bytesBlockHeight == nil {
		return nil
	}
	return b.QueryBlockByBlockHeight(ByteUtil.BytesToUint64(bytesBlockHeight))
}

func (b *BlockchainDatabase) QueryTransactionByTransactionHeight(transactionHeight uint64) *Model.Transaction {
	byteTransaction := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionHeightToTransactionKey(transactionHeight))
	if byteTransaction == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToTransaction(byteTransaction)
}
func (b *BlockchainDatabase) QueryTransactionByTransactionHash(transactionHash string) *Model.Transaction {
	transactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionHashToTransactionHeightKey(transactionHash))
	if transactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(transactionHeight))
}
func (b *BlockchainDatabase) QuerySourceTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.Transaction {
	sourceTransactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSourceTransactionHeightKey(transactionHash, transactionOutputIndex))
	if sourceTransactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(sourceTransactionHeight))
}
func (b *BlockchainDatabase) QueryDestinationTransactionByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.Transaction {
	destinationTransactionHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToDestinationTransactionHeightKey(transactionHash, transactionOutputIndex))
	if destinationTransactionHeight == nil {
		return nil
	}
	return b.QueryTransactionByTransactionHeight(ByteUtil.BytesToUint64(destinationTransactionHeight))
}

func (b *BlockchainDatabase) QueryTransactionOutputByTransactionOutputHeight(transactionOutputHeight uint64) *Model.TransactionOutput {
	bytesTransactionOutput := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputHeightToTransactionOutputKey(transactionOutputHeight))
	if bytesTransactionOutput == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToTransactionOutput(bytesTransactionOutput)
}
func (b *BlockchainDatabase) QueryTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))

}
func (b *BlockchainDatabase) QueryUnspentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QuerySpentTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *Model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(transactionHash, transactionOutputIndex))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}

func (b *BlockchainDatabase) QueryTransactionOutputByAddress(address string) *Model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressToTransactionOutputHeightKey(address))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QueryUnspentTransactionOutputByAddress(address string) *Model.TransactionOutput {
	bytesTransactionOutputHeight := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressToUnspentTransactionOutputHeightKey(address))
	if bytesTransactionOutputHeight == nil {
		return nil
	}
	return b.QueryTransactionOutputByTransactionOutputHeight(ByteUtil.BytesToUint64(bytesTransactionOutputHeight))
}
func (b *BlockchainDatabase) QuerySpentTransactionOutputByAddress(address string) *Model.TransactionOutput {
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
func (b *BlockchainDatabase) createBlockWriteBatch(block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) *KvDbUtil.KvWriteBatch {
	//b.fillBlockProperty(block)
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

/*func (b *BlockchainDatabase) fillBlockProperty(block *Model.Block) {
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
}*/

func (b *BlockchainDatabase) storeHash(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeAddress(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeBlockchainHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockchainHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height-1))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionCount := b.QueryBlockchainTransactionHeight()
	bytesBlockchainTransactionCountKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount+BlockTool.GetTransactionCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount-BlockTool.GetTransactionCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactionOutputCount := b.QueryBlockchainTransactionOutputHeight()
	bytesBlockchainTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionOutputHeightKey()
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount+BlockTool.GetTransactionOutputCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount-BlockTool.GetTransactionOutputCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockHeightToBlock(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(block.Height)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHeightKey, EncodeDecodeTool.EncodeBlock(block))
	} else {
		kvWriteBatch.Delete(blockHeightKey)
	}
}
func (b *BlockchainDatabase) storeBlockHashToBlockHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	blockHashBlockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(block.Hash)
	if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
		kvWriteBatch.Put(blockHashBlockHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Delete(blockHashBlockHeightKey)
	}
}
func (b *BlockchainDatabase) storeTransactionHeightToTransaction(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHeightToTransactionKey := BlockchainDatabaseKeyTool.BuildTransactionHeightToTransactionKey(transaction.TransactionHeight)
			if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
				kvWriteBatch.Put(transactionHeightToTransactionKey, EncodeDecodeTool.EncodeTransaction(transaction))
			} else {
				kvWriteBatch.Delete(transactionHeightToTransactionKey)
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionHashToTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputHeightToTransactionOutput(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputHeightToTransactionOutputKey := BlockchainDatabaseKeyTool.BuildTransactionOutputHeightToTransactionOutputKey(output.TransactionOutputHeight)
					if BlockchainActionEnum.ADD_BLOCK == blockchainActionEnum {
						kvWriteBatch.Put(transactionOutputHeightToTransactionOutputKey, EncodeDecodeTool.EncodeTransactionOutput(output))
					} else {
						kvWriteBatch.Delete(transactionOutputHeightToTransactionOutputKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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
func (b *BlockchainDatabase) storeAddressToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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

func (b *BlockchainDatabase) storeAddressToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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

func (b *BlockchainDatabase) storeAddressToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *Model.Block, blockchainActionEnum BlockchainActionEnum.BlockchainActionEnum) {
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

/**
 * 校验区块新产生的哈希
 */
func (b *BlockchainDatabase) checkBlockNewHash(block *Model.Block) bool {
	//校验哈希作为主键的正确性
	//新产生的哈希不能有重复
	if BlockTool.IsExistDuplicateNewHash(block) {
		LogUtil.Debug("区块数据异常，区块中新产生的哈希有重复。")
		return false
	}

	//新产生的哈希不能被区块链使用过了
	//校验区块Hash是否已经被使用了
	blockHash := block.Hash
	if b.isHashUsed(blockHash) {
		LogUtil.Debug("区块数据异常，区块Hash已经被使用了。")
		return false
	}
	//校验每一笔交易新产生的Hash是否正确
	blockTransactions := block.Transactions
	if blockTransactions != nil {
		for _, transaction := range blockTransactions {
			if !b.checkTransactionNewHash(transaction) {
				return false
			}
		}
	}
	return true
}

/**
 * 区块中校验新产生的哈希
 */
func (b *BlockchainDatabase) checkTransactionNewHash(transaction *Model.Transaction) bool {
	//校验哈希作为主键的正确性
	//校验交易Hash是否已经被使用了
	transactionHash := transaction.TransactionHash
	if b.isHashUsed(transactionHash) {
		LogUtil.Debug("交易数据异常，交易Hash已经被使用了。")
		return false
	}
	return true
}

/**
 * 哈希是否已经被区块链系统使用了？
 */
func (b *BlockchainDatabase) isHashUsed(hash string) bool {
	bytesHash := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildHashKey(hash))
	return bytesHash != nil
}

/**
 * 校验区块新产生的地址
 */
func (b *BlockchainDatabase) checkBlockNewAddress(block *Model.Block) bool {
	//校验地址作为主键的正确性
	//新产生的地址不能有重复
	if BlockTool.IsExistDuplicateNewAddress(block) {
		LogUtil.Debug("区块数据异常，区块中新产生的地址有重复。")
		return false
	}
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			if !b.checkTransactionNewAddress(transaction) {
				return false
			}
		}
	}
	return true
}

/**
 * 区块中校验新产生的哈希
 */
func (b *BlockchainDatabase) checkTransactionNewAddress(transaction *Model.Transaction) bool {
	//区块新产生的地址不能有重复
	if TransactionTool.IsExistDuplicateNewAddress(transaction) {
		return false
	}
	//区块新产生的地址不能被使用过了
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			address := output.Address
			if b.isAddressUsed(address) {
				LogUtil.Debug("区块数据异常，地址[" + address + "]重复。")
				return false
			}
		}
	}
	return true
}
func (b *BlockchainDatabase) isAddressUsed(address string) bool {
	bytesAddress := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildAddressKey(address))
	return bytesAddress != nil
}

//region 双花攻击
/**
 * 校验双花
 * 双花指的是同一笔UTXO被花费两次或多次。
 */
func (b *BlockchainDatabase) checkBlockDoubleSpend(block *Model.Block) bool {
	//双花交易：区块内部存在重复的[未花费交易输出]
	if BlockTool.IsExistDuplicateUtxo(block) {
		LogUtil.Debug("区块数据异常：发生双花交易。")
		return false
	}
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			if !b.checkTransactionDoubleSpend(transaction) {
				LogUtil.Debug("区块数据异常：发生双花交易。")
				return false
			}
		}
	}
	return true
}

/**
 * 校验双花
 */
func (b *BlockchainDatabase) checkTransactionDoubleSpend(transaction *Model.Transaction) bool {
	//双花交易：交易内部存在重复的[未花费交易输出]
	if TransactionTool.IsExistDuplicateUtxo(transaction) {
		LogUtil.Debug("交易数据异常，检测到双花攻击。")
		return false
	}
	//双花交易：交易内部使用了[已经花费的[未花费交易输出]]
	if !b.checkStxoIsUtxo(transaction) {
		LogUtil.Debug("交易数据异常：发生双花交易。")
		return false
	}
	return true
}

/**
 * 检查[花费的交易输出]是否都是[未花费的交易输出]
 */
func (b *BlockchainDatabase) checkStxoIsUtxo(transaction *Model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil {
		for _, transactionInput := range inputs {
			unspentTransactionOutput := transactionInput.UnspentTransactionOutput
			transactionOutput := b.QueryUnspentTransactionOutputByTransactionOutputId(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
			if transactionOutput == nil {
				LogUtil.Debug("交易数据异常：交易输入不是未花费交易输出。")
				return false
			}
		}
	}
	return true
}

//endregion
/**
 * 检验交易脚本，即校验交易输入能解锁交易输出吗？即用户花费的是自己的钱吗？
 * 校验用户花费的是自己的钱吗，用户只可以花费自己的钱。专业点的说法，校验UTXO所有权，用户只可以花费自己拥有的UTXO。
 * 用户如何能证明自己拥有这个UTXO，只要用户能创建出一个能解锁(该UTXO对应的交易输出脚本)的交易输入脚本，就证明了用户拥有该UTXO。
 * 这是因为锁(交易输出脚本)是用户创建的，自然只有该用户有对应的钥匙(交易输入脚本)，自然意味着有钥匙的用户拥有这把锁的所有权。
 */
func (b *BlockchainDatabase) checkTransactionScript(transaction *Model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil && len(inputs) != 0 {
		for _, transactionInput := range inputs {
			//锁(交易输出脚本)
			outputScript := transactionInput.UnspentTransactionOutput.OutputScript
			//钥匙(交易输入脚本)
			inputScript := transactionInput.InputScript
			//完整脚本
			script := ScriptTool.CreateScript(inputScript, outputScript)
			//执行脚本
			scriptExecuteResult := b.VirtualMachine.ExecuteScript(transaction, script)
			/*fmt.Println(ByteUtil.HexStringToBytes(*scriptExecuteResult.Pop()))
			fmt.Println(BooleanCodeEnum.TRUE.Code)*/

			//脚本执行结果是个栈，如果栈有且只有一个元素，且这个元素是0x01，则解锁成功。
			//executeSuccess := scriptExecuteResult.Size() == 1 && ByteUtil.IsEquals(BooleanCodeEnum.TRUE.Code, ByteUtil.HexStringToBytes(*scriptExecuteResult.Pop()))
			//TODO java...
			executeSuccess := scriptExecuteResult.Size() == 1 && StringUtil.IsEquals(ByteUtil.BytesToHexString(BooleanCodeEnum.TRUE.Code), *scriptExecuteResult.Pop())
			if !executeSuccess {
				return false
			}
		}
	}
	return true
}
