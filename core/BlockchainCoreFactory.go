package core

/**
 * 创建BlockchainCore实例
 *
 * @param corePath BlockchainCore数据存放位置
 */
func CreateBlockchainCore(corePath string) *BlockchainCore {

	coreConfiguration := &CoreConfiguration{CorePath: corePath}
	incentive := &Incentive{}
	consensus := &Consensus{}
	virtualMachine := &VirtualMachine{}
	blockchainDatabase := &BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration, VirtualMachine: virtualMachine}

	unconfirmedTransactionDatabase := &UnconfirmedTransactionDatabase{CoreConfiguration: coreConfiguration}
	wallet := &Wallet{CoreConfiguration: coreConfiguration, BlockchainDatabase: blockchainDatabase}
	miner := &Miner{CoreConfiguration: coreConfiguration, Wallet: wallet, BlockchainDatabase: blockchainDatabase, UnconfirmedTransactionDatabase: unconfirmedTransactionDatabase}
	blockchainCore := &BlockchainCore{CoreConfiguration: coreConfiguration, BlockchainDatabase: blockchainDatabase, UnconfirmedTransactionDatabase: unconfirmedTransactionDatabase, Wallet: wallet, Miner: miner}
	return blockchainCore
}
