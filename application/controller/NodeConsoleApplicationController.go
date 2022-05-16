package controller

/*
 @author x.king xdotking@gmail.com
*/
import (
	"helloworldcoin-go/application/vo"
	"helloworldcoin-go/netcore"
	"helloworldcoin-go/netcore/model"
	"net/http"
)

type NodeConsoleApplicationController struct {
	blockchainNetCore *netcore.BlockchainNetCore
}

func NewNodeConsoleApplicationController(blockchainNetCore *netcore.BlockchainNetCore) *NodeConsoleApplicationController {
	var b NodeConsoleApplicationController
	b.blockchainNetCore = blockchainNetCore
	return &b
}

func (n *NodeConsoleApplicationController) IsMinerActive(rw http.ResponseWriter, req *http.Request) {
	isMinerActive := n.blockchainNetCore.GetBlockchainCore().GetMiner().IsActive()

	var response vo.IsMinerActiveResponse
	response.MinerInActiveState = isMinerActive

	success(rw, response)
}
func (n *NodeConsoleApplicationController) ActiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Active()
	var response vo.ActiveMinerResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) DeactiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Deactive()
	var response vo.DeactiveMinerResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) IsAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchBlock := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchBlock()
	var response vo.IsAutoSearchBlockResponse
	response.AutoSearchBlock = isAutoSearchBlock

	success(rw, response)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchBlock()
	var response vo.ActiveAutoSearchBlockResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchBlock()
	var response vo.DeactiveAutoSearchBlockResponse

	success(rw, response)
}

func (n *NodeConsoleApplicationController) AddNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.AddNodeRequest{}).(*vo.AddNodeRequest)

	ip := request.Ip
	if n.blockchainNetCore.GetNodeService().QueryNode(ip) != nil {
		//节点存在，认为是成功添加。
		var response vo.AddNodeResponse
		response.AddNodeSuccess = true
		success(rw, response)
		return
	}
	var node model.Node
	node.Ip = ip
	node.BlockchainHeight = 0
	n.blockchainNetCore.GetNodeService().AddNode(&node)
	var response vo.AddNodeResponse
	response.AddNodeSuccess = true

	success(rw, response)
}
func (n *NodeConsoleApplicationController) UpdateNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.UpdateNodeRequest{}).(*vo.UpdateNodeRequest)

	var node model.Node
	node.Ip = request.Ip
	node.BlockchainHeight = request.BlockchainHeight
	n.blockchainNetCore.GetNodeService().UpdateNode(&node)
	var response vo.UpdateNodeResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) DeleteNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteNodeRequest{}).(*vo.DeleteNodeRequest)

	n.blockchainNetCore.GetNodeService().DeleteNode(request.Ip)
	var response vo.DeleteNodeResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) QueryAllNodes(rw http.ResponseWriter, req *http.Request) {
	nodes := n.blockchainNetCore.GetNodeService().QueryAllNodes()

	var nodeVos []vo.NodeVo
	if nodes != nil {
		for _, node := range nodes {
			var nodeVo vo.NodeVo
			nodeVo.Ip = node.Ip
			nodeVo.BlockchainHeight = node.BlockchainHeight
			nodeVos = append(nodeVos, nodeVo)
		}
	}

	var response vo.QueryAllNodesResponse
	response.Nodes = nodeVos

	success(rw, response)
}

func (n *NodeConsoleApplicationController) IsAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchNode := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchNode()
	var response vo.IsAutoSearchNodeResponse
	response.AutoSearchNode = isAutoSearchNode

	success(rw, response)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchNode()
	var response vo.ActiveAutoSearchNodeResponse

	success(rw, response)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchNode()
	var response vo.DeactiveAutoSearchNodeResponse

	success(rw, response)
}

func (n *NodeConsoleApplicationController) DeleteBlocks(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteBlocksRequest{}).(*vo.DeleteBlocksRequest)

	n.blockchainNetCore.GetBlockchainCore().DeleteBlocks(request.BlockHeight)
	var response vo.DeleteBlocksResponse

	success(rw, response)
}

func (n *NodeConsoleApplicationController) GetMinerMineMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	maxBlockHeight := n.blockchainNetCore.GetBlockchainCore().GetMiner().GetMinerMineMaxBlockHeight()
	var response vo.GetMinerMineMaxBlockHeightResponse
	response.MaxBlockHeight = maxBlockHeight

	success(rw, response)
}
func (n *NodeConsoleApplicationController) SetMinerMineMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SetMinerMineMaxBlockHeightRequest{}).(*vo.SetMinerMineMaxBlockHeightRequest)

	height := request.MaxBlockHeight
	n.blockchainNetCore.GetBlockchainCore().GetMiner().SetMinerMineMaxBlockHeight(height)
	var response vo.SetMinerMineMaxBlockHeightResponse

	success(rw, response)
}
