package model

import "helloworld-blockchain-go/dto"

type Payer struct {
	privateKey             string
	transactionHash        string
	transactionOutputIndex uint64
	value                  uint64
	address                string
}
type Payee struct {
	address string
	value   uint64
}
type AutoBuildTransactionRequest struct {
	Payee []Payee
}
type AutoBuildTransactionResponse struct {
	buildTransactionSuccess bool
	message                 string
	transactionHash         string
	fee                     uint64
	payers                  []Payer
	nonChangePayees         []Payee
	changePayee             Payee
	payees                  []Payee
	transaction             dto.TransactionDto
}
