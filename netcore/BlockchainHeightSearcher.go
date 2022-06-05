package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/netcore-client/client"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type BlockchainHeightSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewBlockchainHeightSearcher(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *BlockchainHeightSearcher {
	var blockchainHeightSearcher BlockchainHeightSearcher
	blockchainHeightSearcher.netCoreConfiguration = netCoreConfiguration
	blockchainHeightSearcher.nodeService = nodeService
	return &blockchainHeightSearcher
}

func (b *BlockchainHeightSearcher) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("'search for node‘s Blockchain Height in the blockchain network' error.", e)
		}
	}()
	for {
		b.searchBlockchainHeight()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetBlockchainHeightSearchTimeInterval())
	}
}

func (b *BlockchainHeightSearcher) searchBlockchainHeight() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		nodeClient := client.NewNodeClient(node.Ip)
		var getBlockchainHeightRequest dto.GetBlockchainHeightRequest
		getBlockchainHeightResponse := nodeClient.GetBlockchainHeight(getBlockchainHeightRequest)
		if getBlockchainHeightResponse != nil {
			node.BlockchainHeight = getBlockchainHeightResponse.BlockchainHeight
			b.nodeService.UpdateNode(node)
		}
	}
}
