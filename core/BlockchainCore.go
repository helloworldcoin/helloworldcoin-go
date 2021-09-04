package core

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
)

type BlockchainCore struct {
	blockchainDatabase             *BlockchainDatabase
	unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
	coreConfiguration              *CoreConfiguration
	wallet                         *Wallet
	miner                          *Miner
}

func NewBlockchainCore(coreConfiguration *CoreConfiguration, blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase, wallet *Wallet, miner *Miner) *BlockchainCore {
	var blockchainCore BlockchainCore
	blockchainCore.coreConfiguration = coreConfiguration
	blockchainCore.blockchainDatabase = blockchainDatabase
	blockchainCore.unconfirmedTransactionDatabase = unconfirmedTransactionDatabase
	blockchainCore.wallet = wallet
	blockchainCore.miner = miner
	return &blockchainCore
}
func (b *BlockchainCore) GetUnconfirmedTransactionDatabase() *UnconfirmedTransactionDatabase {
	return b.unconfirmedTransactionDatabase
}

func (b *BlockchainCore) Start() {
	go b.miner.Start()
}
func (b *BlockchainCore) QueryBlockByBlockHeight(blockHeight uint64) *model.Block {
	return b.blockchainDatabase.QueryBlockByBlockHeight(blockHeight)
}
func (b *BlockchainCore) AddBlockDto(blockDto *dto.BlockDto) bool {
	return b.blockchainDatabase.AddBlockDto(blockDto)
}

func (b *BlockchainCore) QueryBlockchainHeight() uint64 {
	return b.blockchainDatabase.QueryBlockchainHeight()
}

func (b *BlockchainCore) PostTransaction(transactionDto *dto.TransactionDto) {
	b.unconfirmedTransactionDatabase.InsertTransaction(transactionDto)
}

func (b *BlockchainCore) QueryTailBlock() *model.Block {
	return b.blockchainDatabase.QueryTailBlock()
}

func (b *BlockchainCore) DeleteTailBlock() {
	b.blockchainDatabase.DeleteTailBlock()
}

func (b *BlockchainCore) AddBlock(block *model.Block) bool {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return b.AddBlockDto(blockDto)
}

func (b *BlockchainCore) DeleteBlocks(blockHeight uint64) {
	b.blockchainDatabase.DeleteBlocks(blockHeight)
}

func (b *BlockchainCore) GetMiner() *Miner {
	return b.miner
}

func (b *BlockchainCore) GetWallet() *Wallet {
	return b.wallet
}
func (b *BlockchainCore) AutoBuildTransaction(request *model.AutoBuildTransactionRequest) *model.AutoBuildTransactionResponse {
	return b.wallet.AutoBuildTransaction(request)
}

func (b *BlockchainCore) QueryUnconfirmedTransactions(from uint64, size uint64) []*dto.TransactionDto {
	return b.unconfirmedTransactionDatabase.SelectTransactions(from, size)
}

func (b *BlockchainCore) QueryBlockByBlockHash(blockHash string) *model.Block {
	return b.blockchainDatabase.QueryBlockByBlockHash(blockHash)
}

func (b *BlockchainCore) GetBlockchainDatabase() *BlockchainDatabase {
	return b.blockchainDatabase
}

func (b *BlockchainCore) QueryTransactionByTransactionHash(transactionHash string) *model.Transaction {
	return b.blockchainDatabase.QueryTransactionByTransactionHash(transactionHash)

}

func (b *BlockchainCore) QueryUnconfirmedTransactionByTransactionHash(transactionHash string) *dto.TransactionDto {
	return b.unconfirmedTransactionDatabase.SelectTransactionByTransactionHash(transactionHash)
}

func (b *BlockchainCore) QueryTransactionByTransactionHeight(transactionHeight uint64) *model.Transaction {
	return b.blockchainDatabase.QueryTransactionByTransactionHeight(transactionHeight)
}

func (b *BlockchainCore) QueryTransactionOutputByAddress(address string) *model.TransactionOutput {
	return b.blockchainDatabase.QueryTransactionOutputByAddress(address)
}
func (b *BlockchainCore) BlockDto2Block(blockDto *dto.BlockDto) *model.Block {
	return b.blockchainDatabase.BlockDto2Block(blockDto)
}
func (b *BlockchainCore) TransactionDto2Transaction(transactionDto *dto.TransactionDto) *model.Transaction {
	return b.blockchainDatabase.TransactionDto2Transaction(transactionDto)
}
