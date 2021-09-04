package BlockchainCoreFactory

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
)

/*
 @author king 409060350@qq.com
*/

/**
 * 创建BlockchainCore实例
 */
func CreateDefaultBlockchainCore() *core.BlockchainCore {
	return CreateBlockchainCore(ResourcePathTool.GetDataRootPath())
}

/**
 * 创建BlockchainCore实例
 *
 * @param corePath BlockchainCore数据存放位置
 */
func CreateBlockchainCore(corePath string) *core.BlockchainCore {

	coreConfiguration := core.NewCoreConfiguration(corePath)
	incentive := core.NewIncentive()
	consensus := core.NewConsensus()
	virtualMachine := core.NewVirtualMachine()
	blockchainDatabase := core.NewBlockchainDatabase(consensus, incentive, virtualMachine, coreConfiguration)

	unconfirmedTransactionDatabase := core.NewUnconfirmedTransactionDatabase(coreConfiguration)
	wallet := core.NewWallet(coreConfiguration, blockchainDatabase)
	miner := core.NewMiner(coreConfiguration, wallet, blockchainDatabase, unconfirmedTransactionDatabase)
	blockchainCore := core.NewBlockchainCore(coreConfiguration, blockchainDatabase, unconfirmedTransactionDatabase, wallet, miner)
	return blockchainCore
}
