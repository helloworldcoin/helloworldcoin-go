package core

/*
 @author king 409060350@qq.com
*/

type BlockchainCoreFactory struct {
}

/**
 * 创建BlockchainCore实例
 *
 * @param corePath BlockchainCore数据存放位置
 */
func (b *BlockchainCoreFactory) CreateBlockchainCore(corePath string) *BlockchainCore {

	coreConfiguration := &CoreConfiguration{corePath: corePath}
	incentive := &Incentive{}
	consensus := &Consensus{}
	virtualMachine := &VirtualMachine{}
	blockchainDatabase := NewBlockchainDatabase(consensus, incentive, virtualMachine, coreConfiguration)

	unconfirmedTransactionDatabase := NewUnconfirmedTransactionDatabase(coreConfiguration)
	wallet := NewWallet(coreConfiguration, blockchainDatabase)
	miner := NewMiner(coreConfiguration, wallet, blockchainDatabase, unconfirmedTransactionDatabase)
	blockchainCore := NewBlockchainCore(coreConfiguration, blockchainDatabase, unconfirmedTransactionDatabase, wallet, miner)
	return blockchainCore
}
