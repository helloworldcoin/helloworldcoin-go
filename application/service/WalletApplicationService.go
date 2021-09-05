package service

/*
 @author king 409060350@qq.com
*/
import (
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/core/model"
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

func (w *WalletApplicationService) SubmitTransactionToBlockchainNetwork(request *vo.SubmitTransactionToBlockchainNetworkRequest) *vo.SubmitTransactionToBlockchainNetworkResponse {
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
			nodeClient := client.NewNodeClient(node.Ip)
			postTransactionResponse := nodeClient.PostTransaction(postTransactionRequest)
			if postTransactionResponse != nil {
				successSubmitNodes = append(successSubmitNodes, node.Ip)
			} else {
				failedSubmitNodes = append(failedSubmitNodes, node.Ip)
			}
		}
	}

	var response vo.SubmitTransactionToBlockchainNetworkResponse
	response.Transaction = transactionDto
	response.SuccessSubmitNodes = successSubmitNodes
	response.FailedSubmitNodes = failedSubmitNodes
	return &response
}

func (w *WalletApplicationService) AutomaticBuildTransaction(request *vo.AutomaticBuildTransactionRequest) *vo.AutomaticBuildTransactionResponse {
	var autoBuildTransactionRequest model.AutoBuildTransactionRequest
	autoBuildTransactionRequest.NonChangePayees = payeeVos2payees(request.NonChangePayees)

	autoBuildTransactionResponse := w.blockchainNetCore.GetBlockchainCore().AutoBuildTransaction(&autoBuildTransactionRequest)

	var response vo.AutomaticBuildTransactionResponse
	response.BuildTransactionSuccess = autoBuildTransactionResponse.BuildTransactionSuccess
	response.Message = autoBuildTransactionResponse.Message
	response.TransactionHash = autoBuildTransactionResponse.TransactionHash
	response.Fee = autoBuildTransactionResponse.Fee
	response.Transaction = autoBuildTransactionResponse.Transaction
	response.Payers = payers2payerVos(autoBuildTransactionResponse.Payers)
	response.Payees = payees2payeeVos(autoBuildTransactionResponse.Payees)
	response.ChangePayee = payee2payeeVo(autoBuildTransactionResponse.ChangePayee)
	response.NonChangePayees = payees2payeeVos(autoBuildTransactionResponse.NonChangePayees)

	return &response
}
func payeeVos2payees(payeeVos []*vo.PayeeVo) []*model.Payee {
	var payees []*model.Payee
	if payeeVos != nil {
		for _, payeeVo := range payeeVos {
			payee := payeeVo2payee(payeeVo)
			payees = append(payees, payee)
		}
	}
	return payees
}
func payeeVo2payee(payeeVo *vo.PayeeVo) *model.Payee {
	var payee model.Payee
	payee.Address = payeeVo.Address
	payee.Value = payeeVo.Value
	return &payee
}

func payers2payerVos(payers []*model.Payer) []*vo.PayerVo {
	var payerVos []*vo.PayerVo
	if payers != nil {
		for _, payer := range payers {
			payerVo := payer2payerVo(payer)
			payerVos = append(payerVos, payerVo)
		}
	}
	return payerVos
}
func payer2payerVo(payer *model.Payer) *vo.PayerVo {
	var payerVo vo.PayerVo
	payerVo.Address = payer.Address
	payerVo.PrivateKey = payer.PrivateKey
	payerVo.TransactionHash = payer.TransactionHash
	payerVo.TransactionOutputIndex = payer.TransactionOutputIndex
	payerVo.Value = payer.Value
	return &payerVo
}
func payees2payeeVos(payees []*model.Payee) []*vo.PayeeVo {
	var payeeVos []*vo.PayeeVo
	if payees != nil {
		for _, payee := range payees {
			payeeVo := payee2payeeVo(payee)
			payeeVos = append(payeeVos, payeeVo)
		}
	}
	return payeeVos
}
func payee2payeeVo(payee *model.Payee) *vo.PayeeVo {
	var payeeVo vo.PayeeVo
	payeeVo.Address = payee.Address
	payeeVo.Value = payee.Value
	return &payeeVo
}
