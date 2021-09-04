package BlockchainNetCoreFactory

/*
 @author king 409060350@qq.com
*/
import (
	"helloworld-blockchain-go/core/BlockchainCoreFactory"
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/dao"
	"helloworld-blockchain-go/netcore/server"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/FileUtil"
)

/**
 * 创建[区块链网络版核心]实例
 */
func CreateDefaultBlockchainNetCore() *netcore.BlockchainNetCore {
	return CreateBlockchainNetCore(ResourcePathTool.GetDataRootPath())
}

/**
 * 创建[区块链网络版核心]实例
 *
 * @param netcorePath 区块链数据存放位置
 */
func CreateBlockchainNetCore(netcorePath string) *netcore.BlockchainNetCore {
	netCoreConfiguration := configuration.NewNetCoreConfiguration(netcorePath)

	blockchainCorePath := FileUtil.NewPath(netcorePath, "BlockchainCore")
	blockchainCore := BlockchainCoreFactory.CreateBlockchainCore(blockchainCorePath)
	slaveBlockchainCorePath := FileUtil.NewPath(netcorePath, "SlaveBlockchainCore")
	slaveBlockchainCore := BlockchainCoreFactory.CreateBlockchainCore(slaveBlockchainCorePath)

	nodeDao := dao.NewNodeDao(netCoreConfiguration)
	nodeService := service.NewNodeService(nodeDao)
	nodeServer := server.NewNodeServer(netCoreConfiguration, blockchainCore, nodeService)

	seedNodeInitializer := netcore.NewSeedNodeInitializer(netCoreConfiguration, nodeService)
	nodeSearcher := netcore.NewNodeSearcher(netCoreConfiguration, nodeService)
	nodeBroadcaster := netcore.NewNodeBroadcaster(netCoreConfiguration, nodeService)
	nodeCleaner := netcore.NewNodeCleaner(netCoreConfiguration, nodeService)

	blockchainHeightSearcher := netcore.NewBlockchainHeightSearcher(netCoreConfiguration, nodeService)
	blockchainHeightBroadcaster := netcore.NewBlockchainHeightBroadcaster(netCoreConfiguration, blockchainCore, nodeService)

	blockSearcher := netcore.NewBlockSearcher(netCoreConfiguration, blockchainCore, slaveBlockchainCore, nodeService)
	blockBroadcaster := netcore.NewBlockBroadcaster(netCoreConfiguration, blockchainCore, nodeService)

	unconfirmedTransactionsSearcher := netcore.NewUnconfirmedTransactionsSearcher(netCoreConfiguration, blockchainCore, nodeService)

	blockchainNetCore := netcore.NewBlockchainNetCore(netCoreConfiguration, nodeService, blockchainCore, nodeServer, seedNodeInitializer, nodeSearcher, nodeBroadcaster, nodeCleaner, blockchainHeightSearcher, blockchainHeightBroadcaster, blockSearcher, blockBroadcaster, unconfirmedTransactionsSearcher)
	return blockchainNetCore
}
