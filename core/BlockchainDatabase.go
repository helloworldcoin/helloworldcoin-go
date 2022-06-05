package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/BlockchainAction"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/model/script/BooleanCode"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/BlockchainDatabaseKeyTool"
	"helloworldcoin-go/core/tool/ScriptDtoTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/SizeTool"
	"helloworldcoin-go/core/tool/StructureTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/EncodeDecodeTool"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/StringUtil"
	"sync"
)

const BLOCKCHAIN_DATABASE_NAME = "BlockchainDatabase"

type BlockchainDatabase struct {
	consensus         *Consensus
	incentive         *Incentive
	virtualMachine    *VirtualMachine
	coreConfiguration *CoreConfiguration
}

func NewBlockchainDatabase(consensus *Consensus, incentive *Incentive, virtualMachine *VirtualMachine, coreConfiguration *CoreConfiguration) *BlockchainDatabase {
	var blockchainDatabase BlockchainDatabase
	blockchainDatabase.consensus = consensus
	blockchainDatabase.incentive = incentive
	blockchainDatabase.virtualMachine = virtualMachine
	blockchainDatabase.coreConfiguration = coreConfiguration
	return &blockchainDatabase
}
func (b *BlockchainDatabase) GetIncentive() *Incentive {
	return b.incentive
}
func (b *BlockchainDatabase) GetConsensus() *Consensus {
	return b.consensus
}

var Mutex = sync.Mutex{}

func (b *BlockchainDatabase) AddBlockDto(blockDto *dto.BlockDto) bool {
	Mutex.Lock()
	defer Mutex.Unlock()

	block := b.BlockDto2Block(blockDto)
	checkBlock := b.CheckBlock(block)
	if !checkBlock {
		return false
	}
	kvWriteBatch := b.createBlockWriteBatch(block, BlockchainAction.ADD_BLOCK)
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
	kvWriteBatch := b.createBlockWriteBatch(tailBlock, BlockchainAction.DELETE_BLOCK)
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
		kvWriteBatch := b.createBlockWriteBatch(tailBlock, BlockchainAction.DELETE_BLOCK)
		KvDbUtil.Write(b.getBlockchainDatabasePath(), kvWriteBatch)
	}
}

