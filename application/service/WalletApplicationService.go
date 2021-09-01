package service

import (
	"helloworld-blockchain-go/application/vo/transaction"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/netcore-client/client"
)

type WalletApplicationService struct {
	blockchainNetCore *netcore.BlockchainNetCore
}

func NewWalletApplicationService(blockchainNetCore *netcore.BlockchainNetCore) *WalletApplicationService {
	var b WalletApplicationService
	b.blockchainNetCore = blockchainNetCore
	return &b
}

func (w *WalletApplicationService) SubmitTransactionToBlockchainNetwork(request *transaction.SubmitTransactionToBlockchainNetworkRequest) *transaction.SubmitTransactionToBlockchainNetworkResponse {
	transactionDto := request.Transaction
	//将交易提交到本地区块链
	w.blockchainNetCore.GetBlockchainCore().PostTransaction(transactionDto)
	//提交交易到网络
	nodes := w.blockchainNetCore.GetNodeService().QueryAllNodes()
	var successSubmitNodes []string
	var failedSubmitNodes []string
	if nodes != nil {
		for _, node := range nodes {
			var postTransactionRequest dto.PostTransactionRequest
			postTransactionRequest.Transaction = transactionDto
			nodeClient := client.NodeClient{Ip: node.Ip}
			postTransactionResponse := nodeClient.PostTransaction(postTransactionRequest)
			if postTransactionResponse != nil {
				successSubmitNodes = append(successSubmitNodes, node.Ip)
			} else {
				failedSubmitNodes = append(failedSubmitNodes, node.Ip)
			}
		}
	}

	var response transaction.SubmitTransactionToBlockchainNetworkResponse
	response.Transaction = transactionDto
	response.SuccessSubmitNodes = successSubmitNodes
	response.FailedSubmitNodes = failedSubmitNodes
	return &response
}
