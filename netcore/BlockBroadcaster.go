package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/netcore-client/client"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type BlockBroadcaster struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	nodeService          *service.NodeService
}

func NewBlockBroadcaster(netCoreConfiguration *configuration.NetCoreConfiguration, blockchainCore *core.BlockchainCore, nodeService *service.NodeService) *BlockBroadcaster {
	var blockBroadcaster BlockBroadcaster
	blockBroadcaster.netCoreConfiguration = netCoreConfiguration
	blockBroadcaster.blockchainCore = blockchainCore
	blockBroadcaster.nodeService = nodeService
	return &blockBroadcaster
}

func (b *BlockBroadcaster) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("在区块链网络中广播自己的区块出现异常", e)
		}
	}()
	for {
		b.broadcastBlock()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetBlockBroadcastTimeInterval())
	}
}

func (b *BlockBroadcaster) broadcastBlock() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		block := b.blockchainCore.QueryTailBlock()
		if block == nil {
			return
		}
		if block.Height <= node.BlockchainHeight {
			continue
		}
		blockDto := Model2DtoTool.Block2BlockDto(block)
		nodeClient := client.NewNodeClient(node.Ip)
		var postBlockRequest dto.PostBlockRequest
		postBlockRequest.Block = blockDto
		nodeClient.PostBlock(postBlockRequest)
	}
}
