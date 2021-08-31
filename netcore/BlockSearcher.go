package netcore

import (
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/BlockDtoTool"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/model"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/SystemUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type BlockSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	blockchainCore       *core.BlockchainCore
	slaveBlockchainCore  *core.BlockchainCore
	nodeService          *service.NodeService
}

func (b *BlockSearcher) start() {
	defer func() {
		if err := recover(); err != nil {
			SystemUtil.ErrorExit("在区块链网络中同步节点的区块出现异常", err)
		}
	}()
	for {
		b.searchBlocks()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetSearchBlockTimeInterval())
	}
}

/**
 * 搜索新的区块，并同步这些区块到本地区块链系统
 */
func (b *BlockSearcher) searchBlocks() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}
	for _, node := range nodes {
		b.searchNodeBlocks(b.blockchainCore, b.slaveBlockchainCore, node)
	}
}

/**
 * 搜索远程节点的区块到本地，未分叉同步至主链，硬分叉不同步，软分叉同步至从链
 */
func (b *BlockSearcher) searchNodeBlocks(masterBlockchainCore *core.BlockchainCore, slaveBlockchainCore *core.BlockchainCore, node *model.Node) {
	if !b.netCoreConfiguration.IsAutoSearchBlock() {
		return
	}
	masterBlockchainHeight := masterBlockchainCore.QueryBlockchainHeight()
	//本地区块链高度大于等于远程节点区块链高度，此时远程节点没有可以同步到本地区块链的区块。
	if masterBlockchainHeight >= node.BlockchainHeight {
		return
	}
	//本地区块链与node区块链是否分叉？
	fork := b.isForkNode(masterBlockchainCore, node)
	if fork {
		isHardFork := b.isHardForkNode(masterBlockchainCore, node)
		if !isHardFork {
			//求分叉区块的高度
			forkBlockHeight := b.getForkBlockHeight(masterBlockchainCore, node)
			//复制"主区块链核心"的区块至"从区块链核心"
			b.duplicateBlockchainCore(masterBlockchainCore, slaveBlockchainCore)
			//删除"从区块链核心"已分叉区块
			slaveBlockchainCore.DeleteBlocks(forkBlockHeight)
			//同步远程节点区块至从区块链核心
			b.synchronizeBlocks(slaveBlockchainCore, node, forkBlockHeight)
			//同步从区块链核心的区块至主区块链核心
			b.promoteBlockchainCore(slaveBlockchainCore, masterBlockchainCore)
		}
	} else {
		//未分叉，同步远程节点区块至主区块链核心
		nextBlockHeight := masterBlockchainCore.QueryBlockchainHeight() + 1
		b.synchronizeBlocks(masterBlockchainCore, node, nextBlockHeight)
	}
}

/**
 * 复制区块链核心的区块，操作完成后，'来源区块链核心'区块数据不发生变化，'去向区块链核心'的区块数据与'来源区块链核心'的区块数据保持一致。
 * @param fromBlockchainCore 来源区块链核心
 * @param toBlockchainCore 去向区块链核心
 */
func (b *BlockSearcher) duplicateBlockchainCore(fromBlockchainCore *core.BlockchainCore, toBlockchainCore *core.BlockchainCore) {
	//删除'去向区块链核心'区块
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
	//增加'去向区块链核心'区块
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
 * 增加"去向区块链核心"的区块，操作完成后，"来源区块链核心"的区块不发生变化，"去向区块链核心"的高度不变或者增长。
 * @param fromBlockchainCore "来源区块链核心"
 * @param toBlockchainCore "去向区块链核心"
 */
func (b *BlockSearcher) promoteBlockchainCore(fromBlockchainCore *core.BlockchainCore, toBlockchainCore *core.BlockchainCore) {
	//此时，"去向区块链核心高度"大于"来源区块链核心高度"，"去向区块链核心高度"不能增加，结束逻辑。
	if toBlockchainCore.QueryBlockchainHeight() >= fromBlockchainCore.QueryBlockchainHeight() {
		return
	}
	//硬分叉
	if b.isHardFork(toBlockchainCore, fromBlockchainCore) {
		return
	}
	//此时，"去向区块链核心高度"小于"来源区块链核心高度"，且未硬分叉，可以增加"去向区块链核心高度"
	b.duplicateBlockchainCore(fromBlockchainCore, toBlockchainCore)
}

func (b *BlockSearcher) getForkBlockHeight(blockchainCore *core.BlockchainCore, node *model.Node) uint64 {
	//求分叉区块的高度，此时已知分叉了，从当前高度依次递减1，判断高度相同的区块的是否相等，若相等，(高度+1)即开始分叉高度。
	masterBlockchainHeight := blockchainCore.QueryBlockchainHeight()
	forkBlockHeight := masterBlockchainHeight
	for {
		if forkBlockHeight <= GenesisBlockSetting.HEIGHT {
			break
		}
		var getBlockRequest dto.GetBlockRequest
		getBlockRequest.BlockHeight = forkBlockHeight
		nodeClient := client.NodeClient{Ip: node.Ip}

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
		nodeClient := client.NodeClient{Ip: node.Ip}
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
	nodeClient := client.NodeClient{Ip: node.Ip}
	getBlockResponse := nodeClient.GetBlock(getBlockRequest)
	//没有查询到区块，这里则认为远程节点没有该高度的区块存在，远程节点的高度没有本地区块链高度高，所以不分叉。
	if getBlockResponse == nil {
		return false
	}
	blockDto := getBlockResponse.Block
	if blockDto == nil {
		return false
	}
	blockHash := BlockDtoTool.CalculateBlockHash(blockDto)
	return !StringUtil.IsEquals(block.Hash, blockHash)
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
	nodeClient := client.NodeClient{Ip: node.Ip}
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
