package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core"
	"helloworldcoin-go/core/tool/BlockDtoTool"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/netcore-client/client"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/model"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/StringUtil"
	"helloworldcoin-go/util/ThreadUtil"
)

type BlockSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	slaveBlockchainCore  *core.BlockchainCore
	nodeService          *service.NodeService
}

func NewBlockSearcher(netCoreConfiguration *configuration.NetCoreConfiguration, blockchainCore *core.BlockchainCore, slaveBlockchainCore *core.BlockchainCore, nodeService *service.NodeService) *BlockSearcher {
	var blockSearcher BlockSearcher
	blockSearcher.netCoreConfiguration = netCoreConfiguration
	blockSearcher.blockchainCore = blockchainCore
	blockSearcher.slaveBlockchainCore = slaveBlockchainCore
	blockSearcher.nodeService = nodeService
	return &blockSearcher
}

func (b *BlockSearcher) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("'search for blocks in the blockchain network' error.", e)
		}
	}()
	for {
		b.searchNodesBlocks()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetBlockSearchTimeInterval())
	}
}

/**
 * Search for new blocks and sync these blocks to the local blockchain system
 */
func (b *BlockSearcher) searchNodesBlocks() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		b.searchNodeBlocks(b.blockchainCore, b.slaveBlockchainCore, node)
	}
}

/**
 * search blocks
 */
func (b *BlockSearcher) searchNodeBlocks(masterBlockchainCore *core.BlockchainCore, slaveBlockchainCore *core.BlockchainCore, node *model.Node) {
	if !b.netCoreConfiguration.IsAutoSearchBlock() {
		return
	}
	masterBlockchainHeight := masterBlockchainCore.QueryBlockchainHeight()
	if masterBlockchainHeight >= node.BlockchainHeight {
		return
	}
	fork := b.isForkNode(masterBlockchainCore, node)
	if fork {
		isHardFork := b.isHardForkNode(masterBlockchainCore, node)
		if !isHardFork {
			forkBlockHeight := b.getForkBlockHeight(masterBlockchainCore, node)
			b.duplicateBlockchainCore(masterBlockchainCore, slaveBlockchainCore)
			slaveBlockchainCore.DeleteBlocks(forkBlockHeight)
			b.synchronizeBlocks(slaveBlockchainCore, node, forkBlockHeight)
			b.promoteBlockchainCore(slaveBlockchainCore, masterBlockchainCore)
		}
	} else {
		nextBlockHeight := masterBlockchainCore.QueryBlockchainHeight() + 1
		b.synchronizeBlocks(masterBlockchainCore, node, nextBlockHeight)
	}
}

/**
 * duplicate BlockchainCore
 */
func (b *BlockSearcher) duplicateBlockchainCore(fromBlockchainCore *core.BlockchainCore, toBlockchainCore *core.BlockchainCore) {
	//delete blocks
	for {
		toBlockchainTailBlock := toBlockchainCore.QueryTailBlock()
		if toBlockchainTailBlock == nil {
			break
		}
		fromBlockchainBlock := fromBlockchainCore.QueryBlockByBlockHeight(toBlockchainTailBlock.Height)
		if BlockTool.IsBlockEquals(fromBlockchainBlock, toBlockchainTailBlock) {
			break
		}
		toBlockchainCore.DeleteTailBlock()
	}
	//add blocks
	for {
		toBlockchainHeight := toBlockchainCore.QueryBlockchainHeight()
		nextBlock := fromBlockchainCore.QueryBlockByBlockHeight(toBlockchainHeight + 1)
		if nextBlock == nil {
			break
		}
		toBlockchainCore.AddBlock(nextBlock)
	}
}

/**
 * promote BlockchainCore
 */
func (b *BlockSearcher) promoteBlockchainCore(fromBlockchainCore *core.BlockchainCore, toBlockchainCore *core.BlockchainCore) {
	if toBlockchainCore.QueryBlockchainHeight() >= fromBlockchainCore.QueryBlockchainHeight() {
		return
	}
	//hard fork
	if b.isHardFork(toBlockchainCore, fromBlockchainCore) {
		return
	}
	b.duplicateBlockchainCore(fromBlockchainCore, toBlockchainCore)
}

