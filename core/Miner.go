package core

import (
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/RandomUtil"
	"helloworldcoin-go/util/SleepUtil"
	"helloworldcoin-go/util/TimeUtil"
)

type Miner struct {
	coreConfiguration              CoreConfiguration
	wallet                         Wallet
	blockchainDataBase             BlockchainDatabase
	unconfirmedTransactionDataBase UnconfirmedTransactionDatabase
}

func (i *Miner) start() {
	for {
		SleepUtil.Sleep(10)
		if !i.isActive() {
			continue
		}
		minerAccount := i.wallet.CreateAccount()
		block := BuildMiningBlock(&i.blockchainDataBase, &i.unconfirmedTransactionDataBase, &minerAccount)
		startTimestamp := TimeUtil.CurrentMillisecondTimestamp()
		for {
			if !i.isActive() {
				break
			}
			//在挖矿的期间，可能收集到新的交易。每隔一定的时间，重新组装挖矿中的区块，这样新收集到交易就可以被放进挖矿中的区块了。
			if TimeUtil.CurrentMillisecondTimestamp()-startTimestamp > i.coreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			//随机数
			block.Nonce = HexUtil.BytesToHexString(RandomUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			//挖矿成功
			if i.blockchainDataBase.Consensus.CheckConsensus(&i.blockchainDataBase, block) {
				i.wallet.SaveAccount(&minerAccount)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				isAddBlockToBlockchainSuccess := i.blockchainDataBase.AddBlockDto(blockDto)
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