func (b *BlockchainDatabase) CheckBlock(block *model.Block) bool {

	//check block structure
	if !StructureTool.CheckBlockStructure(block) {
		LogUtil.Debug("The block data is abnormal. Please verify the block structure.")
		return false
	}
	//check block size
	if !SizeTool.CheckBlockSize(block) {
		LogUtil.Debug("The block data is abnormal, please check the size of the block.")
		return false
	}

	//check business
	previousBlock := b.QueryTailBlock()
	//check block height
	if !BlockTool.CheckBlockHeight(previousBlock, block) {
		LogUtil.Debug("Wrong block height for block write.")
		return false
	}
	//check previous block hash
	if !BlockTool.CheckPreviousBlockHash(previousBlock, block) {
		LogUtil.Debug("The previous block hash of the block write was wrong.")
		return false
	}
	//check block timestamp
	if !BlockTool.CheckBlockTimestamp(previousBlock, block) {
		LogUtil.Debug("Block generation is too late.")
		return false
	}

	//check block new hash
	if !b.checkBlockNewHash(block) {
		LogUtil.Debug("The block data is abnormal, and the newly generated hash in the block is abnormal.")
		return false
	}
	//check block new address
	if !b.checkBlockNewAddress(block) {
		LogUtil.Debug("The block data is abnormal, and the newly generated hash in the block is abnormal.")
		return false
	}

	//check block double spend
	if !b.checkBlockDoubleSpend(block) {
		LogUtil.Debug("The block data is abnormal, and a double-spending attack is detected.")
		return false
	}
	//check consensus
	if !b.consensus.CheckConsensus(b, block) {
		LogUtil.Debug("The block data is abnormal and the consensus rules are not met.")
		return false
	}
	//check incentive
	if !b.incentive.CheckIncentive(b, block) {
		LogUtil.Debug("The block data is abnormal, and the incentive verification fails.")
		return false
	}
	//check transaction
	for _, transaction := range block.Transactions {
		transactionCanAddToNextBlock := b.CheckTransaction(transaction)
		if !transactionCanAddToNextBlock {
			LogUtil.Debug("The block data is abnormal, and the transaction is abnormal.")
			return false
		}
	}
	return true
}
func (b *BlockchainDatabase) CheckTransaction(transaction *model.Transaction) bool {
	//check Transaction Structure
	if !StructureTool.CheckTransactionStructure(transaction) {
		LogUtil.Debug("The transaction data is abnormal, please verify the structure of the transaction.")
		return false
	}
	//check Transaction Size
	if !SizeTool.CheckTransactionSize(transaction) {
		LogUtil.Debug("The transaction data is abnormal, please check the size of the transaction.")
		return false
	}

	//Check if the address in the transaction is a P2PKH address
	if !TransactionTool.CheckPayToPublicKeyHashAddress(transaction) {
		return false
	}
	//Check if the address in the transaction is a P2PKH address
	if !TransactionTool.CheckPayToPublicKeyHashScript(transaction) {
		return false
	}

	//business verification
	//check Transaction New Hash
	if !b.checkTransactionNewHash(transaction) {
		LogUtil.Debug("The block data is abnormal, and the newly generated hash in the block is abnormal.")
		return false
	}
	//check Transaction New Address
	if !b.checkTransactionNewAddress(transaction) {
		LogUtil.Debug("The block data is abnormal, and the newly generated hash in the block is abnormal.")
		return false
	}
	//check Transaction Value
	if !TransactionTool.CheckTransactionValue(transaction) {
		LogUtil.Debug("The block data is abnormal and the transaction amount is illegal")
		return false
	}
	//check Transaction Double Spend
	if !b.checkTransactionDoubleSpend(transaction) {
		LogUtil.Debug("The transaction data is abnormal, and a double-spending attack is detected.")
		return false
	}
	//check Transaction Script
	if !b.checkTransactionScript(transaction) {
		LogUtil.Debug("Transaction verification failed: transaction [input script] unlock transaction [output script] exception.")
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
	return EncodeDecodeTool.Decode(bytesBlock, model.Block{}).(*model.Block)
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
	return EncodeDecodeTool.Decode(byteTransaction, model.Transaction{}).(*model.Transaction)
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
	return EncodeDecodeTool.Decode(bytesTransactionOutput, model.TransactionOutput{}).(*model.TransactionOutput)
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

func (b *BlockchainDatabase) getBlockchainDatabasePath() string {
	return FileUtil.NewPath(b.coreConfiguration.getCorePath(), BLOCKCHAIN_DATABASE_NAME)
}
func (b *BlockchainDatabase) createBlockWriteBatch(block *model.Block, blockchainAction BlockchainAction.BlockchainAction) *KvDbUtil.KvWriteBatch {
	//b.fillBlockProperty(block)
	kvWriteBatch := new(KvDbUtil.KvWriteBatch)

	b.storeHash(kvWriteBatch, block, blockchainAction)
	b.storeAddress(kvWriteBatch, block, blockchainAction)

	b.storeBlockchainHeight(kvWriteBatch, block, blockchainAction)
	b.storeBlockchainTransactionHeight(kvWriteBatch, block, blockchainAction)
	b.storeBlockchainTransactionOutputHeight(kvWriteBatch, block, blockchainAction)

	b.storeBlockHeightToBlock(kvWriteBatch, block, blockchainAction)
	b.storeBlockHashToBlockHeight(kvWriteBatch, block, blockchainAction)

	b.storeTransactionHeightToTransaction(kvWriteBatch, block, blockchainAction)
	b.storeTransactionHashToTransactionHeight(kvWriteBatch, block, blockchainAction)

	b.storeTransactionOutputHeightToTransactionOutput(kvWriteBatch, block, blockchainAction)
	b.storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	b.storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	b.storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	b.storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch, block, blockchainAction)
	b.storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch, block, blockchainAction)

	b.storeAddressToTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	b.storeAddressToUnspentTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	b.storeAddressToSpentTransactionOutputHeight(kvWriteBatch, block, blockchainAction)
	return kvWriteBatch
}

