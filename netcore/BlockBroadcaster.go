package netcore

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type BlockBroadcaster struct {
	netCoreConfiguration service.NetCoreConfiguration
	blockchainCore       core.BlockchainCore
	nodeService          service.NodeService
}

func (b *BlockBroadcaster) start() {
	defer func() {
		if err := recover(); err != nil {
			SystemUtil.ErrorExit("在区块链网络中广播自己的区块出现异常", err)
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
		nodeClient := client.NodeClient{node.Ip}
		var postBlockRequest dto.PostBlockRequest
		postBlockRequest.Block = blockDto
		nodeClient.PostBlock(postBlockRequest)
	}
}
