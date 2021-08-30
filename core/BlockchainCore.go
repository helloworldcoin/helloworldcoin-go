package core

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/dto"
)

type BlockchainCore struct {
	BlockchainDatabase             *BlockchainDatabase
	UnconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
	CoreConfiguration              *CoreConfiguration
	Wallet                         *Wallet
	Miner                          *Miner
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
