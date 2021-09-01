package transaction

import (
	"helloworld-blockchain-go/application/vo/framwork"
	"helloworld-blockchain-go/dto"
)

type TransactionInputVo struct {
	address                string
	value                  uint64
	inputScript            string
	transactionHash        string
	transactionOutputIndex uint64
}
type TransactionOutputVo struct {
	address                string
	value                  uint64
	outputScript           string
	transactionHash        string
	transactionOutputIndex uint64
}
type TransactionVo struct {
	blockHeight     uint64
	blockHash       string
	confirmCount    uint64
	transactionHash string
	blockTime       string

	transactionFee          uint64
	transactionType         string
	transactionInputCount   uint64
	transactionOutputCount  uint64
	transactionInputValues  uint64
	transactionOutputValues uint64

	transactionInputs  []*TransactionInputVo
	transactionOutputs []*TransactionOutputVo

	inputScripts  []string
	outputScripts []string
}
type TransactionOutputDetailVo struct {
	value           uint64
	spent           bool
	transactionType string

	fromBlockHeight            uint64
	fromBlockHash              string
	fromTransactionHash        string
	fromTransactionOutputIndex uint64
	fromOutputScript           string

	toBlockHeight           uint64
	toBlockHash             string
	toTransactionHash       string
	toTransactionInputIndex uint64
	toInputScript           string

	inputTransaction  TransactionVo
	outputTransaction TransactionVo
}
type UnconfirmedTransactionVo struct {
	transactionHash string
	inputs          []*TransactionInputVo
	outputs         []*TransactionOutputVo
}
type QueryTransactionByTransactionHashRequest struct {
	transactionHash string
}
type QueryTransactionByTransactionHashResponse struct {
	transaction TransactionVo
}
type QueryTransactionOutputByAddressRequest struct {
	address string
}
type QueryTransactionOutputByAddressResponse struct {
	transactionOutputDetail TransactionOutputDetailVo
}
type QueryTransactionOutputByTransactionOutputIdRequest struct {
	transactionHash        string
	transactionOutputIndex uint64
}
type QueryTransactionOutputByTransactionOutputIdResponse struct {
	transactionOutputDetail TransactionOutputDetailVo
}
type QueryTransactionsByBlockHashTransactionHeightRequest struct {
	blockHash     string
	pageCondition framwork.PageCondition
}
type QueryTransactionsByBlockHashTransactionHeightResponse struct {
	transactions []TransactionVo
}
type QueryUnconfirmedTransactionByTransactionHashRequest struct {
	transactionHash string
}
type QueryUnconfirmedTransactionByTransactionHashResponse struct {
	transaction UnconfirmedTransactionVo
}
type QueryUnconfirmedTransactionsRequest struct {
	pageCondition framwork.PageCondition
}
type QueryUnconfirmedTransactionsResponse struct {
	unconfirmedTransactions []*UnconfirmedTransactionVo
}
type SubmitTransactionToBlockchainNetworkRequest struct {
	Transaction *dto.TransactionDto `json:"transaction"`
}
type SubmitTransactionToBlockchainNetworkResponse struct {
	//交易
	Transaction *dto.TransactionDto `json:"transaction"`
	//交易成功提交的节点
	SuccessSubmitNodes []string `json:"successSubmitNodes"`
	//交易提交失败的节点
	FailedSubmitNodes []string `json:"failedSubmitNodes"`
}