func (b *BlockSearcher) getForkBlockHeight(blockchainCore *core.BlockchainCore, node *model.Node) uint64 {
	masterBlockchainHeight := blockchainCore.QueryBlockchainHeight()
	forkBlockHeight := masterBlockchainHeight
	for {
		if forkBlockHeight <= GenesisBlockSetting.HEIGHT {
			break
		}
		var getBlockRequest dto.GetBlockRequest
		getBlockRequest.BlockHeight = forkBlockHeight
		nodeClient := client.NewNodeClient(node.Ip)

		getBlockResponse := nodeClient.GetBlock(getBlockRequest)
		if getBlockResponse == nil {
			break
		}
		remoteBlock := getBlockResponse.Block
		if remoteBlock == nil {
			break
		}
		localBlock := blockchainCore.QueryBlockByBlockHeight(forkBlockHeight)
		if BlockDtoTool.IsBlockEquals(Model2DtoTool.Block2BlockDto(localBlock), remoteBlock) {
			break
		}
		forkBlockHeight--
	}
	forkBlockHeight++
	return forkBlockHeight
}

func (b *BlockSearcher) synchronizeBlocks(blockchainCore *core.BlockchainCore, node *model.Node, startBlockHeight uint64) {
	for {
		var getBlockRequest dto.GetBlockRequest
		getBlockRequest.BlockHeight = startBlockHeight
		nodeClient := client.NewNodeClient(node.Ip)
		getBlockResponse := nodeClient.GetBlock(getBlockRequest)
		if getBlockResponse == nil {
			break
		}
		remoteBlock := getBlockResponse.Block
		if remoteBlock == nil {
			break
		}
		isAddBlockSuccess := blockchainCore.AddBlockDto(remoteBlock)
		if !isAddBlockSuccess {
			break
		}
		startBlockHeight++
	}
}

func (b *BlockSearcher) isForkNode(blockchainCore *core.BlockchainCore, node *model.Node) bool {
	block := blockchainCore.QueryTailBlock()
	if block == nil {
		return false
	}
	var getBlockRequest dto.GetBlockRequest
	getBlockRequest.BlockHeight = block.Height
	nodeClient := client.NewNodeClient(node.Ip)
	getBlockResponse := nodeClient.GetBlock(getBlockRequest)
	//no block with this height exist, so no fork.
	if getBlockResponse == nil {
		return false
	}
	blockDto := getBlockResponse.Block
	if blockDto == nil {
		return false
	}
	blockHash := BlockDtoTool.CalculateBlockHash(blockDto)
	return !StringUtil.Equals(block.Hash, blockHash)
}

func (b *BlockSearcher) isHardFork(blockchainCore1 *core.BlockchainCore, blockchainCore2 *core.BlockchainCore) bool {
	var longer *core.BlockchainCore
	var shorter *core.BlockchainCore
	if blockchainCore1.QueryBlockchainHeight() >= blockchainCore2.QueryBlockchainHeight() {
		longer = blockchainCore1
		shorter = blockchainCore2
	} else {
		longer = blockchainCore2
		shorter = blockchainCore1
	}

	shorterBlockchainHeight := shorter.QueryBlockchainHeight()
	if shorterBlockchainHeight < b.netCoreConfiguration.GetHardForkBlockCount() {
		return false
	}

	criticalPointBlocHeight := shorterBlockchainHeight - b.netCoreConfiguration.GetHardForkBlockCount() + 1
	longerBlock := longer.QueryBlockByBlockHeight(criticalPointBlocHeight)
	shorterBlock := shorter.QueryBlockByBlockHeight(criticalPointBlocHeight)
	return !BlockTool.IsBlockEquals(longerBlock, shorterBlock)
}

func (b *BlockSearcher) isHardForkNode(blockchainCore *core.BlockchainCore, node *model.Node) bool {
	blockchainHeight := blockchainCore.QueryBlockchainHeight()
	if blockchainHeight < b.netCoreConfiguration.GetHardForkBlockCount() {
		return false
	}
	criticalPointBlocHeight := blockchainHeight - b.netCoreConfiguration.GetHardForkBlockCount() + 1
	if criticalPointBlocHeight <= GenesisBlockSetting.HEIGHT {
		return false
	}
	var getBlockRequest dto.GetBlockRequest
	getBlockRequest.BlockHeight = criticalPointBlocHeight
	nodeClient := client.NewNodeClient(node.Ip)
	getBlockResponse := nodeClient.GetBlock(getBlockRequest)
	if getBlockResponse == nil {
		return false
	}
	remoteBlock := getBlockResponse.Block
	if remoteBlock == nil {
		return false
	}
	localBlock := blockchainCore.QueryBlockByBlockHeight(criticalPointBlocHeight)
	return !BlockDtoTool.IsBlockEquals(Model2DtoTool.Block2BlockDto(localBlock), remoteBlock)
}
