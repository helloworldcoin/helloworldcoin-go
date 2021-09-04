package netcore

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type UnconfirmedTransactionsSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	nodeService          *service.NodeService
}

func (u *UnconfirmedTransactionsSearcher) start() {
	defer func() {
		if err := recover(); err != nil {
			SystemUtil.ErrorExit("在区块链网络中搜寻未确认交易出现异常", err)
		}
	}()
	for {
		u.searchUnconfirmedTransactions()
		ThreadUtil.MillisecondSleep(u.netCoreConfiguration.GetSearchUnconfirmedTransactionsTimeInterval())
	}
}

func (u *UnconfirmedTransactionsSearcher) searchUnconfirmedTransactions() {
	nodes := u.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		nodeClient := client.NewNodeClient(node.Ip)
		var getUnconfirmedTransactionsRequest dto.GetUnconfirmedTransactionsRequest
		getUnconfirmedTransactionsResponse := nodeClient.GetUnconfirmedTransactions(getUnconfirmedTransactionsRequest)
		if getUnconfirmedTransactionsResponse == nil {
			continue
		}
		transactions := getUnconfirmedTransactionsResponse.Transactions
		if transactions == nil {
			continue
		}
		for _, transaction := range transactions {
			u.blockchainCore.GetUnconfirmedTransactionDatabase().InsertTransaction(transaction)
		}
	}
}
