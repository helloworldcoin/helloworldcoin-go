package controller

import (
	"helloworld-blockchain-go/application/vo/block"
	"helloworld-blockchain-go/application/vo/miner"
	"helloworld-blockchain-go/application/vo/node"
	"helloworld-blockchain-go/application/vo/synchronizer"
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

	var response miner.IsMinerActiveResponse
	response.MinerInActiveState = isMineActive
	s := CreateSuccessResponse("", response)

	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Active()
	var response miner.ActiveMinerResponse
	response.ActiveMinerSuccess = true
	s := CreateSuccessResponse("", response)

	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveMiner(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetBlockchainCore().GetMiner().Deactive()
	var response miner.DeactiveMinerResponse
	response.DeactiveMinerSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) IsAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchBlock := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchBlock()
	var response synchronizer.IsAutoSearchBlockResponse
	response.AutoSearchBlock = isAutoSearchBlock

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchBlock()
	var response synchronizer.ActiveAutoSearchBlockResponse
	response.ActiveAutoSearchBlockSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchBlock(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchBlock()
	var response synchronizer.DeactiveAutoSearchBlockResponse
	response.DeactiveAutoSearchBlockSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) AddNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), node.AddNodeRequest{}).(*node.AddNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		//return Response.CreateFailResponse("节点IP不能为空")
	}
	if n.blockchainNetCore.GetNodeService().QueryNode(ip) != nil {
		//return Response.createFailResponse("节点已经存在，不需要重复添加")
	}
	var nodeTemp model.Node
	nodeTemp.Ip = ip
	nodeTemp.BlockchainHeight = 0
	n.blockchainNetCore.GetNodeService().AddNode(&nodeTemp)
	var response node.AddNodeResponse
	response.AddNodeSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) UpdateNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), node.UpdateNodeRequest{}).(*node.UpdateNodeRequest)

	ip := request.Ip
	if StringUtil.IsNullOrEmpty(ip) {
		//return Response.createFailResponse("节点IP不能为空");
	}
	var nodeTemp model.Node
	nodeTemp.Ip = ip
	nodeTemp.BlockchainHeight = request.BlockchainHeight
	n.blockchainNetCore.GetNodeService().UpdateNode(&nodeTemp)
	var response node.UpdateNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeleteNode(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), node.DeleteNodeRequest{}).(*node.DeleteNodeRequest)

	n.blockchainNetCore.GetNodeService().DeleteNode(request.Ip)
	var response node.DeleteNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) QueryAllNodes(rw http.ResponseWriter, req *http.Request) {
	nodes := n.blockchainNetCore.GetNodeService().QueryAllNodes()

	var nodeVos []node.NodeVo
	if nodes != nil {
		for _, nodeTemp := range nodes {
			var nodeVo node.NodeVo
			nodeVo.Ip = nodeTemp.Ip
			nodeVo.BlockchainHeight = nodeTemp.BlockchainHeight
			nodeVos = append(nodeVos, nodeVo)
		}
	}

	var response node.QueryAllNodesResponse
	response.Nodes = nodeVos

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) IsAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	isAutoSearchNode := n.blockchainNetCore.GetNetCoreConfiguration().IsAutoSearchNode()
	var response node.IsAutoSearchNodeResponse
	response.AutoSearchNode = isAutoSearchNode

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) ActiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().ActiveAutoSearchNode()
	var response node.ActiveAutoSearchNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) DeactiveAutoSearchNode(rw http.ResponseWriter, req *http.Request) {
	n.blockchainNetCore.GetNetCoreConfiguration().DeactiveAutoSearchNode()
	var response node.DeactiveAutoSearchNodeResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) DeleteBlocks(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), block.DeleteBlocksRequest{}).(*block.DeleteBlocksRequest)

	n.blockchainNetCore.GetBlockchainCore().DeleteBlocks(request.BlockHeight)
	var response block.DeleteBlocksResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (n *NodeConsoleApplicationController) GetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	maxBlockHeight := n.blockchainNetCore.GetBlockchainCore().GetMiner().GetMaxBlockHeight()
	var response miner.GetMaxBlockHeightResponse
	response.MaxBlockHeight = maxBlockHeight

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (n *NodeConsoleApplicationController) SetMaxBlockHeight(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), miner.SetMaxBlockHeightRequest{}).(*miner.SetMaxBlockHeightRequest)

	height := request.MaxBlockHeight
	n.blockchainNetCore.GetBlockchainCore().GetMiner().SetMaxBlockHeight(height)
	var response miner.SetMaxBlockHeightResponse

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
