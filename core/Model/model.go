package Model

type Block struct {
	Timestamp         uint64
	PreviousBlockHash string
	MerkleTreeRoot    string
	Nonce             string
	Transactions      []Transaction
}
type Transaction struct {
	TransactionHash string
	Inputs          []TransactionInput
	Outputs         []TransactionOutput
}
type TransactionInput struct {
	TransactionHash        string
	TransactionOutputIndex uint64
	InputScript            []string
}
type TransactionOutput struct {
	Value        uint64
	OutputScript []string
}
