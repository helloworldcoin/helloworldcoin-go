package Model

import (
	"helloworld-blockchain-go/core/Model/Script"
	"helloworld-blockchain-go/core/Model/TransactionType"
)

type Block struct {
	Timestamp      uint64
	PreviousHash   string
	MerkleTreeRoot string
	Nonce          string
	Transactions   []*Transaction

	Height                    uint64
	Hash                      string
	TransactionCount          uint64
	PreviousTransactionHeight uint64
	Difficulty                string
}
type Transaction struct {
	TransactionHash string
	Inputs          []*TransactionInput
	Outputs         []*TransactionOutput

	TransactionIndex  uint64
	TransactionHeight uint64
	BlockHeight       uint64

	TransactionType TransactionType.TransactionType
}
type TransactionInput struct {
	UnspentTransactionOutput TransactionOutput
	InputScript              Script.InputScript
}
type TransactionOutput struct {
	Value        uint64
	OutputScript Script.OutputScript

	BlockHeight             uint64
	BlockHash               string
	TransactionHeight       uint64
	TransactionHash         string
	TransactionOutputIndex  uint64
	TransactionIndex        uint64
	TransactionOutputHeight uint64
	Address                 string
}
