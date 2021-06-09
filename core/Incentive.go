package core

import (
	"helloworldcoin-go/core/Model"
)

type Incentive struct {
}

func (incentive *Incentive) incentiveValue(blockchainDataBase *BlockchainDatabase, block *Model.Block) uint64 {
	//TODO
	return 50
}

func (incentive *Incentive) checkIncentive(blockchainDataBase *BlockchainDatabase, block *Model.Block) bool {
	//TODO
	return true
}
