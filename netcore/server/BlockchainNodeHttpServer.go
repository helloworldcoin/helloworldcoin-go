package server

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/dto/API"
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
	http.HandleFunc("/hello", b.HelloServer)
	http.HandleFunc(API.GET_BLOCK, b.getBlock)
	http.HandleFunc(API.POST_BLOCK, b.postBlock)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (b *BlockchainNodeHttpServer) HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func (b *BlockchainNodeHttpServer) getBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetBlockRequest{}).(*dto.GetBlockRequest)
	blockByBlockHeight := b.BlockchainCore.QueryBlockByBlockHeight(request.BlockHeight)
	block := Model2DtoTool.Block2BlockDto(blockByBlockHeight)
	var response dto.GetBlockResponse
	response.Block = *block
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *BlockchainNodeHttpServer) postBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostBlockRequest{}).(*dto.PostBlockRequest)
	b.BlockchainCore.AddBlockDto(&request.Block)
	var response dto.PostBlockResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
