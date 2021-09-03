package controller

import (
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"io"
	"io/ioutil"
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
	s := CreateSuccessResponse("", response)

	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Active()
	var response vo.ActiveMinerResponse
	response.ActiveMinerSuccess = true
	s := CreateSuccessResponse("", response)

	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Deactive()
	var response vo.DeactiveMinerResponse
	response.DeactiveMinerSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) IsAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchBlock := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchBlock()
	var response vo.IsAutoSearchBlockResponse
	response.AutoSearchBlock = isAutoSearchBlock

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchBlock()
	var response vo.ActiveAutoSearchBlockResponse
	response.ActiveAutoSearchBlockSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchBlock()
	var response vo.DeactiveAutoSearchBlockResponse
	response.DeactiveAutoSearchBlockSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) AddNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.AddNodeRequest{}).(*vo.AddNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		//return Response.CreateFailResponse("节点IP不能为空")
	}
	if n.blockchainNetCore.GetNodeService().QueryNode(ip) != nil {
		//return Response.createFailResponse("节点已经存在，不需要重复添加")
	}
	var node model.Node
	node.Ip = ip
	node.BlockchainHeight = 0
	n.blockchainNetCore.GetNodeService().AddNode(&node)
	var response vo.AddNodeResponse
	response.AddNodeSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) UpdateNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.UpdateNodeRequest{}).(*vo.UpdateNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		//return Response.createFailResponse("节点IP不能为空");
	}
	var node model.Node
	node.Ip = ip
	node.BlockchainHeight = request.BlockchainHeight
	n.blockchainNetCore.GetNodeService().UpdateNode(&node)
	var response vo.UpdateNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeleteNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.DeleteNodeRequest{}).(*vo.DeleteNodeRequest)

	n.blockchainNetCore.GetNodeService().DeleteNode(request.Ip)
	var response vo.DeleteNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
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

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) IsAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchNode := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchNode()
	var response vo.IsAutoSearchNodeResponse
	response.AutoSearchNode = isAutoSearchNode

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchNode()
	var response vo.ActiveAutoSearchNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchNode()
	var response vo.DeactiveAutoSearchNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) DeleteBlocks(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.DeleteBlocksRequest{}).(*vo.DeleteBlocksRequest)

	n.blockchainNetCore.GetBlockchainCore().DeleteBlocks(request.BlockHeight)
	var response vo.DeleteBlocksResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) GetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	maxBlockHeight := n.blockchainNetCore.GetBlockchainCore().GetMiner().GetMaxBlockHeight()
	var response vo.GetMaxBlockHeightResponse
	response.MaxBlockHeight = maxBlockHeight

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) SetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.SetMaxBlockHeightRequest{}).(*vo.SetMaxBlockHeightRequest)

	height := request.MaxBlockHeight
	n.blockchainNetCore.GetBlockchainCore().GetMiner().SetMaxBlockHeight(height)
	var response vo.SetMaxBlockHeightResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
