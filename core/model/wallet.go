package model

/*
 @author x.king xdotking@gmail.com
*/

import "helloworldcoin-go/netcore-dto/dto"

type Payer struct {
	PrivateKey             string
	TransactionHash        string
	TransactionOutputIndex uint64
	Value                  uint64
	Address                string
}
type Payee struct {
	Address string
	Value   uint64
}
type AutoBuildTransactionRequest struct {
	NonChangePayees []*Payee
}
type AutoBuildTransactionResponse struct {
	BuildTransactionSuccess bool
	TransactionHash         string
	Fee                     uint64
	Payers                  []*Payer
	NonChangePayees         []*Payee
	ChangePayee             *Payee
	Payees                  []*Payee
	Transaction             *dto.TransactionDto
}
