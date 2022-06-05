package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/netcore-client/client"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/model"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type NodeSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewNodeSearcher(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *NodeSearcher {
	var nodeSearcher NodeSearcher
	nodeSearcher.netCoreConfiguration = netCoreConfiguration
	nodeSearcher.nodeService = nodeService
	return &nodeSearcher
}

func (b *NodeSearcher) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("'search for nodes in the blockchain network' error.", e)
		}
	}()
	for {
		b.searchNodes()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetNodeSearchTimeInterval())
	}
}

func (b *NodeSearcher) searchNodes() {
	if !b.netCoreConfiguration.IsAutoSearchNode() {
		return
	}

	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		if !b.netCoreConfiguration.IsAutoSearchNode() {
			return
		}
		nodeClient := client.NewNodeClient(node.Ip)
		var getNodesRequest dto.GetNodesRequest
		getNodesResponse := nodeClient.GetNodes(getNodesRequest)
		b.handleGetNodesResponse(getNodesResponse)
	}
}
func (b *NodeSearcher) handleGetNodesResponse(getNodesResponse *dto.GetNodesResponse) {
	if getNodesResponse == nil {
		return
	}
	nodes := getNodesResponse.Nodes
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		if !b.netCoreConfiguration.IsAutoSearchNode() {
			return
		}
		nodeClient := client.NewNodeClient(node.Ip)
		var pingRequest dto.PingRequest
		pingResponse := nodeClient.PingNode(pingRequest)
		if pingResponse != nil {
			var n model.Node
			n.Ip = node.Ip
			n.BlockchainHeight = 0
			b.nodeService.AddNode(&n)
		}
	}
}
