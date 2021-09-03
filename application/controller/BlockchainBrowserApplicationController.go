package controller

import (
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/TimeUtil"
	"io"
	"io/ioutil"
	"net/http"
)

type BlockchainBrowserApplicationController struct {
	blockchainNetCore                   *netcore.BlockchainNetCore
	blockchainBrowserApplicationService *service.BlockchainBrowserApplicationService
}

func NewBlockchainBrowserApplicationController(blockchainNetCore *netcore.BlockchainNetCore, blockchainBrowserApplicationService *service.BlockchainBrowserApplicationService) *BlockchainBrowserApplicationController {
	var b BlockchainBrowserApplicationController
	b.blockchainNetCore = blockchainNetCore
	b.blockchainBrowserApplicationService = blockchainBrowserApplicationService
	return &b
}

func (b *BlockchainBrowserApplicationController) QueryTransactionByTransactionHash(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryTransactionByTransactionHashRequest{}).(*vo.QueryTransactionByTransactionHashRequest)

	transactionVo := b.blockchainBrowserApplicationService.QueryTransactionByTransactionHash(request.TransactionHash)
	if transactionVo == nil {
		FailedHttpResponse(rw, "根据交易哈希未能查询到交易。")
		return
	}

	var response vo.QueryTransactionByTransactionHashResponse
	response.Transaction = transactionVo

	SuccessHttpResponse(rw, "", response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionsByBlockHashTransactionHeight(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryTransactionsByBlockHashTransactionHeightRequest{}).(*vo.QueryTransactionsByBlockHashTransactionHeightRequest)

	pageCondition := request.PageCondition
	if StringUtil.IsNullOrEmpty(request.BlockHash) {
		FailedHttpResponse(rw, "区块哈希不能是空。")
		return
	}
	transactionVos := b.blockchainBrowserApplicationService.QueryTransactionListByBlockHashTransactionHeight(request.BlockHash, pageCondition.From, pageCondition.Size)
	var response vo.QueryTransactionsByBlockHashTransactionHeightResponse
	response.Transactions = transactionVos

	SuccessHttpResponse(rw, "", response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionOutputByAddress(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryTransactionOutputByAddressRequest{}).(*vo.QueryTransactionOutputByAddressRequest)

	transactionOutputDetailVo := b.blockchainBrowserApplicationService.QueryTransactionOutputByAddress(request.Address)
	var response vo.QueryTransactionOutputByAddressResponse
	response.TransactionOutputDetail = transactionOutputDetailVo

	SuccessHttpResponse(rw, "", response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionOutputByTransactionOutputId(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryTransactionOutputByTransactionOutputIdRequest{}).(*vo.QueryTransactionOutputByTransactionOutputIdRequest)

	transactionOutputDetailVo := b.blockchainBrowserApplicationService.QueryTransactionOutputByTransactionOutputId(request.TransactionHash, request.TransactionOutputIndex)
	var response vo.QueryTransactionOutputByTransactionOutputIdResponse
	response.TransactionOutputDetail = transactionOutputDetailVo

	SuccessHttpResponse(rw, "", response)
}
func (b *BlockchainBrowserApplicationController) QueryBlockchainHeight(rw http.ResponseWriter, req *http.Request) {

	blockchainHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	var response vo.QueryBlockchainHeightResponse
	response.BlockchainHeight = blockchainHeight

	SuccessHttpResponse(rw, "", response)
}

func (b *BlockchainBrowserApplicationController) QueryUnconfirmedTransactionByTransactionHash(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryUnconfirmedTransactionByTransactionHashRequest{}).(*vo.QueryUnconfirmedTransactionByTransactionHashRequest)

	unconfirmedTransactionVo := b.blockchainBrowserApplicationService.QueryUnconfirmedTransactionByTransactionHash(request.TransactionHash)
	if unconfirmedTransactionVo == nil {
		FailedHttpResponse(rw, "交易哈希["+request.TransactionHash+"]不是未确认交易。")
		return
	}
	var response vo.QueryUnconfirmedTransactionByTransactionHashResponse
	response.Transaction = unconfirmedTransactionVo

	SuccessHttpResponse(rw, "", response)
}

func (b *BlockchainBrowserApplicationController) QueryUnconfirmedTransactions(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryUnconfirmedTransactionsRequest{}).(*vo.QueryUnconfirmedTransactionsRequest)

	pageCondition := request.PageCondition
	transactionDtos := b.blockchainNetCore.GetBlockchainCore().QueryUnconfirmedTransactions(pageCondition.From, pageCondition.Size)
	if transactionDtos == nil {
		FailedHttpResponse(rw, "未查询到未确认的交易。")
		return
	}

	var unconfirmedTransactionVos []*vo.UnconfirmedTransactionVo
	for _, transactionDto := range transactionDtos {
		unconfirmedTransactionVo := b.blockchainBrowserApplicationService.QueryUnconfirmedTransactionByTransactionHash(TransactionDtoTool.CalculateTransactionHash(transactionDto))
		if unconfirmedTransactionVo != nil {
			unconfirmedTransactionVos = append(unconfirmedTransactionVos, unconfirmedTransactionVo)
		}
	}
	var response vo.QueryUnconfirmedTransactionsResponse
	response.UnconfirmedTransactions = unconfirmedTransactionVos

	SuccessHttpResponse(rw, "", response)
}

