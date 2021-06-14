package core

import (
	"fmt"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/RandomUtil"
	"helloworldcoin-go/util/JsonUtil"
	"helloworldcoin-go/util/SleepUtil"
	"helloworldcoin-go/util/TimeUtil"
)

type Miner struct {
	CoreConfiguration              *CoreConfiguration
	Wallet                         *Wallet
	BlockchainDatabase             *BlockchainDatabase
	UnconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
}

func (i *Miner) Start() {
	for {
		SleepUtil.Sleep(10)
		if !i.isActive() {
			continue
		}
		minerAccount := i.Wallet.CreateAccount()
		block := BuildMiningBlock(i.BlockchainDatabase, i.UnconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.CurrentMillisecondTimestamp()
		for {
			if !i.isActive() {
				break
			}
			//在挖矿的期间，可能收集到新的交易。每隔一定的时间，重新组装挖矿中的区块，这样新收集到交易就可以被放进挖矿中的区块了。
			if TimeUtil.CurrentMillisecondTimestamp()-startTimestamp > i.CoreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			//随机数
			block.Nonce = HexUtil.BytesToHexString(RandomUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			//挖矿成功
			if i.BlockchainDatabase.Consensus.CheckConsensus(i.BlockchainDatabase, block) {
				i.Wallet.SaveAccount(minerAccount)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				fmt.Println(JsonUtil.ToJson(blockDto))
				isAddBlockToBlockchainSuccess := i.BlockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					//LogUtil.debug("挖矿成功，但是区块放入区块链失败。")
				}
				break
			}
		}
	}
}
func (i *Miner) isActive() bool {
	//TODO
	return true
}
