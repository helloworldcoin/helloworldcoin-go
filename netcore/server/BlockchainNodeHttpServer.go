package server

import (
	"fmt"
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/dto/API"
	"helloworld-blockchain-go/setting/BlockSetting"
	"helloworld-blockchain-go/util/JsonUtil"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type BlockchainNodeHttpServer struct {
	BlockchainCore *core.BlockchainCore
}

func (b *BlockchainNodeHttpServer) Start() {
	http.HandleFunc(API.PING, b.ping)
	http.HandleFunc(API.GET_NODES, b.getNodes)
	http.HandleFunc(API.GET_BLOCK, b.getBlock)
	http.HandleFunc(API.POST_BLOCK, b.postBlock)
	http.HandleFunc(API.POST_BLOCKCHAIN_HEIGHT, b.postBlockchainHeight)
	http.HandleFunc(API.GET_BLOCKCHAIN_HEIGHT, b.getBlockchainHeight)
	http.HandleFunc(API.POST_TRANSACTION, b.postTransaction)
	http.HandleFunc(API.GET_UNCONFIRMED_TRANSACTIONS, b.getUnconfirmedTransactions)

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (b *BlockchainNodeHttpServer) getBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetBlockRequest{}).(*dto.GetBlockRequest)
	blockByBlockHeight := b.BlockchainCore.QueryBlockByBlockHeight(request.BlockHeight)
	block := Model2DtoTool.Block2BlockDto(blockByBlockHeight)
	var response dto.GetBlockResponse
	response.Block = block
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) postBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostBlockRequest{}).(*dto.PostBlockRequest)
	b.BlockchainCore.AddBlockDto(request.Block)
	var response dto.PostBlockResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) postBlockchainHeight(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostBlockchainHeightRequest{}).(*dto.PostBlockchainHeightRequest)

	fmt.Println(request)
	/*	var node Node
		node.setIp(requestIp)
		node.setBlockchainHeight(request.getBlockchainHeight())
		nodeService.updateNode(node)*/
	var response dto.PostBlockchainHeightResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) getBlockchainHeight(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetBlockchainHeightRequest{}).(*dto.GetBlockchainHeightRequest)
	fmt.Println(request)
	blockchainHeight := b.BlockchainCore.QueryBlockchainHeight()
	var response dto.GetBlockchainHeightResponse
	response.BlockchainHeight = blockchainHeight
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) postTransaction(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostTransactionRequest{}).(*dto.PostTransactionRequest)
	b.BlockchainCore.PostTransaction(request.Transaction)
	var response dto.PostTransactionResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) getUnconfirmedTransactions(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetUnconfirmedTransactionsRequest{}).(*dto.GetUnconfirmedTransactionsRequest)
	fmt.Println(request)
	unconfirmedTransactionDatabase := b.BlockchainCore.UnconfirmedTransactionDatabase
	transactions := unconfirmedTransactionDatabase.SelectTransactions(1, BlockSetting.BLOCK_MAX_TRANSACTION_COUNT)
	var response dto.GetUnconfirmedTransactionsResponse
	response.Transactions = transactions
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}

func (b *BlockchainNodeHttpServer) ping(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PingRequest{}).(*dto.PingRequest)
	fmt.Println(request)
	//TODO
	var response dto.PingResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}

func (b *BlockchainNodeHttpServer) getNodes(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetNodesRequest{}).(*dto.GetNodesRequest)
	fmt.Println(request)
	//TODO
	var response dto.GetNodesResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
