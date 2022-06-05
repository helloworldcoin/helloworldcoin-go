package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core"
	"helloworldcoin-go/core/tool/ResourceTool"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/dao"
	"helloworldcoin-go/netcore/server"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/util/FileUtil"
)

type BlockchainNetCore struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService

	blockchainCore *core.BlockchainCore
	nodeServer     *server.NodeServer

	seedNodeInitializer *SeedNodeInitializer
	nodeSearcher        *NodeSearcher
	nodeBroadcaster     *NodeBroadcaster
	nodeCleaner         *NodeCleaner

	blockchainHeightSearcher    *BlockchainHeightSearcher
	blockchainHeightBroadcaster *BlockchainHeightBroadcaster

	blockSearcher    *BlockSearcher
	blockBroadcaster *BlockBroadcaster

	unconfirmedTransactionsSearcher *UnconfirmedTransactionsSearcher
}

func NewDefaultBlockchainNetCore() *BlockchainNetCore {
	return NewBlockchainNetCore(ResourceTool.GetDataRootPath())
}
func NewBlockchainNetCore(netcorePath string) *BlockchainNetCore {
	netCoreConfiguration := configuration.NewNetCoreConfiguration(netcorePath)

	blockchainCorePath := FileUtil.NewPath(netcorePath, "BlocoreainCore")
	blockchainCore := core.NewBlockchainCore(blockchainCorePath)
	slaveBlockchainCorePath := FileUtil.NewPath(netcorePath, "SlaveBlockchainCore")
	slaveBlockchainCore := core.NewBlockchainCore(slaveBlockchainCorePath)

	nodeDao := dao.NewNodeDao(netCoreConfiguration)
	nodeService := service.NewNodeService(nodeDao)
	nodeServer := server.NewNodeServer(netCoreConfiguration, blockchainCore, nodeService)

	seedNodeInitializer := NewSeedNodeInitializer(netCoreConfiguration, nodeService)
	nodeSearcher := NewNodeSearcher(netCoreConfiguration, nodeService)
	nodeBroadcaster := NewNodeBroadcaster(netCoreConfiguration, nodeService)
	nodeCleaner := NewNodeCleaner(netCoreConfiguration, nodeService)

	blockchainHeightSearcher := NewBlockchainHeightSearcher(netCoreConfiguration, nodeService)
	blockchainHeightBroadcaster := NewBlockchainHeightBroadcaster(netCoreConfiguration, blockchainCore, nodeService)

	blockSearcher := NewBlockSearcher(netCoreConfiguration, blockchainCore, slaveBlockchainCore, nodeService)
	blockBroadcaster := NewBlockBroadcaster(netCoreConfiguration, blockchainCore, nodeService)

	unconfirmedTransactionsSearcher := NewUnconfirmedTransactionsSearcher(netCoreConfiguration, blockchainCore, nodeService)

	blockchainNetCore := NewBlockchainNetCore0(netCoreConfiguration, nodeService, blockchainCore, nodeServer, seedNodeInitializer, nodeSearcher, nodeBroadcaster, nodeCleaner, blockchainHeightSearcher, blockchainHeightBroadcaster, blockSearcher, blockBroadcaster, unconfirmedTransactionsSearcher)
	return blockchainNetCore
}
func NewBlockchainNetCore0(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService,
	blockchainCore *core.BlockchainCore, nodeServer *server.NodeServer, seedNodeInitializer *SeedNodeInitializer, nodeSearcher *NodeSearcher,
	nodeBroadcaster *NodeBroadcaster, nodeCleaner *NodeCleaner,
	blockchainHeightSearcher *BlockchainHeightSearcher, blockchainHeightBroadcaster *BlockchainHeightBroadcaster,
	blockSearcher *BlockSearcher, blockBroadcaster *BlockBroadcaster,
	unconfirmedTransactionsSearcher *UnconfirmedTransactionsSearcher) *BlockchainNetCore {
	var blockchainNetCore BlockchainNetCore
	blockchainNetCore.netCoreConfiguration = netCoreConfiguration
	blockchainNetCore.nodeService = nodeService

	blockchainNetCore.blockchainCore = blockchainCore
	blockchainNetCore.nodeServer = nodeServer

	blockchainNetCore.seedNodeInitializer = seedNodeInitializer
	blockchainNetCore.nodeSearcher = nodeSearcher
	blockchainNetCore.nodeBroadcaster = nodeBroadcaster
	blockchainNetCore.nodeCleaner = nodeCleaner

	blockchainNetCore.blockchainHeightSearcher = blockchainHeightSearcher
	blockchainNetCore.blockchainHeightBroadcaster = blockchainHeightBroadcaster

	blockchainNetCore.blockSearcher = blockSearcher
	blockchainNetCore.blockBroadcaster = blockBroadcaster

	blockchainNetCore.unconfirmedTransactionsSearcher = unconfirmedTransactionsSearcher
	return &blockchainNetCore
}
func (b *BlockchainNetCore) GetBlockchainCore() *core.BlockchainCore {
	return b.blockchainCore
}
func (b *BlockchainNetCore) GetNodeService() *service.NodeService {
	return b.nodeService
}
func (b *BlockchainNetCore) GetNetCoreConfiguration() *configuration.NetCoreConfiguration {
	return b.netCoreConfiguration
}

func (b *BlockchainNetCore) Start() {
	go b.blockchainCore.Start()
	go b.nodeServer.Start()

	go b.seedNodeInitializer.start()
	go b.nodeBroadcaster.start()
	go b.nodeSearcher.start()
	go b.nodeCleaner.start()

	go b.blockchainHeightBroadcaster.start()
	go b.blockchainHeightSearcher.start()

	go b.blockBroadcaster.start()
	go b.blockSearcher.start()

	go b.unconfirmedTransactionsSearcher.start()
}
