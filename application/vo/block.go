package vo

/*
 @author king 409060350@qq.com
*/

type BlockVo struct {
	Height              uint64 `json:"height"`
	ConfirmCount        uint64 `json:"confirmCount"`
	BlockSize           string `json:"blockSize"`
	TransactionCount    uint64 `json:"transactionCount"`
	Time                string `json:"time"`
	MinerIncentiveValue uint64 `json:"minerIncentiveValue"`

	Difficulty        string `json:"difficulty"`
	Nonce             string `json:"nonce"`
	Hash              string `json:"hash"`
	PreviousBlockHash string `json:"previousBlockHash"`
	NextBlockHash     string `json:"nextBlockHash"`
	MerkleTreeRoot    string `json:"merkleTreeRoot"`
}
type BlockVo2 struct {
	Height              uint64 `json:"height"`
	BlockSize           string `json:"blockSize"`
	TransactionCount    uint64 `json:"transactionCount"`
	MinerIncentiveValue uint64 `json:"minerIncentiveValue"`
	Time                string `json:"time"`
	Hash                string `json:"hash"`
}
type DeleteBlocksRequest struct {
	BlockHeight uint64 `json:"blockHeight"`
}
type DeleteBlocksResponse struct {
}
type QueryBlockByBlockHashRequest struct {
	BlockHash string `json:"blockHash"`
}
type QueryBlockByBlockHashResponse struct {
	Block *BlockVo `json:"block"`
}
type QueryBlockByBlockHeightRequest struct {
	BlockHeight uint64 `json:"blockHeight"`
}
type QueryBlockByBlockHeightResponse struct {
	Block *BlockVo `json:"block"`
}
type QueryTop10BlocksRequest struct {
}
type QueryTop10BlocksResponse struct {
	Blocks []BlockVo2 `json:"blocks"`
}
