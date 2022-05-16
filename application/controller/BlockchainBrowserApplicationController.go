package controller

/*
 @author x.king xdotking@gmail.com
*/
import (
	"helloworldcoin-go/application/service"
	"helloworldcoin-go/application/vo"
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/SizeTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/netcore"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/TimeUtil"
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
	request := GetRequest(req, vo.QueryTransactionByTransactionHashRequest{}).(*vo.QueryTransactionByTransactionHashRequest)

	transactionVo := b.blockchainBrowserApplicationService.QueryTransactionByTransactionHash(request.TransactionHash)
	var response vo.QueryTransactionByTransactionHashResponse
	response.Transaction = transactionVo

	success(rw, response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionsByBlockHashTransactionHeight(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryTransactionsByBlockHashTransactionHeightRequest{}).(*vo.QueryTransactionsByBlockHashTransactionHeightRequest)

	pageCondition := request.PageCondition
	transactionVos := b.blockchainBrowserApplicationService.QueryTransactionListByBlockHashTransactionHeight(request.BlockHash, pageCondition.From, pageCondition.Size)
	var response vo.QueryTransactionsByBlockHashTransactionHeightResponse
	response.Transactions = transactionVos

	success(rw, response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionOutputByAddress(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryTransactionOutputByAddressRequest{}).(*vo.QueryTransactionOutputByAddressRequest)

	transactionOutputVo3 := b.blockchainBrowserApplicationService.QueryTransactionOutputByAddress(request.Address)
	var response vo.QueryTransactionOutputByAddressResponse
	response.TransactionOutput = transactionOutputVo3

	success(rw, response)
}
func (b *BlockchainBrowserApplicationController) QueryTransactionOutputByTransactionOutputId(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryTransactionOutputByTransactionOutputIdRequest{}).(*vo.QueryTransactionOutputByTransactionOutputIdRequest)

	transactionOutputVo3 := b.blockchainBrowserApplicationService.QueryTransactionOutputByTransactionOutputId(request.TransactionHash, request.TransactionOutputIndex)
	var response vo.QueryTransactionOutputByTransactionOutputIdResponse
	response.TransactionOutput = transactionOutputVo3

	success(rw, response)
}
func (b *BlockchainBrowserApplicationController) QueryBlockchainHeight(rw http.ResponseWriter, req *http.Request) {

	blockchainHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	var response vo.QueryBlockchainHeightResponse
	response.BlockchainHeight = blockchainHeight

	success(rw, response)
}

func (b *BlockchainBrowserApplicationController) QueryUnconfirmedTransactionByTransactionHash(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryUnconfirmedTransactionByTransactionHashRequest{}).(*vo.QueryUnconfirmedTransactionByTransactionHashRequest)

	unconfirmedTransactionVo := b.blockchainBrowserApplicationService.QueryUnconfirmedTransactionByTransactionHash(request.TransactionHash)
	var response vo.QueryUnconfirmedTransactionByTransactionHashResponse
	response.Transaction = unconfirmedTransactionVo

	success(rw, response)
}

func (b *BlockchainBrowserApplicationController) QueryUnconfirmedTransactions(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryUnconfirmedTransactionsRequest{}).(*vo.QueryUnconfirmedTransactionsRequest)

	pageCondition := request.PageCondition
	transactionDtos := b.blockchainNetCore.GetBlockchainCore().QueryUnconfirmedTransactions(pageCondition.From, pageCondition.Size)
	if transactionDtos == nil {
		var response vo.QueryUnconfirmedTransactionsResponse
		success(rw, response)
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

	success(rw, response)
}

func (b *BlockchainBrowserApplicationController) QueryBlockByBlockHeight(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryBlockByBlockHeightRequest{}).(*vo.QueryBlockByBlockHeightRequest)

	blockVo := b.blockchainBrowserApplicationService.QueryBlockViewByBlockHeight(request.BlockHeight)
	var response vo.QueryBlockByBlockHeightResponse
	response.Block = blockVo

	success(rw, response)
}
func (b *BlockchainBrowserApplicationController) QueryBlockByBlockHash(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.QueryBlockByBlockHashRequest{}).(*vo.QueryBlockByBlockHashRequest)

	block1 := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHash(request.BlockHash)
	if block1 == nil {
		var response vo.QueryBlockByBlockHashResponse
		success(rw, response)
		return
	}
	blockVo := b.blockchainBrowserApplicationService.QueryBlockViewByBlockHeight(block1.Height)
	var response vo.QueryBlockByBlockHashResponse
	response.Block = blockVo

	success(rw, response)
}

func (b *BlockchainBrowserApplicationController) QueryLatest10Blocks(rw http.ResponseWriter, req *http.Request) {
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
		blockVo.BlockSize = SizeTool.CalculateBlockSize(block)
		blockVo.TransactionCount = BlockTool.GetTransactionCount(block)
		blockVo.MinerIncentiveValue = BlockTool.GetWritedIncentiveValue(block)
		blockVo.Time = TimeUtil.FormatMillisecondTimestamp(block.Timestamp)
		blockVo.Hash = block.Hash
		blockVos = append(blockVos, blockVo)
	}

	var response vo.QueryLatest10BlocksResponse
	response.Blocks = blockVos

	success(rw, response)
}
