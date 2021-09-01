package miner

type ActiveMinerRequest struct {
}
type ActiveMinerResponse struct {
	ActiveMinerSuccess bool `json:"activeMinerSuccess"`
}

type DeactiveMinerRequest struct {
}
type DeactiveMinerResponse struct {
	DeactiveMinerSuccess bool `json:"deactiveMinerSuccess"`
}

type IsMinerActiveRequest struct {
}
type IsMinerActiveResponse struct {
	MinerInActiveState bool `json:"minerInActiveState"`
}
type GetMaxBlockHeightRequest struct {
}
type GetMaxBlockHeightResponse struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMaxBlockHeightRequest struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMaxBlockHeightResponse struct {
}
