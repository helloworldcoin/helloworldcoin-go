package core

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/core/Model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
	"helloworld-blockchain-go/util/TimeUtil"
)

type Miner struct {
	CoreConfiguration              *CoreConfiguration
	Wallet                         *Wallet
	BlockchainDatabase             *BlockchainDatabase
	UnconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
}

func (i *Miner) Start() {
	for {
		ThreadUtil.MillisecondSleep(10)
		if !i.isActive() {
			continue
		}
		minerAccount := i.Wallet.CreateAccount()
		block := i.buildMiningBlock(i.BlockchainDatabase, i.UnconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.MillisecondTimestamp()
		for {
			if !i.isActive() {
				break
			}
			//在挖矿的期间，可能收集到新的交易。每隔一定的时间，重新组装挖矿中的区块，这样新收集到交易就可以被放进挖矿中的区块了。
			if TimeUtil.MillisecondTimestamp()-startTimestamp > i.CoreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			//随机数
			block.Nonce = ByteUtil.BytesToHexString(ByteUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			//挖矿成功
			if i.BlockchainDatabase.Consensus.CheckConsensus(i.BlockchainDatabase, block) {
				i.Wallet.SaveAccount(minerAccount)
				LogUtil.Debug("祝贺您！挖矿成功！！！区块高度:" + StringUtil.ValueOfUint64(block.Height) + ",区块哈希:" + block.Hash)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				isAddBlockToBlockchainSuccess := i.BlockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					LogUtil.Debug("挖矿成功，但是区块放入区块链失败。")
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

func (i *Miner) buildMiningBlock(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase, minerAccount *AccountUtil.Account) *Model.Block {
	timestamp := TimeUtil.MillisecondTimestamp()

	tailBlock := blockchainDatabase.QueryTailBlock()
	var nonNonceBlock Model.Block
	nonNonceBlock.Timestamp = timestamp

	if tailBlock == nil {
		nonNonceBlock.Height = GenesisBlockSetting.HEIGHT + uint64(1)
		nonNonceBlock.PreviousHash = GenesisBlockSetting.HASH
	} else {
		nonNonceBlock.Height = tailBlock.Height + uint64(1)
		nonNonceBlock.PreviousHash = tailBlock.Hash
	}
	packingTransactions := i.packingTransactions(blockchainDatabase, unconfirmedTransactionDatabase)

	incentive := blockchainDatabase.Incentive
	incentiveValue := incentive.IncentiveValue(blockchainDatabase, &nonNonceBlock)

	mineAwardTransaction := i.buildIncentiveTransaction(minerAccount.Address, incentiveValue)
	var mineAwardTransactions []*Model.Transaction
	mineAwardTransactions = append(mineAwardTransactions, mineAwardTransaction)

	packingTransactions = append(mineAwardTransactions, packingTransactions...)
	nonNonceBlock.Transactions = packingTransactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(&nonNonceBlock)
	nonNonceBlock.MerkleTreeRoot = merkleTreeRoot

	//计算挖矿难度
	nonNonceBlock.Difficulty = blockchainDatabase.Consensus.CalculateDifficult(blockchainDatabase, &nonNonceBlock)
	return &nonNonceBlock
}

func (i *Miner) buildIncentiveTransaction(address string, incentiveValue uint64) *Model.Transaction {
	var transaction Model.Transaction
	transaction.TransactionType = TransactionType.GENESIS_TRANSACTION

	var outputs []*Model.TransactionOutput
	var output Model.TransactionOutput
	output.Address = address
	output.Value = incentiveValue
	output.OutputScript = ScriptTool.CreatePayToPublicKeyHashOutputScript(address)
	outputs = append(outputs, &output)

	transaction.Outputs = outputs
	transaction.TransactionHash = TransactionTool.CalculateTransactionHash(transaction)
	return &transaction
}
func (i *Miner) packingTransactions(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) []*Model.Transaction {
	forMineBlockTransactionDtos := unconfirmedTransactionDatabase.SelectTransactions(uint64(1), uint64(10000))

	var transactions []*Model.Transaction
	if forMineBlockTransactionDtos != nil {
		for _, transactionDto := range forMineBlockTransactionDtos {
			//TODO exception
			transaction := TransactionDto2Transaction(blockchainDatabase, &transactionDto)
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}