/*func (b *BlockchainDatabase) fillBlockProperty(block *model.Block) {
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

func (b *BlockchainDatabase) storeHash(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	blockHashKey := BlockchainDatabaseKeyTool.BuildHashKey(block.Hash)

	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(blockHashKey, blockHashKey)
	} else {
		kvWriteBatch.Delete(blockHashKey)
	}
	transactions := block.Transactions
	for _, transaction := range transactions {
		transactionHashKey := BlockchainDatabaseKeyTool.BuildHashKey(transaction.TransactionHash)
		if BlockchainAction.ADD_BLOCK == blockchainAction {
			kvWriteBatch.Put(transactionHashKey, transactionHashKey)
		} else {
			kvWriteBatch.Delete(transactionHashKey)
		}
	}
}
func (b *BlockchainDatabase) storeAddress(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	for _, transaction := range transactions {
		outputs := transaction.Outputs
		for _, output := range outputs {
			addressKey := BlockchainDatabaseKeyTool.BuildAddressKey(output.Address)
			if BlockchainAction.ADD_BLOCK == blockchainAction {
				kvWriteBatch.Put(addressKey, addressKey)
			} else {
				kvWriteBatch.Delete(addressKey)
			}
		}

	}
}
func (b *BlockchainDatabase) storeBlockchainHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	blockchainHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainHeightKey()
	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Put(blockchainHeightKey, ByteUtil.Uint64ToBytes(block.Height-1))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactionCount := b.QueryBlockchainTransactionHeight()
	bytesBlockchainTransactionCountKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionHeightKey()
	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount+BlockTool.GetTransactionCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionCountKey, ByteUtil.Uint64ToBytes(transactionCount-BlockTool.GetTransactionCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockchainTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactionOutputCount := b.QueryBlockchainTransactionOutputHeight()
	bytesBlockchainTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildBlockchainTransactionOutputHeightKey()
	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount+BlockTool.GetTransactionOutputCount(block)))
	} else {
		kvWriteBatch.Put(bytesBlockchainTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutputCount-BlockTool.GetTransactionOutputCount(block)))
	}
}
func (b *BlockchainDatabase) storeBlockHeightToBlock(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	blockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHeightToBlockKey(block.Height)
	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(blockHeightKey, EncodeDecodeTool.Encode(block))
	} else {
		kvWriteBatch.Delete(blockHeightKey)
	}
}
func (b *BlockchainDatabase) storeBlockHashToBlockHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	blockHashBlockHeightKey := BlockchainDatabaseKeyTool.BuildBlockHashToBlockHeightKey(block.Hash)
	if BlockchainAction.ADD_BLOCK == blockchainAction {
		kvWriteBatch.Put(blockHashBlockHeightKey, ByteUtil.Uint64ToBytes(block.Height))
	} else {
		kvWriteBatch.Delete(blockHashBlockHeightKey)
	}
}
func (b *BlockchainDatabase) storeTransactionHeightToTransaction(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHeightToTransactionKey := BlockchainDatabaseKeyTool.BuildTransactionHeightToTransactionKey(transaction.TransactionHeight)
			if BlockchainAction.ADD_BLOCK == blockchainAction {
				kvWriteBatch.Put(transactionHeightToTransactionKey, EncodeDecodeTool.Encode(transaction))
			} else {
				kvWriteBatch.Delete(transactionHeightToTransactionKey)
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionHashToTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			transactionHashToTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionHashToTransactionHeightKey(transaction.TransactionHash)
			if BlockchainAction.ADD_BLOCK == blockchainAction {
				kvWriteBatch.Put(transactionHashToTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
			} else {
				kvWriteBatch.Delete(transactionHashToTransactionHeightKey)
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputHeightToTransactionOutput(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputHeightToTransactionOutputKey := BlockchainDatabaseKeyTool.BuildTransactionOutputHeightToTransactionOutputKey(output.TransactionOutputHeight)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Put(transactionOutputHeightToTransactionOutputKey, EncodeDecodeTool.Encode(output))
					} else {
						kvWriteBatch.Delete(transactionOutputHeightToTransactionOutputKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, output := range outputs {
					transactionOutputIdToTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToTransactionOutputHeightKey(output.TransactionHash, output.TransactionOutputIndex)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Put(transactionOutputIdToTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToUnspentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
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
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Put(transactionOutputIdToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToUnspentTransactionOutputHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToSpentTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSpentTransactionOutputHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
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
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Delete(transactionOutputIdToSpentTransactionOutputHeightKey)
					} else {
						kvWriteBatch.Put(transactionOutputIdToSpentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(output.TransactionOutputHeight))
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToSourceTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			outputs := transaction.Outputs
			if outputs != nil {
				for _, transactionOutput := range outputs {
					transactionOutputIdToToSourceTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToSourceTransactionHeightKey(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Put(transactionOutputIdToToSourceTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToSourceTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeTransactionOutputIdToDestinationTransactionHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			inputs := transaction.Inputs
			if inputs != nil {
				for _, transactionInput := range inputs {
					unspentTransactionOutput := transactionInput.UnspentTransactionOutput
					transactionOutputIdToToDestinationTransactionHeightKey := BlockchainDatabaseKeyTool.BuildTransactionOutputIdToDestinationTransactionHeightKey(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
					if BlockchainAction.ADD_BLOCK == blockchainAction {
						kvWriteBatch.Put(transactionOutputIdToToDestinationTransactionHeightKey, ByteUtil.Uint64ToBytes(transaction.TransactionHeight))
					} else {
						kvWriteBatch.Delete(transactionOutputIdToToDestinationTransactionHeightKey)
					}
				}
			}
		}
	}
}
func (b *BlockchainDatabase) storeAddressToTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
	transactions := block.Transactions
	if transactions == nil {
		return
	}
	for _, transaction := range transactions {
		outputs := transaction.Outputs
		if outputs != nil {
			for _, transactionOutput := range outputs {
				addressToTransactionOutputHeightKey := BlockchainDatabaseKeyTool.BuildAddressToTransactionOutputHeightKey(transactionOutput.Address)
				if blockchainAction == BlockchainAction.ADD_BLOCK {
					kvWriteBatch.Put(addressToTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (b *BlockchainDatabase) storeAddressToUnspentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
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
				if blockchainAction == BlockchainAction.ADD_BLOCK {
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
				if blockchainAction == BlockchainAction.ADD_BLOCK {
					kvWriteBatch.Put(addressToUnspentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(transactionOutput.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToUnspentTransactionOutputHeightKey)
				}
			}
		}
	}
}

func (b *BlockchainDatabase) storeAddressToSpentTransactionOutputHeight(kvWriteBatch *KvDbUtil.KvWriteBatch, block *model.Block, blockchainAction BlockchainAction.BlockchainAction) {
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
				if blockchainAction == BlockchainAction.ADD_BLOCK {
					kvWriteBatch.Put(addressToSpentTransactionOutputHeightKey, ByteUtil.Uint64ToBytes(utxo.TransactionOutputHeight))
				} else {
					kvWriteBatch.Delete(addressToSpentTransactionOutputHeightKey)
				}
			}
		}
	}
}

//region block hash related
/**
 * check Block New Hash
 */
