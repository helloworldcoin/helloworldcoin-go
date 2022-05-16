package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/SizeTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/setting/BlockSetting"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/StringUtil"
	"helloworldcoin-go/util/ThreadUtil"
	"helloworldcoin-go/util/TimeUtil"
)

type Miner struct {
	coreConfiguration              *CoreConfiguration
	wallet                         *Wallet
	blockchainDatabase             *BlockchainDatabase
	unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
}

func NewMiner(coreConfiguration *CoreConfiguration, wallet *Wallet, blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) *Miner {
	var miner Miner
	miner.coreConfiguration = coreConfiguration
	miner.wallet = wallet
	miner.blockchainDatabase = blockchainDatabase
	miner.unconfirmedTransactionDatabase = unconfirmedTransactionDatabase
	return &miner
}

func (i *Miner) Start() {
	for {
		ThreadUtil.MillisecondSleep(10)
		if !i.IsActive() {
			continue
		}
		if i.isMiningHeightExceedsLimit() {
			continue
		}

		blockChainHeight := i.blockchainDatabase.QueryBlockchainHeight()

		if blockChainHeight >= i.coreConfiguration.getMinerMineMaxBlockHeight() {
			continue
		}

		minerAccount := i.wallet.CreateAccount()
		block := i.buildMiningBlock(i.blockchainDatabase, i.unconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.MillisecondTimestamp()
		for {
			if !i.IsActive() {
				break
			}
			if TimeUtil.MillisecondTimestamp()-startTimestamp > i.coreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			block.Nonce = ByteUtil.BytesToHexString(ByteUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			if i.blockchainDatabase.consensus.CheckConsensus(i.blockchainDatabase, block) {
				i.wallet.SaveAccount(minerAccount)
				LogUtil.Debug("Congratulations! Mining success! Block height:" + StringUtil.ValueOfUint64(block.Height) + ", Block hash:" + block.Hash)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				isAddBlockToBlockchainSuccess := i.blockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					LogUtil.Debug("Mining succeeded, but the block failed to be put into the blockchain.")
				}
				break
			}
		}
	}
}
func (i *Miner) IsActive() bool {
	return i.coreConfiguration.IsMinerActive()
}
func (i *Miner) Deactive() {
	i.coreConfiguration.deactiveMiner()
}
func (i *Miner) Active() {
	i.coreConfiguration.activeMiner()
}

func (i *Miner) SetMinerMineMaxBlockHeight(maxHeight uint64) {
	i.coreConfiguration.setMinerMineMaxBlockHeight(maxHeight)
}

func (i *Miner) GetMinerMineMaxBlockHeight() uint64 {
	return i.coreConfiguration.getMinerMineMaxBlockHeight()
}

func (i *Miner) buildMiningBlock(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase, minerAccount *AccountUtil.Account) *model.Block {
	timestamp := TimeUtil.MillisecondTimestamp()

	tailBlock := blockchainDatabase.QueryTailBlock()
	var nonNonceBlock model.Block
	nonNonceBlock.Timestamp = timestamp

	if tailBlock == nil {
		nonNonceBlock.Height = GenesisBlockSetting.HEIGHT + uint64(1)
		nonNonceBlock.PreviousHash = GenesisBlockSetting.HASH
	} else {
		nonNonceBlock.Height = tailBlock.Height + uint64(1)
		nonNonceBlock.PreviousHash = tailBlock.Hash
	}
	packingTransactions := i.packingTransactions(blockchainDatabase, unconfirmedTransactionDatabase)

	incentive := blockchainDatabase.incentive
	incentiveValue := incentive.IncentiveValue(blockchainDatabase, &nonNonceBlock)

	mineAwardTransaction := i.buildIncentiveTransaction(minerAccount.Address, incentiveValue)
	var mineAwardTransactions []*model.Transaction
	mineAwardTransactions = append(mineAwardTransactions, mineAwardTransaction)

	packingTransactions = append(mineAwardTransactions, packingTransactions...)
	nonNonceBlock.Transactions = packingTransactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(&nonNonceBlock)
	nonNonceBlock.MerkleTreeRoot = merkleTreeRoot

	nonNonceBlock.Difficulty = blockchainDatabase.consensus.CalculateDifficult(blockchainDatabase, &nonNonceBlock)
	return &nonNonceBlock
}

func (i *Miner) buildIncentiveTransaction(address string, incentiveValue uint64) *model.Transaction {
	var transaction model.Transaction
	transaction.TransactionType = TransactionType.COINBASE_TRANSACTION

	var outputs []*model.TransactionOutput
	var output model.TransactionOutput
	output.Address = address
	output.Value = incentiveValue
	output.OutputScript = ScriptTool.CreatePayToPublicKeyHashOutputScript(address)
	outputs = append(outputs, &output)

	transaction.Outputs = outputs
	transaction.TransactionHash = TransactionTool.CalculateTransactionHash(&transaction)
	return &transaction
}
func (i *Miner) packingTransactions(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) []*model.Transaction {
	forMineBlockTransactionDtos := unconfirmedTransactionDatabase.SelectTransactions(uint64(1), uint64(10000))

	var transactions []*model.Transaction
	var backupTransactions []*model.Transaction

	if forMineBlockTransactionDtos != nil {
		for _, transactionDto := range forMineBlockTransactionDtos {
			defer func() {
				if e := recover(); e != nil {
					transactionHash := TransactionDtoTool.CalculateTransactionHash(transactionDto)
					LogUtil.Error("Abnormal transaction, transaction hash:"+transactionHash, e)
					unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
				}
			}()
			transaction := blockchainDatabase.TransactionDto2Transaction(transactionDto)
			transactions = append(transactions, transaction)
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}
	for _, transaction := range backupTransactions {
		checkTransaction := blockchainDatabase.CheckTransaction(transaction)
		if checkTransaction {
			transactions = append(transactions, transaction)
		} else {
			transactionHash := TransactionTool.CalculateTransactionHash(transaction)
			LogUtil.Debug("Abnormal transaction, transaction hash:" + transactionHash)
			unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}

	//prevent double spending
	transactionOutputIdSet := make(map[string]bool)
	for _, transaction := range backupTransactions {
		inputs := transaction.Inputs
		if inputs != nil {
			canAdd := true
			for _, transactionInput := range inputs {
				unspentTransactionOutput := transactionInput.UnspentTransactionOutput
				transactionOutputId := TransactionTool.GetTransactionOutputId(unspentTransactionOutput)
				_, exists := transactionOutputIdSet[transactionOutputId]
				if exists {
					canAdd = false
					break
				} else {
					transactionOutputIdSet[transactionOutputId] = true
				}
			}
			if canAdd {
				transactions = append(transactions, transaction)
			} else {
				transactionHash := TransactionTool.CalculateTransactionHash(transaction)
				LogUtil.Debug("Abnormal transaction, transaction hash:" + transactionHash)
				unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
			}
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}

	//Prevent an address used multiple times
	addressSet := make(map[string]bool)
	for _, transaction := range backupTransactions {
		outputs := transaction.Outputs
		if outputs != nil {
			canAdd := true
			for _, output := range outputs {
				address := output.Address
				_, exists := addressSet[address]
				if exists {
					canAdd = false
					break
				} else {
					addressSet[address] = true
				}
			}
			if canAdd {
				transactions = append(transactions, transaction)
			} else {
				transactionHash := TransactionTool.CalculateTransactionHash(transaction)
				LogUtil.Debug("Abnormal transaction, transaction hash:" + transactionHash)
				unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
			}
		}
	}

	//TODO TransactionTool.sortByTransactionFeeRateDescend(transactions);

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}

	size := uint64(0)
	for i := 0; i < len(backupTransactions); i++ {
		if uint64(i+1) > BlockSetting.BLOCK_MAX_TRANSACTION_COUNT-1 {
			break
		}
		transaction := backupTransactions[i]
		size += SizeTool.CalculateTransactionSize(transaction)
		if size > BlockSetting.BLOCK_MAX_CHARACTER_COUNT {
			break
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
func (i *Miner) isMiningHeightExceedsLimit() bool {
	blockChainHeight := i.blockchainDatabase.QueryBlockchainHeight()
	return blockChainHeight >= i.coreConfiguration.getMinerMineMaxBlockHeight()
}
