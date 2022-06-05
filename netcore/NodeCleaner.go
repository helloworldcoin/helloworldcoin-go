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

type NodeCleaner struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewNodeCleaner(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *NodeCleaner {
	var nodeCleaner NodeCleaner
	nodeCleaner.netCoreConfiguration = netCoreConfiguration
	nodeCleaner.nodeService = nodeService
	return &nodeCleaner
}

func (b *NodeCleaner) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("'clean up dead nodes' error.", e)
		}
	}()
	for {
		b.cleanDeadNodes()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetNodeCleanTimeInterval())
	}
}

func (b *NodeCleaner) cleanDeadNodes() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		nodeClient := client.NewNodeClient(node.Ip)
		var pingRequest dto.PingRequest
		pingResponse := nodeClient.PingNode(pingRequest)
		if pingResponse == nil {
			b.nodeService.DeleteNode(node.Ip)
		}
	}
}