func (b *BlockchainDatabase) checkBlockNewHash(block *model.Block) bool {
	if BlockTool.IsExistDuplicateNewHash(block) {
		LogUtil.Debug("The block data is abnormal, exist duplicate hash.")
		return false
	}

	blockHash := block.Hash
	if b.isHashUsed(blockHash) {
		LogUtil.Debug("The block data is abnormal, and the block hash has been used.")
		return false
	}
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
func (b *BlockchainDatabase) checkTransactionNewHash(transaction *model.Transaction) bool {
	transactionHash := transaction.TransactionHash
	if b.isHashUsed(transactionHash) {
		return false
	}
	return true
}
func (b *BlockchainDatabase) isHashUsed(hash string) bool {
	bytesHash := KvDbUtil.Get(b.getBlockchainDatabasePath(), BlockchainDatabaseKeyTool.BuildHashKey(hash))
	return bytesHash != nil
}

//endregion

//region address related
func (b *BlockchainDatabase) checkBlockNewAddress(block *model.Block) bool {
	if BlockTool.IsExistDuplicateNewAddress(block) {
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
func (b *BlockchainDatabase) checkTransactionNewAddress(transaction *model.Transaction) bool {
	if TransactionTool.IsExistDuplicateNewAddress(transaction) {
		return false
	}
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			address := output.Address
			if b.isAddressUsed(address) {
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

//endregion

//region Double spend attack
func (b *BlockchainDatabase) checkBlockDoubleSpend(block *model.Block) bool {
	if BlockTool.IsExistDuplicateUtxo(block) {
		LogUtil.Debug("Abnormal block data: a double-spend transaction occurred.")
		return false
	}
	transactions := block.Transactions
	if transactions != nil {
		for _, transaction := range transactions {
			if !b.checkTransactionDoubleSpend(transaction) {
				LogUtil.Debug("Abnormal block data: a double-spend transaction occurred.")
				return false
			}
		}
	}
	return true
}

/**
 * check Transaction Double Spend
 */
func (b *BlockchainDatabase) checkTransactionDoubleSpend(transaction *model.Transaction) bool {
	//Double spend transaction: there is a duplicate [unspent transaction output] inside the transaction
	if TransactionTool.IsExistDuplicateUtxo(transaction) {
		LogUtil.Debug("The transaction data is abnormal, and a double-spending attack is detected.")
		return false
	}
	//Double spend transaction: transaction uses [spent [unspent transaction output]] inside the transaction
	if !b.checkStxoIsUtxo(transaction) {
		LogUtil.Debug("The transaction data is abnormal, and a double-spending attack is detected.")
		return false
	}
	return true
}

/**
 * Check if [spent transaction outputs] are all [unspent transaction outputs] ?
 */
func (b *BlockchainDatabase) checkStxoIsUtxo(transaction *model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil {
		for _, transactionInput := range inputs {
			unspentTransactionOutput := transactionInput.UnspentTransactionOutput
			transactionOutput := b.QueryUnspentTransactionOutputByTransactionOutputId(unspentTransactionOutput.TransactionHash, unspentTransactionOutput.TransactionOutputIndex)
			if transactionOutput == nil {
				LogUtil.Debug("Transaction data exception: transaction input is not unspent transaction output.")
				return false
			}
		}
	}
	return true
}

//endregion

func (b *BlockchainDatabase) checkTransactionScript(transaction *model.Transaction) bool {
	inputs := transaction.Inputs
	if inputs != nil && len(inputs) != 0 {
		for _, transactionInput := range inputs {
			outputScript := transactionInput.UnspentTransactionOutput.OutputScript
			inputScript := transactionInput.InputScript
			script := ScriptTool.CreateScript(inputScript, outputScript)
			result := b.virtualMachine.Execute(transaction, script)
			//TODO java...
			executeSuccess := result.Size() == 1 && StringUtil.Equals(ByteUtil.BytesToHexString(BooleanCode.TRUE.Code), *result.Pop())
			if !executeSuccess {
				return false
			}
		}
	}
	return true
}

//region dto to model
func (b *BlockchainDatabase) BlockDto2Block(blockDto *dto.BlockDto) *model.Block {
	previousHash := blockDto.PreviousHash
	previousBlock := b.QueryBlockByBlockHash(previousHash)
	block := &model.Block{}
	block.Timestamp = blockDto.Timestamp
	block.PreviousHash = previousHash
	block.Nonce = blockDto.Nonce

	blockHeight := BlockTool.GetNextBlockHeight(previousBlock)
	block.Height = blockHeight
	transactions := b.transactionDtos2Transactions(blockDto.Transactions)
	block.Transactions = transactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(block)
	block.MerkleTreeRoot = merkleTreeRoot

	blockHash := BlockTool.CalculateBlockHash(block)
	block.Hash = blockHash

	difficult := b.consensus.CalculateDifficult(b, block)
	block.Difficulty = difficult

	b.fillBlockProperty(block)

	if !b.consensus.CheckConsensus(b, block) {
		//TODO throw new RuntimeException("Check Block Consensus Failed.")
		return nil
	}
	return block
}
func (b *BlockchainDatabase) transactionDtos2Transactions(transactionDtos []*dto.TransactionDto) []*model.Transaction {
	var transactions []*model.Transaction
	if transactionDtos != nil {
		for _, transactionDto := range transactionDtos {
			transaction := b.TransactionDto2Transaction(transactionDto)
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}
func (b *BlockchainDatabase) TransactionDto2Transaction(transactionDto *dto.TransactionDto) *model.Transaction {
	var inputs []*model.TransactionInput
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			unspentTransactionOutput := b.QueryUnspentTransactionOutputByTransactionOutputId(transactionInputDto.TransactionHash, transactionInputDto.TransactionOutputIndex)
			if unspentTransactionOutput == nil {
				//TODO throw new RuntimeException("Illegal transaction. the transaction input is not an unspent transaction output.");
				return nil
			}
			var transactionInput model.TransactionInput
			transactionInput.UnspentTransactionOutput = unspentTransactionOutput
			transactionInput.InputScript = b.InputScriptDto2InputScript(transactionInputDto.InputScript)
			inputs = append(inputs, &transactionInput)
		}
	}
	var outputs []*model.TransactionOutput
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			transactionOutput := b.transactionOutputDto2TransactionOutput(transactionOutputDto)
			outputs = append(outputs, transactionOutput)
		}
	}
	transaction := new(model.Transaction)
	transactionType := b.obtainTransactionDto(transactionDto)
	transaction.TransactionType = transactionType
	transaction.TransactionHash = TransactionDtoTool.CalculateTransactionHash(transactionDto)
	transaction.Inputs = inputs
	transaction.Outputs = outputs
	return transaction
}
func (b *BlockchainDatabase) transactionOutputDto2TransactionOutput(transactionOutputDto *dto.TransactionOutputDto) *model.TransactionOutput {
	var transactionOutput model.TransactionOutput
	publicKeyHash := ScriptDtoTool.GetPublicKeyHashFromPayToPublicKeyHashOutputScript(transactionOutputDto.OutputScript)
	address := AccountUtil.AddressFromPublicKeyHash(publicKeyHash)
	transactionOutput.Address = address
	transactionOutput.Value = transactionOutputDto.Value
	transactionOutput.OutputScript = b.OutputScriptDto2OutputScript(transactionOutputDto.OutputScript)
	return &transactionOutput
}
func (b *BlockchainDatabase) obtainTransactionDto(transactionDto *dto.TransactionDto) TransactionType.TransactionType {
	if transactionDto.Inputs == nil || len(transactionDto.Inputs) == 0 {
		return TransactionType.COINBASE_TRANSACTION
	}
	return TransactionType.STANDARD_TRANSACTION
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
	if transactions != nil {
		for _, transaction := range transactions {
			transactionIndex = transactionIndex + 1
			transactionHeight = transactionHeight + 1
			transaction.BlockHeight = blockHeight
			transaction.TransactionIndex = transactionIndex
			transaction.TransactionHeight = transactionHeight

			outputs := transaction.Outputs
			if outputs != nil {
				for i := 0; i < len(outputs); i = i + 1 {
					transactionOutputHeight = transactionOutputHeight + 1
					output := outputs[i]
					output.BlockHeight = blockHeight
					output.BlockHash = blockHash
					output.TransactionHeight = transactionHeight
					output.TransactionHash = transaction.TransactionHash
					output.TransactionOutputIndex = uint64(i) + uint64(1)
					output.TransactionIndex = transaction.TransactionIndex
					output.TransactionOutputHeight = transactionOutputHeight
				}
			}
		}
	}
}
func (b *BlockchainDatabase) OutputScriptDto2OutputScript(outputScriptDto *dto.OutputScriptDto) *model.OutputScript {
	var outputScript model.OutputScript
	outputScript = append(outputScript, *outputScriptDto...)
	return &outputScript
}
func (b *BlockchainDatabase) InputScriptDto2InputScript(inputScriptDto *dto.InputScriptDto) *model.InputScript {
	var inputScript model.InputScript
	inputScript = append(inputScript, *inputScriptDto...)
	return &inputScript
}

//endregion
