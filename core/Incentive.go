package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/setting/IncentiveSetting"
)

type Incentive struct {
}

func (incentive *Incentive) IncentiveValue(blockchainDatabase *BlockchainDatabase, block *model.Block) uint64 {
	minerSubsidy := getMinerSubsidy(block)
	minerFee := BlockTool.GetBlockFee(block)
	return minerSubsidy + minerFee
}

func (incentive *Incentive) checkIncentive(blockchainDatabase *BlockchainDatabase, block *model.Block) bool {
	writeIncentiveValue := BlockTool.GetWritedIncentiveValue(block)
	targetIncentiveValue := incentive.IncentiveValue(blockchainDatabase, block)
	if writeIncentiveValue != targetIncentiveValue {
		return false
	}
	return true
}

func getMinerSubsidy(block *model.Block) uint64 {
	initCoin := IncentiveSetting.BLOCK_INIT_INCENTIVE
	multiple := (block.Height - uint64(1)) / IncentiveSetting.INCENTIVE_HALVING_INTERVAL

	for multiple > 0 {
		initCoin = initCoin / uint64(2)
		multiple = multiple - uint64(1)
	}
	return initCoin
}
