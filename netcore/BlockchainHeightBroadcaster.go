package netcore

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type BlockchainHeightBroadcaster struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	nodeService          *service.NodeService
}

func NewBlockchainHeightBroadcaster(netCoreConfiguration *configuration.NetCoreConfiguration, blockchainCore *core.BlockchainCore, nodeService *service.NodeService) *BlockchainHeightBroadcaster {
	var blockchainHeightBroadcaster BlockchainHeightBroadcaster
	blockchainHeightBroadcaster.netCoreConfiguration = netCoreConfiguration
	blockchainHeightBroadcaster.blockchainCore = blockchainCore
	blockchainHeightBroadcaster.nodeService = nodeService
	return &blockchainHeightBroadcaster
}

func (b *BlockchainHeightBroadcaster) start() {
	defer func() {
		if e := recover(); e != nil {
			SystemUtil.ErrorExit("在区块链网络中广播自身区块链高度异常", e)
		}
	}()
	for {
		b.broadcastBlockchainHeight()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetBlockchainHeightBroadcastTimeInterval())
	}
}

func (b *BlockchainHeightBroadcaster) broadcastBlockchainHeight() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		blockchainHeight := b.blockchainCore.QueryBlockchainHeight()
		if blockchainHeight <= node.BlockchainHeight {
			continue
		}
		nodeClient := client.NewNodeClient(node.Ip)
		var postBlockchainHeightRequest dto.PostBlockchainHeightRequest
		postBlockchainHeightRequest.BlockchainHeight = blockchainHeight
		nodeClient.PostBlockchainHeight(postBlockchainHeightRequest)
	}
}
