package dto

type BlockDto struct {
	Timestamp    uint64
	PreviousHash string
	Transactions []TransactionDto
	Nonce        string
}
type TransactionDto struct {
	Inputs  []TransactionInputDto
	Outputs []TransactionOutputDto
}
type TransactionInputDto struct {
	TransactionHash        string
	TransactionOutputIndex uint64
	InputScript            []string
}
type TransactionOutputDto struct {
	OutputScript []string
	Value        uint64
}
