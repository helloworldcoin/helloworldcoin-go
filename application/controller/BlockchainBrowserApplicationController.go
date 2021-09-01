package controller

import (
	"helloworld-blockchain-go/application/vo/block"
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/TimeUtil"
	"io"
	"net/http"
)

type BlockchainBrowserApplicationController struct {
	blockchainNetCore *netcore.BlockchainNetCore
	blockchainCore    *core.BlockchainCore
	// BlockchainBrowserApplicationService blockchainBrowserApplicationService
}

func NewBlockchainBrowserApplicationController(blockchainNetCore *netcore.BlockchainNetCore, blockchainCore *core.BlockchainCore) *BlockchainBrowserApplicationController {
	var b BlockchainBrowserApplicationController
	b.blockchainNetCore = blockchainNetCore
	b.blockchainCore = blockchainCore
	return &b
}

func (b *BlockchainBrowserApplicationController) QueryTop10Blocks(w http.ResponseWriter, req *http.Request) {

	var blocks []*Model.Block
	blockHeight := b.blockchainCore.QueryBlockchainHeight()
	for {
		if blockHeight <= GenesisBlockSetting.HEIGHT {
			break
		}
		block := b.blockchainCore.QueryBlockByBlockHeight(blockHeight)
		blocks = append(blocks, block)
		if len(blocks) >= 10 {
			break
		}
		blockHeight--
	}

	var blockVos []block.BlockVo2
	for _, block1 := range blocks {
		var blockVo block.BlockVo2
		blockVo.Height = block1.Height
		blockVo.BlockSize = "100字符" //TODO SizeTool.CalculateBlockSize(block1) + "字符" TODO
		blockVo.TransactionCount = BlockTool.GetTransactionCount(block1)
		blockVo.MinerIncentiveValue = BlockTool.GetWritedIncentiveValue(block1)
		blockVo.Time = TimeUtil.FormatMillisecondTimestamp(block1.Timestamp)
		blockVo.Hash = block1.Hash
		blockVos = append(blockVos, blockVo)
	}

	var response block.QueryTop10BlocksResponse
	response.Blocks = blockVos
	r := "{\"status\":\"SUCCESS\",\"message\":\"message\",\"data\":" + JsonUtil.ToString(response) + "}"
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, r)
}
