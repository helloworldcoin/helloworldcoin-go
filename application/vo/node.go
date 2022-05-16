package vo

/*
 @author x.king xdotking@gmail.com
*/

type ActiveAutoSearchNodeRequest struct {
}
type ActiveAutoSearchNodeResponse struct {
}
type AddNodeRequest struct {
	Ip string `json:"ip"`
}
type AddNodeResponse struct {
	AddNodeSuccess bool `json:"addNodeSuccess"`
}
type DeactiveAutoSearchNodeRequest struct {
}
type DeactiveAutoSearchNodeResponse struct {
}
type DeleteNodeRequest struct {
	Ip string `json:"ip"`
}
type DeleteNodeResponse struct {
}
type IsAutoSearchNodeRequest struct {
}
type IsAutoSearchNodeResponse struct {
	AutoSearchNode bool `json:"autoSearchNode"`
}
type NodeVo struct {
	Ip               string `json:"ip"`
	BlockchainHeight uint64 `json:"blockchainHeight"`
}
type QueryAllNodesRequest struct {
}
type QueryAllNodesResponse struct {
	Nodes []NodeVo `json:"nodes"`
}
type QueryBlockchainHeightRequest struct {
}
type QueryBlockchainHeightResponse struct {
	BlockchainHeight uint64 `json:"blockchainHeight"`
}
type UpdateNodeRequest struct {
	Ip               string `json:"ip"`
	BlockchainHeight uint64 `json:"blockchainHeight"`
}
type UpdateNodeResponse struct {
}
