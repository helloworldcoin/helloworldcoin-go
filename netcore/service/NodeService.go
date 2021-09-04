package service

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/netcore/dao"
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/netcore/po"
)

type NodeService struct {
	nodeDao *dao.NodeDao
}

func NewNodeService(nodeDao *dao.NodeDao) *NodeService {
	return &NodeService{nodeDao}
}

func (n *NodeService) DeleteNode(ip string) {
	n.nodeDao.DeleteNode(ip)
}

func (n *NodeService) QueryAllNodes() []*model.Node {
	nodePos := n.nodeDao.QueryAllNodes()
	return n.nodePo2Nodes(nodePos)
}

func (n *NodeService) AddNode(node *model.Node) {
	if n.nodeDao.QueryNode(node.Ip) != nil {
		return
	}
	nodePo := n.node2NodePo(node)
	n.nodeDao.AddNode(nodePo)
}

func (n *NodeService) UpdateNode(node *model.Node) {
	nodePo := n.nodeDao.QueryNode(node.Ip)
	if nodePo == nil {
		return
	}
	nodePo = n.node2NodePo(node)
	n.nodeDao.UpdateNode(nodePo)
}

func (n *NodeService) QueryNode(ip string) *model.Node {
	nodePo := n.nodeDao.QueryNode(ip)
	if nodePo == nil {
		return nil
	}
	return n.nodePo2Node(nodePo)
}

func (n *NodeService) nodePo2Nodes(nodePos []*po.NodePo) []*model.Node {
	var nodes []*model.Node
	if nodePos != nil {
		for _, nodePo := range nodePos {
			node := n.nodePo2Node(nodePo)
			nodes = append(nodes, node)
		}
	}
	return nodes
}
func (n *NodeService) nodePo2Node(nodePo *po.NodePo) *model.Node {
	var node model.Node
	node.Ip = nodePo.Ip
	node.BlockchainHeight = nodePo.BlockchainHeight
	return &node
}
func (n *NodeService) node2NodePo(node *model.Node) *po.NodePo {
	var nodePo po.NodePo
	nodePo.Ip = node.Ip
	nodePo.BlockchainHeight = node.BlockchainHeight
	return &nodePo
}
