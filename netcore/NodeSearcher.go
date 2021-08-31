package netcore

import (
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type NodeSearcher struct {
	netCoreConfiguration configuration.NetCoreConfiguration
	nodeService          service.NodeService
}

func (b *NodeSearcher) start() {
	defer func() {
		if err := recover(); err != nil {
			SystemUtil.ErrorExit("在区块链网络中搜索新的节点出现异常", err)
		}
	}()
	for {
		b.searchNodes()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetSearchNodeTimeInterval())
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
		nodeClient := client.NodeClient{Ip: node.Ip}
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
		nodeClient := client.NodeClient{Ip: node.Ip}
		var pingRequest dto.PingRequest
		pingResponse := nodeClient.PingNode(pingRequest)
		if pingResponse != nil {
			var n model.Node
			n.Ip = node.Ip
			n.BlockchainHeight = 0
			b.nodeService.AddNode(&n)
			LogUtil.Debug("自动机制发现节点[" + node.Ip + "]，已在节点数据库中添加了该节点。")
		}
	}
}
