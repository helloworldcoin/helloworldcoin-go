package dao

import (
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/netcore/po"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

type NodeDao struct {
	netCoreConfiguration *service.NetCoreConfiguration
}

const NODE_DATABASE_NAME = "NodeDatabase"

func (n NodeDao) QueryNode(ip string) *po.NodePo {
	nodePos := n.QueryAllNodes()
	if nodePos != nil {
		for _, nodePo := range nodePos {
			if StringUtil.IsEquals(ip, nodePo.Ip) {
				return nodePo
			}
		}
	}
	return nil
}

func (n NodeDao) AddNode(node *po.NodePo) {
	KvDbUtil.Put(n.getNodeDatabasePath(), n.getKeyByNodePo(node), n.encode(node))
}

func (n NodeDao) UpdateNode(node *po.NodePo) {
	KvDbUtil.Put(n.getNodeDatabasePath(), n.getKeyByNodePo(node), n.encode(node))
}

func (n NodeDao) DeleteNode(ip string) {
	KvDbUtil.Delete(n.getNodeDatabasePath(), n.getKeyByIp(ip))
}

func (n NodeDao) QueryAllNodes() []*po.NodePo {
	var nodePos []*po.NodePo
	//获取所有
	bytesNodePos := KvDbUtil.Gets(n.getNodeDatabasePath(), 1, 100000000)
	if bytesNodePos != nil {
		for e := bytesNodePos.Front(); e != nil; e = e.Next() {
			nodePo := n.decodeToNodePo(e.Value.([]byte))
			nodePos = append(nodePos, &nodePo)
		}
	}
	return nodePos
}
func (n NodeDao) getNodeDatabasePath() string {
	return FileUtil.NewPath(n.netCoreConfiguration.NetCorePath, NODE_DATABASE_NAME)
}
func (n NodeDao) getKeyByNodePo(node *po.NodePo) []byte {
	return n.getKeyByIp(node.Ip)
}
func (n NodeDao) getKeyByIp(ip string) []byte {
	return ByteUtil.StringToUtf8Bytes(ip)
}
func (n NodeDao) encode(node *po.NodePo) []byte {
	return ByteUtil.StringToUtf8Bytes(JsonUtil.ToString(*node))
}
func (n NodeDao) decodeToNodePo(bytesNodePo []byte) po.NodePo {
	return JsonUtil.ToObject(ByteUtil.Utf8BytesToString(bytesNodePo), po.NodePo{}).(po.NodePo)
}
