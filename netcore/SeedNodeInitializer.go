package netcore

import (
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/setting/NetworkSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type SeedNodeInitializer struct {
	netCoreConfiguration service.NetCoreConfiguration
	nodeService          service.NodeService
}

func (s *SeedNodeInitializer) start() {
	defer func() {
		if err := recover(); err != nil {
			SystemUtil.ErrorExit("定时将种子节点加入区块链网络出现异常", err)
		}
	}()
	for {
		s.addSeedNodes()
		ThreadUtil.MillisecondSleep(s.netCoreConfiguration.GetAddSeedNodeTimeInterval())
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
