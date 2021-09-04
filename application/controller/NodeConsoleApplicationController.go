package controller

/*
 @author king 409060350@qq.com
*/
import (
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/util/StringUtil"
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

func (n *NodeConsoleApplicationController) IsMineActive(rw http.ResponseWriter, req *http.Request) {
	isMineActive := n.blockchainNetCore.GetBlockchainCore().GetMiner().IsActive()

	var response vo.IsMinerActiveResponse
	response.MinerInActiveState = isMineActive

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) ActiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Active()
	var response vo.ActiveMinerResponse
	response.ActiveMinerSuccess = true

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) DeactiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Deactive()
	var response vo.DeactiveMinerResponse
	response.DeactiveMinerSuccess = true

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) IsAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchBlock := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchBlock()
	var response vo.IsAutoSearchBlockResponse
	response.AutoSearchBlock = isAutoSearchBlock

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchBlock()
	var response vo.ActiveAutoSearchBlockResponse
	response.ActiveAutoSearchBlockSuccess = true

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchBlock()
	var response vo.DeactiveAutoSearchBlockResponse
	response.DeactiveAutoSearchBlockSuccess = true

	SuccessHttpResponse(rw, "", response)
}

func (n *NodeConsoleApplicationController) AddNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.AddNodeRequest{}).(*vo.AddNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		FailedHttpResponse(rw, "节点IP不能为空。")
		return
	}
	if n.blockchainNetCore.GetNodeService().QueryNode(ip) != nil {
		FailedHttpResponse(rw, "节点已经存在，不需要重复添加。")
		return
	}
	var node model.Node
	node.Ip = ip
	node.BlockchainHeight = 0
	n.blockchainNetCore.GetNodeService().AddNode(&node)
	var response vo.AddNodeResponse
	response.AddNodeSuccess = true

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) UpdateNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.UpdateNodeRequest{}).(*vo.UpdateNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		FailedHttpResponse(rw, "节点IP不能为空。")
		return
	}
	var node model.Node
	node.Ip = ip
	node.BlockchainHeight = request.BlockchainHeight
	n.blockchainNetCore.GetNodeService().UpdateNode(&node)
	var response vo.UpdateNodeResponse

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) DeleteNode(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteNodeRequest{}).(*vo.DeleteNodeRequest)

	n.blockchainNetCore.GetNodeService().DeleteNode(request.Ip)
	var response vo.DeleteNodeResponse

	SuccessHttpResponse(rw, "", response)
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

	SuccessHttpResponse(rw, "", response)
}

func (n *NodeConsoleApplicationController) IsAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchNode := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchNode()
	var response vo.IsAutoSearchNodeResponse
	response.AutoSearchNode = isAutoSearchNode

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchNode()
	var response vo.ActiveAutoSearchNodeResponse

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchNode()
	var response vo.DeactiveAutoSearchNodeResponse

	SuccessHttpResponse(rw, "", response)
}

func (n *NodeConsoleApplicationController) DeleteBlocks(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteBlocksRequest{}).(*vo.DeleteBlocksRequest)

	n.blockchainNetCore.GetBlockchainCore().DeleteBlocks(request.BlockHeight)
	var response vo.DeleteBlocksResponse

	SuccessHttpResponse(rw, "", response)
}

func (n *NodeConsoleApplicationController) GetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	maxBlockHeight := n.blockchainNetCore.GetBlockchainCore().GetMiner().GetMaxBlockHeight()
	var response vo.GetMaxBlockHeightResponse
	response.MaxBlockHeight = maxBlockHeight

	SuccessHttpResponse(rw, "", response)
}
func (n *NodeConsoleApplicationController) SetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SetMaxBlockHeightRequest{}).(*vo.SetMaxBlockHeightRequest)

	height := request.MaxBlockHeight
	n.blockchainNetCore.GetBlockchainCore().GetMiner().SetMaxBlockHeight(height)
	var response vo.SetMaxBlockHeightResponse

	SuccessHttpResponse(rw, "", response)
}
