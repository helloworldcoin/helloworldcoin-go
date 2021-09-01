package block

type BlockVo2 struct {
	Height              uint64 `json:"height"`
	BlockSize           string `json:"blockSize"`
	TransactionCount    uint64 `json:"transactionCount"`
	MinerIncentiveValue uint64 `json:"minerIncentiveValue"`
	Time                string `json:"time"`
	Hash                string `json:"hash"`
}
