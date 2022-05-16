package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/model"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/setting/NetworkSetting"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type SeedNodeInitializer struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewSeedNodeInitializer(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *SeedNodeInitializer {
	var seedNodeInitializer SeedNodeInitializer
	seedNodeInitializer.netCoreConfiguration = netCoreConfiguration
	seedNodeInitializer.nodeService = nodeService
	return &seedNodeInitializer
}

func (s *SeedNodeInitializer) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("定时将种子节点加入区块链网络出现异常", e)
		}
	}()
	for {
		s.addSeedNodes()
		ThreadUtil.MillisecondSleep(s.netCoreConfiguration.GetSeedNodeInitializeTimeInterval())
	}
}

func (s *SeedNodeInitializer) addSeedNodes() {
	if !s.netCoreConfiguration.IsAutoSearchNode() {
		return
	}

	for _, seedNode := range NetworkSetting.SEED_NODES {
		var node model.Node
		node.Ip = seedNode
		node.BlockchainHeight = 0
		s.nodeService.AddNode(&node)
		LogUtil.Debug("种子节点初始化器提示您:种子节点[" + node.Ip + "]加入了区块链网络。")
	}
}
