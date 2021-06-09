package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/Model/BlockchainActionEnum"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
	"sync"
)

const BLOCKCHAIN_DATABASE_NAME = "BlockchainDatabase"

type BlockchainDatabase struct {
	consensus         Consensus
	incentive         Incentive
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

func (blockchainDatabase *BlockchainDatabase) QueryBlockchainHeight() int {
	return 1

}
func (blockchainDatabase *BlockchainDatabase) QueryBlockchainTransactionHeight() int {
	return 1

}
func (blockchainDatabase *BlockchainDatabase) QueryBlockchainTransactionOutputHeight() int {
	return 1

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
func (blockchainDatabase *BlockchainDatabase) QuerySourceTransactionByTransactionOutputId(transactionOutputId *Model.TransactionOutputId) *Model.Transaction {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QueryDestinationTransactionByTransactionOutputId(transactionOutputId *Model.TransactionOutputId) *Model.Transaction {
	return nil

}

func (blockchainDatabase *BlockchainDatabase) QueryTransactionOutputByTransactionOutputHeight(transactionOutputHeight uint64) *Model.TransactionOutput {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QueryTransactionOutputByTransactionOutputId(transactionOutputId *Model.TransactionOutputId) *Model.TransactionOutput {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QueryUnspentTransactionOutputByTransactionOutputId(transactionOutputId *Model.TransactionOutputId) *Model.TransactionOutput {
	return nil

}
func (blockchainDatabase *BlockchainDatabase) QuerySpentTransactionOutputByTransactionOutputId(transactionOutputId *Model.TransactionOutputId) *Model.TransactionOutput {
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

	return nil

}
func (blockchainDatabase *BlockchainDatabase) fillBlockProperty(block *Model.Block) {
	transactionIndex := 0
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
			output.TransactionOutputIndex = index + 1
			output.TransactionIndex = transaction.TransactionIndex
			output.TransactionOutputHeight = transactionOutputHeight
		}
	}
}
