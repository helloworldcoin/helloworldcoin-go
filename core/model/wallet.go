package model

/*
 @author king 409060350@qq.com
*/

import "helloworld-blockchain-go/dto"

type Payer struct {
	PrivateKey             string `json:"privateKey"`
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
	Value                  uint64 `json:"value"`
	Address                string `json:"address"`
}
type Payee struct {
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

//TODO web层不应该直接调用
type AutoBuildTransactionRequest struct {
	NonChangePayees []Payee `json:"nonChangePayees"`
}
type AutoBuildTransactionResponse struct {
	BuildTransactionSuccess bool               `json:"buildTransactionSuccess"`
	Message                 string             `json:"message"`
	TransactionHash         string             `json:"transactionHash"`
	Fee                     uint64             `json:"fee"`
	Payers                  []Payer            `json:"payers"`
	NonChangePayees         []Payee            `json:"nonChangePayees"`
	ChangePayee             Payee              `json:"changePayee"`
	Payees                  []Payee            `json:"payees"`
	Transaction             dto.TransactionDto `json:"transaction"`
}
