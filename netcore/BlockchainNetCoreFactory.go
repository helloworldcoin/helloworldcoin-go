package netcore

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/dao"
	"helloworld-blockchain-go/netcore/server"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/FileUtil"
)

/**
 * 创建[区块链网络版核心]实例
 */
func CreateDefaultBlockchainNetCore() *BlockchainNetCore {
	return CreateBlockchainNetCore(ResourcePathTool.GetDataRootPath())
}

/**
 * 创建[区块链网络版核心]实例
 *
 * @param netcorePath 区块链数据存放位置
 */
func CreateBlockchainNetCore(netcorePath string) *BlockchainNetCore {
	netCoreConfiguration := &configuration.NetCoreConfiguration{NetCorePath: netcorePath}

	blockchainCorePath := FileUtil.NewPath(netcorePath, "BlockchainCore")
	blockchainCore := (&core.BlockchainCoreFactory{}).CreateBlockchainCore(blockchainCorePath)
	slaveBlockchainCorePath := FileUtil.NewPath(netcorePath, "SlaveBlockchainCore")
	slaveBlockchainCore := (&core.BlockchainCoreFactory{}).CreateBlockchainCore(slaveBlockchainCorePath)

	nodeDao := dao.NewNodeDao(netCoreConfiguration)
	nodeService := service.NewNodeService(nodeDao)
	nodeServer := server.NewNodeServer(netCoreConfiguration, blockchainCore, nodeService)

	seedNodeInitializer := &SeedNodeInitializer{netCoreConfiguration, nodeService}
	nodeSearcher := &NodeSearcher{netCoreConfiguration, nodeService}
	nodeBroadcaster := &NodeBroadcaster{netCoreConfiguration, nodeService}
	nodeCleaner := &NodeCleaner{netCoreConfiguration, nodeService}

	blockchainHeightSearcher := &BlockchainHeightSearcher{netCoreConfiguration, nodeService}
	blockchainHeightBroadcaster := &BlockchainHeightBroadcaster{netCoreConfiguration, blockchainCore, nodeService}

	blockSearcher := &BlockSearcher{netCoreConfiguration, blockchainCore, slaveBlockchainCore, nodeService}
	blockBroadcaster := &BlockBroadcaster{netCoreConfiguration, blockchainCore, nodeService}

	unconfirmedTransactionsSearcher := &UnconfirmedTransactionsSearcher{netCoreConfiguration, blockchainCore, nodeService}

	blockchainNetCore := BlockchainNetCore{netCoreConfiguration, nodeService, blockchainCore, nodeServer, seedNodeInitializer, nodeSearcher, nodeBroadcaster, nodeCleaner, blockchainHeightSearcher, blockchainHeightBroadcaster, blockSearcher, blockBroadcaster, unconfirmedTransactionsSearcher}
	return &blockchainNetCore
}
