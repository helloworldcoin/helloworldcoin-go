package Model

import (
	"strconv"
)

type Block struct {
	Timestamp         uint64
	PreviousBlockHash string
	MerkleTreeRoot    string
	Nonce             string
	Transactions      []Transaction

	Height                    int
	Hash                      string
	TransactionCount          int
	PreviousTransactionHeight int
}
type Transaction struct {
	TransactionHash string
	Inputs          []TransactionInput
	Outputs         []TransactionOutput

	TransactionIndex  int
	TransactionHeight int
	BlockHeight       int
}
type TransactionInput struct {
	TransactionHash        string
	TransactionOutputIndex uint64
	InputScript            []string
}
type TransactionOutput struct {
	Value        uint64
	OutputScript []string

	BlockHeight             int
	BlockHash               string
	TransactionHeight       int
	TransactionHash         string
	TransactionOutputIndex  int
	TransactionIndex        int
	TransactionOutputHeight int
}
type TransactionOutputId struct {
	TransactionHash        string
	TransactionOutputIndex int
}

func (transactionOutputId *TransactionOutputId) GetTransactionOutputId() string {
	return transactionOutputId.TransactionHash + "|" + strconv.Itoa(transactionOutputId.TransactionOutputIndex)
}
