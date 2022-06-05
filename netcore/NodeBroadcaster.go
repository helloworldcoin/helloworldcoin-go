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

type NodeBroadcaster struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewNodeBroadcaster(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *NodeBroadcaster {
	var nodeBroadcaster NodeBroadcaster
	nodeBroadcaster.netCoreConfiguration = netCoreConfiguration
	nodeBroadcaster.nodeService = nodeService
	return &nodeBroadcaster
}

func (b *NodeBroadcaster) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("'broadcasts itself to the whole network' error.", e)
		}
	}()
	for {
		b.broadcastNode()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetNodeBroadcastTimeInterval())
	}
}

func (b *NodeBroadcaster) broadcastNode() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		nodeClient := client.NewNodeClient(node.Ip)
		var pingRequest dto.PingRequest
		nodeClient.PingNode(pingRequest)
	}
}
