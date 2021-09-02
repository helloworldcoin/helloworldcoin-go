package transaction

import (
	"helloworld-blockchain-go/application/vo/framwork"
	"helloworld-blockchain-go/dto"
)

type TransactionInputVo struct {
	Address                string `json:"address"`
	Value                  uint64 `json:"value"`
	InputScript            string `json:"inputScript"`
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
}
type TransactionOutputVo struct {
	Address                string `json:"address"`
	Value                  uint64 `json:"value"`
	OutputScript           string `json:"outputScript"`
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
}
type TransactionVo struct {
	BlockHeight     uint64 `json:"blockHeight"`
	BlockHash       string `json:"blockHash"`
	ConfirmCount    uint64 `json:"confirmCount"`
	TransactionHash string `json:"transactionHash"`
	BlockTime       string `json:"blockTime"`

	TransactionFee          uint64 `json:"transactionFee"`
	TransactionType         string `json:"transactionType"`
	TransactionInputCount   uint64 `json:"transactionInputCount"`
	TransactionOutputCount  uint64 `json:"transactionOutputCount"`
	TransactionInputValues  uint64 `json:"transactionInputValues"`
	TransactionOutputValues uint64 `json:"transactionOutputValues"`

	TransactionInputs  []*TransactionInputVo  `json:"transactionInputs"`
	TransactionOutputs []*TransactionOutputVo `json:"transactionOutputs"`

	InputScripts  []string `json:"inputScripts"`
	OutputScripts []string `json:"outputScripts"`
}
type TransactionOutputDetailVo struct {
	Value           uint64 `json:"value"`
	Spent           bool   `json:"spent"`
	TransactionType string `json:"transactionType"`

	FromBlockHeight            uint64 `json:"fromBlockHeight"`
	FromBlockHash              string `json:"fromBlockHash"`
	FromTransactionHash        string `json:"fromTransactionHash"`
	FromTransactionOutputIndex uint64 `json:"fromTransactionOutputIndex"`
	FromOutputScript           string `json:"fromOutputScript"`

	ToBlockHeight           uint64 `json:"toBlockHeight"`
	ToBlockHash             string `json:"toBlockHash"`
	ToTransactionHash       string `json:"toTransactionHash"`
	ToTransactionInputIndex uint64 `json:"toTransactionInputIndex"`
	ToInputScript           string `json:"toInputScript"`

	InputTransaction  *TransactionVo `json:"inputTransaction"`
	OutputTransaction *TransactionVo `json:"outputTransaction"`
}
type UnconfirmedTransactionVo struct {
	TransactionHash string                  `json:"transactionHash"`
	Inputs          []*TransactionInputVo2  `json:"inputs"`
	Outputs         []*TransactionOutputVo2 `json:"outputs"`
}
type TransactionInputVo2 struct {
	Value                  uint64 `json:"value"`
	Address                string `json:"address"`
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
}
type TransactionOutputVo2 struct {
	Value   uint64 `json:"value"`
	Address string `json:"address"`
}
type QueryTransactionByTransactionHashRequest struct {
	TransactionHash string `json:"transactionHash"`
}
type QueryTransactionByTransactionHashResponse struct {
	Transaction *TransactionVo `json:"transaction"`
}
type QueryTransactionOutputByAddressRequest struct {
	Address string `json:"address"`
}
type QueryTransactionOutputByAddressResponse struct {
	TransactionOutputDetail *TransactionOutputDetailVo `json:"transactionOutputDetail"`
}
type QueryTransactionOutputByTransactionOutputIdRequest struct {
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
}
type QueryTransactionOutputByTransactionOutputIdResponse struct {
	TransactionOutputDetail *TransactionOutputDetailVo `json:"transactionOutputDetail"`
}
type QueryTransactionsByBlockHashTransactionHeightRequest struct {
	BlockHash     string                  `json:"blockHash"`
	PageCondition *framwork.PageCondition `json:"pageCondition"`
}
type QueryTransactionsByBlockHashTransactionHeightResponse struct {
	Transactions []*TransactionVo `json:"transactions"`
}
type QueryUnconfirmedTransactionByTransactionHashRequest struct {
	TransactionHash string `json:"transactionHash"`
}
type QueryUnconfirmedTransactionByTransactionHashResponse struct {
	Transaction *UnconfirmedTransactionVo `json:"transaction"`
}
type QueryUnconfirmedTransactionsRequest struct {
	PageCondition *framwork.PageCondition `json:"pageCondition"`
}
type QueryUnconfirmedTransactionsResponse struct {
	UnconfirmedTransactions []*UnconfirmedTransactionVo `json:"unconfirmedTransactions"`
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
