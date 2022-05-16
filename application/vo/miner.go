package vo

/*
 @author x.king xdotking@gmail.com
*/

type ActiveMinerRequest struct {
}
type ActiveMinerResponse struct {
}

type DeactiveMinerRequest struct {
}
type DeactiveMinerResponse struct {
}

type IsMinerActiveRequest struct {
}
type IsMinerActiveResponse struct {
	MinerInActiveState bool `json:"minerInActiveState"`
}
type GetMaxBlockHeightRequest struct {
}
type GetMinerMineMaxBlockHeightRequest struct {
}
type GetMinerMineMaxBlockHeightResponse struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMinerMineMaxBlockHeightRequest struct {
	MaxBlockHeight uint64 `json:"maxBlockHeight"`
}
type SetMinerMineMaxBlockHeightResponse struct {
}
