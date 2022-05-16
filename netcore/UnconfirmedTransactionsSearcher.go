package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core"
	"helloworldcoin-go/netcore-client/client"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type UnconfirmedTransactionsSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	nodeService          *service.NodeService
}

func NewUnconfirmedTransactionsSearcher(netCoreConfiguration *configuration.NetCoreConfiguration, blockchainCore *core.BlockchainCore, nodeService *service.NodeService) *UnconfirmedTransactionsSearcher {
	var unconfirmedTransactionsSearcher UnconfirmedTransactionsSearcher
	unconfirmedTransactionsSearcher.netCoreConfiguration = netCoreConfiguration
	unconfirmedTransactionsSearcher.blockchainCore = blockchainCore
	unconfirmedTransactionsSearcher.nodeService = nodeService
	return &unconfirmedTransactionsSearcher
}

func (u *UnconfirmedTransactionsSearcher) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("在区块链网络中搜寻未确认交易出现异常", e)
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
