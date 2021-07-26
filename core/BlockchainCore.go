package core

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/dto"
)

type BlockchainCore struct {
	BlockchainDatabase *BlockchainDatabase
}

func (b BlockchainCore) QueryBlockByBlockHeight(blockHeight uint64) *Model.Block {
	return b.BlockchainDatabase.QueryBlockByBlockHeight(blockHeight)
}
func (b BlockchainCore) AddBlockDto(blockDto *dto.BlockDto) bool {
	return b.BlockchainDatabase.AddBlockDto(blockDto)
}