func (b *BlockchainBrowserApplicationController) QueryBlockByBlockHeight(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryBlockByBlockHeightRequest{}).(*vo.QueryBlockByBlockHeightRequest)

	blockVo := b.blockchainBrowserApplicationService.QueryBlockViewByBlockHeight(request.BlockHeight)
	if blockVo == nil {
		FailedHttpResponse(rw, "区块链中不存在区块高度["+StringUtil.ValueOfUint64(request.BlockHeight)+"]，请检查输入高度。")
		return
	}
	var response vo.QueryBlockByBlockHeightResponse
	response.Block = blockVo

	SuccessHttpResponse(rw, "", response)
}

func (b *BlockchainBrowserApplicationController) QueryBlockByBlockHash(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.QueryBlockByBlockHashRequest{}).(*vo.QueryBlockByBlockHashRequest)

	block1 := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHash(request.BlockHash)
	if block1 == nil {
		FailedHttpResponse(rw, "区块链中不存在区块哈希["+request.BlockHash+"]，请检查输入哈希。")
		return
	}
	blockVo := b.blockchainBrowserApplicationService.QueryBlockViewByBlockHeight(block1.Height)
	var response vo.QueryBlockByBlockHashResponse
	response.Block = blockVo

	SuccessHttpResponse(rw, "", response)
}

func (b *BlockchainBrowserApplicationController) QueryTop10Blocks(rw http.ResponseWriter, req *http.Request) {
	var blocks []*model.Block
	blockHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	for {
		if blockHeight <= GenesisBlockSetting.HEIGHT {
			break
		}
		block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(blockHeight)
		blocks = append(blocks, block)
		if len(blocks) >= 10 {
			break
		}
		blockHeight--
	}

	var blockVos []vo.BlockVo2
	for _, block := range blocks {
		var blockVo vo.BlockVo2
		blockVo.Height = block.Height
		blockVo.BlockSize = "100字符" //TODO SizeTool.CalculateBlockSize(block1) + "字符" TODO
		blockVo.TransactionCount = BlockTool.GetTransactionCount(block)
		blockVo.MinerIncentiveValue = BlockTool.GetWritedIncentiveValue(block)
		blockVo.Time = TimeUtil.FormatMillisecondTimestamp(block.Timestamp)
		blockVo.Hash = block.Hash
		blockVos = append(blockVos, blockVo)
	}

	var response vo.QueryTop10BlocksResponse
	response.Blocks = blockVos

	SuccessHttpResponse(rw, "", response)
}

//TODO
func SuccessHttpResponse(rw http.ResponseWriter, message string, response interface{}) {
	s := "{\"status\":\"success\",\"message\":\"" + message + "\",\"data\":" + JsonUtil.ToString(response) + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func FailedHttpResponse(rw http.ResponseWriter, message string) {
	s := "{\"status\":\"failed\",\"message\":\"" + message + "\",\"data\":null" + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
