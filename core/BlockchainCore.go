package core

import (
	"helloworld-blockchain-go/core/Model"
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
	//TODO 异常
	go b.Miner.Start()
}
func (b BlockchainCore) QueryBlockByBlockHeight(blockHeight uint64) *Model.Block {
	return b.BlockchainDatabase.QueryBlockByBlockHeight(blockHeight)
}
func (b BlockchainCore) AddBlockDto(blockDto *dto.BlockDto) bool {
	return b.BlockchainDatabase.AddBlockDto(blockDto)
}

func (b BlockchainCore) QueryBlockchainHeight() uint64 {
	return b.BlockchainDatabase.QueryBlockchainHeight()
}

func (b BlockchainCore) PostTransaction(transactionDto *dto.TransactionDto) {
	b.UnconfirmedTransactionDatabase.InsertTransaction(transactionDto)
}

func (b BlockchainCore) QueryTailBlock() *Model.Block {
	return b.BlockchainDatabase.QueryTailBlock()
}

func (b BlockchainCore) DeleteTailBlock() {
	b.BlockchainDatabase.DeleteTailBlock()
}

func (b BlockchainCore) AddBlock(block *Model.Block) bool {
	blockDto := Model2DtoTool.Block2BlockDto(block)
	return b.AddBlockDto(blockDto)
}

func (b BlockchainCore) DeleteBlocks(blockHeight uint64) {
	b.BlockchainDatabase.DeleteBlocks(blockHeight)
}

func (b *BlockchainCore) GetMiner() *Miner {
	return b.Miner
}
