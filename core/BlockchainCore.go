package core

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/model/ModelWallet"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
)

type BlockchainCore struct {
	BlockchainDatabase             *BlockchainDatabase
	UnconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
	CoreConfiguration              *CoreConfiguration
	Wallet                         *Wallet
	Miner                          *Miner
}

func (b *BlockchainCore) Start() {
	go b.Miner.Start()
}
func (b *BlockchainCore) QueryBlockByBlockHeight(blockHeight uint64) *model.Block {
	return b.BlockchainDatabase.QueryBlockByBlockHeight(blockHeight)
}
func (b *BlockchainCore) AddBlockDto(blockDto *dto.BlockDto) bool {
	return b.BlockchainDatabase.AddBlockDto(blockDto)
}

func (b *BlockchainCore) QueryBlockchainHeight() uint64 {
	return b.BlockchainDatabase.QueryBlockchainHeight()
}

func (b *BlockchainCore) PostTransaction(transactionDto *dto.TransactionDto) {
	b.UnconfirmedTransactionDatabase.InsertTransaction(transactionDto)
}

func (b *BlockchainCore) QueryTailBlock() *model.Block {
	return b.BlockchainDatabase.QueryTailBlock()
}

func (b *BlockchainCore) DeleteTailBlock() {
	b.BlockchainDatabase.DeleteTailBlock()
}

func (b *BlockchainCore) AddBlock(block *model.Block) bool {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return b.AddBlockDto(blockDto)
}

func (b *BlockchainCore) DeleteBlocks(blockHeight uint64) {
	b.BlockchainDatabase.DeleteBlocks(blockHeight)
}

func (b *BlockchainCore) GetMiner() *Miner {
	return b.Miner
}

func (b *BlockchainCore) GetWallet() *Wallet {
	return b.Wallet
}
func (b *BlockchainCore) AutoBuildTransaction(request *ModelWallet.AutoBuildTransactionRequest) *ModelWallet.AutoBuildTransactionResponse {
	return b.Wallet.AutoBuildTransaction(request)
}

func (b *BlockchainCore) QueryUnconfirmedTransactions(from uint64, size uint64) []*dto.TransactionDto {
	return b.UnconfirmedTransactionDatabase.SelectTransactions(from, size)
}

func (b *BlockchainCore) QueryBlockByBlockHash(blockHash string) *model.Block {
	return b.BlockchainDatabase.QueryBlockByBlockHash(blockHash)
}

func (b *BlockchainCore) GetBlockchainDatabase() *BlockchainDatabase {
	return b.BlockchainDatabase
}

func (b *BlockchainCore) QueryTransactionByTransactionHash(transactionHash string) *model.Transaction {
	return b.BlockchainDatabase.QueryTransactionByTransactionHash(transactionHash)

}

func (b *BlockchainCore) QueryUnconfirmedTransactionByTransactionHash(transactionHash string) *dto.TransactionDto {
	return b.UnconfirmedTransactionDatabase.SelectTransactionByTransactionHash(transactionHash)
}

func (b *BlockchainCore) QueryTransactionByTransactionHeight(transactionHeight uint64) *model.Transaction {
	return b.BlockchainDatabase.QueryTransactionByTransactionHeight(transactionHeight)
}

func (b *BlockchainCore) QueryTransactionOutputByAddress(address string) *model.TransactionOutput {
	return b.BlockchainDatabase.QueryTransactionOutputByAddress(address)
}
