package netcore

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/server"
	"helloworld-blockchain-go/netcore/service"
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

func NewBlockchainNetCore(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService,
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

func (b *BlockchainNetCore) Start() {
	//启动本地的单机区块链
	go b.blockchainCore.Start()
	//启动区块链节点服务器
	go b.nodeServer.Start()

	//种子节点初始化器
	go b.seedNodeInitializer.start()
	//启动节点广播器
	go b.nodeBroadcaster.start()
	//启动节点搜寻器
	go b.nodeSearcher.start()
	//启动节点清理器
	go b.nodeCleaner.start()

	//启动区块链高度广播器
	go b.blockchainHeightBroadcaster.start()
	//启动区块链高度搜索器
	go b.blockchainHeightSearcher.start()

	//启动区块广播器
	go b.blockBroadcaster.start()
	//启动区块搜寻器
	go b.blockSearcher.start()

	//未确认交易搜索器
	go b.unconfirmedTransactionsSearcher.start()
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
