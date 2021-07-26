package dto

type BlockDto struct {
	Timestamp    uint64           `json:"timestamp"`
	PreviousHash string           `json:"previousHash"`
	Transactions []TransactionDto `json:"transactions"`
	Nonce        string           `json:"nonce"`
}
type TransactionDto struct {
	Inputs  []TransactionInputDto  `json:"inputs"`
	Outputs []TransactionOutputDto `json:"outputs"`
}
type TransactionInputDto struct {
	TransactionHash        string         `json:"transactionHash"`
	TransactionOutputIndex uint64         `json:"transactionOutputIndex"`
	InputScript            InputScriptDto `json:"inputScript"`
}
type TransactionOutputDto struct {
	OutputScript OutputScriptDto `json:"outputScript"`
	Value        uint64          `json:"value"`
}
type ScriptDto = []string
type InputScriptDto = ScriptDto
type OutputScriptDto = ScriptDto

type GetBlockRequest struct {
	BlockHeight uint64 `json:"blockHeight"`
}
type GetBlockResponse struct {
	Block BlockDto `json:"block"`
}
type PostBlockRequest struct {
	Block BlockDto `json:"block"`
}
type PostBlockResponse struct {
}
