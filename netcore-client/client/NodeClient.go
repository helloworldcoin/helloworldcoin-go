package client

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/dto/API"
	"helloworld-blockchain-go/setting/NetworkSetting"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/NetUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

type NodeClient struct {
	ip string
}

func NewNodeClient(ip string) *NodeClient {
	var nodeClient NodeClient
	nodeClient.ip = ip
	return &nodeClient
}
func (n *NodeClient) GetIp() string {
	return n.ip
}

func (n *NodeClient) PostTransaction(request dto.PostTransactionRequest) *dto.PostTransactionResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.POST_TRANSACTION)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.PostTransactionResponse{}).(*dto.PostTransactionResponse)
}

func (n *NodeClient) PingNode(request dto.PingRequest) *dto.PingResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.PING)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.PingResponse{}).(*dto.PingResponse)
}

func (n *NodeClient) GetBlock(request dto.GetBlockRequest) *dto.GetBlockResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.GET_BLOCK)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.GetBlockResponse{}).(*dto.GetBlockResponse)
}

func (n *NodeClient) GetNodes(request dto.GetNodesRequest) *dto.GetNodesResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.GET_NODES)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.GetNodesResponse{}).(*dto.GetNodesResponse)
}

func (n *NodeClient) PostBlock(request dto.PostBlockRequest) *dto.PostBlockResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.POST_BLOCK)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.PostBlockResponse{}).(*dto.PostBlockResponse)
}

func (n *NodeClient) PostBlockchainHeight(request dto.PostBlockchainHeightRequest) *dto.PostBlockchainHeightResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.POST_BLOCKCHAIN_HEIGHT)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.PostBlockchainHeightResponse{}).(*dto.PostBlockchainHeightResponse)
}

func (n *NodeClient) GetBlockchainHeight(request dto.GetBlockchainHeightRequest) *dto.GetBlockchainHeightResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.GET_BLOCKCHAIN_HEIGHT)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.GetBlockchainHeightResponse{}).(*dto.GetBlockchainHeightResponse)
}

func (n *NodeClient) GetUnconfirmedTransactions(request dto.GetUnconfirmedTransactionsRequest) *dto.GetUnconfirmedTransactionsResponse {
	defer func() {}()
	requestUrl := n.getUrl(API.GET_UNCONFIRMED_TRANSACTIONS)
	requestBody := JsonUtil.ToString(request)
	responseHtml := NetUtil.Get(requestUrl, requestBody)
	return JsonUtil.ToObject(responseHtml, dto.GetUnconfirmedTransactionsResponse{}).(*dto.GetUnconfirmedTransactionsResponse)
}
func (n *NodeClient) getUrl(api string) string {
	return "http://" + n.ip + ":" + StringUtil.ValueOfUint64(NetworkSetting.PORT) + api
}
