package service

import (
	"helloworld-blockchain-go/application/vo/block"
	"helloworld-blockchain-go/application/vo/transaction"
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/SizeTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/TimeUtil"
)

type BlockchainBrowserApplicationService struct {
	blockchainNetCore *netcore.BlockchainNetCore
}

func NewBlockchainBrowserApplicationService(blockchainNetCore *netcore.BlockchainNetCore) *BlockchainBrowserApplicationService {
	var b BlockchainBrowserApplicationService
	b.blockchainNetCore = blockchainNetCore
	return &b
}
func (b *BlockchainBrowserApplicationService) QueryTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *transaction.TransactionOutputDetailVo {
	//查询交易输出
	transactionOutput := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryTransactionOutputByTransactionOutputId(transactionHash, transactionOutputIndex)
	if transactionOutput == nil {
		return nil
	}

	var transactionOutputDetailVo transaction.TransactionOutputDetailVo
	transactionOutputDetailVo.FromBlockHeight = transactionOutput.BlockHeight
	transactionOutputDetailVo.FromBlockHash = transactionOutput.BlockHash
	transactionOutputDetailVo.FromTransactionHash = transactionOutput.TransactionHash
	transactionOutputDetailVo.Value = transactionOutput.Value
	transactionOutputDetailVo.FromOutputScript = ScriptTool.StringOutputScript(transactionOutput.OutputScript)
	transactionOutputDetailVo.FromTransactionOutputIndex = transactionOutput.TransactionOutputIndex

	//是否是未花费输出
	transactionOutputTemp := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryUnspentTransactionOutputByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
	transactionOutputDetailVo.Spent = transactionOutputTemp == nil

	//来源
	inputTransactionVo := b.QueryTransactionByTransactionHash(transactionOutput.TransactionHash)
	transactionOutputDetailVo.InputTransaction = inputTransactionVo
	transactionOutputDetailVo.TransactionType = inputTransactionVo.TransactionType

	//去向
	var outputTransactionVo *transaction.TransactionVo
	if transactionOutputTemp == nil {
		destinationTransaction := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryDestinationTransactionByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
		inputs := destinationTransaction.Inputs
		if inputs != nil {
			for inputIndex := 0; inputIndex < len(inputs); inputIndex++ {
				transactionInput := inputs[inputIndex]
				unspentTransactionOutput := transactionInput.UnspentTransactionOutput
				if StringUtil.IsEquals(transactionOutput.TransactionHash, unspentTransactionOutput.TransactionHash) &&
					transactionOutput.TransactionOutputIndex == unspentTransactionOutput.TransactionOutputIndex {
					transactionOutputDetailVo.ToTransactionInputIndex = uint64(inputIndex) + uint64(1)
					transactionOutputDetailVo.ToInputScript = ScriptTool.StringInputScript(transactionInput.InputScript)
					break
				}
			}
		}
		outputTransactionVo = b.QueryTransactionByTransactionHash(destinationTransaction.TransactionHash)
		transactionOutputDetailVo.ToBlockHeight = outputTransactionVo.BlockHeight
		transactionOutputDetailVo.ToBlockHash = outputTransactionVo.BlockHash
		transactionOutputDetailVo.ToTransactionHash = outputTransactionVo.TransactionHash
		transactionOutputDetailVo.OutputTransaction = outputTransactionVo
	}
	return &transactionOutputDetailVo

}

func (b *BlockchainBrowserApplicationService) QueryTransactionOutputByAddress(address string) *transaction.TransactionOutputDetailVo {
	transactionOutput := b.blockchainNetCore.GetBlockchainCore().QueryTransactionOutputByAddress(address)
	if transactionOutput == nil {
		return nil
	}
	transactionOutputDetailVo := b.QueryTransactionOutputByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
	return transactionOutputDetailVo
}

func (b *BlockchainBrowserApplicationService) QueryTransactionListByBlockHashTransactionHeight(blockHash string, from uint64, size uint64) []*transaction.TransactionVo {
	block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHash(blockHash)
	var transactionVos []*transaction.TransactionVo
	for i := from; i < from+size; i++ {
		if from < 0 {
			break
		}
		if i > block.TransactionCount {
			break
		}
		transactionHeight := block.PreviousTransactionHeight + i
		transaction := b.blockchainNetCore.GetBlockchainCore().QueryTransactionByTransactionHeight(transactionHeight)
		transactionVo := b.QueryTransactionByTransactionHash(transaction.TransactionHash)
		transactionVos = append(transactionVos, transactionVo)
	}
	return transactionVos
}

func (b *BlockchainBrowserApplicationService) QueryBlockViewByBlockHeight(blockHeight uint64) *block.BlockVo {
	block1 := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(blockHeight)
	if block1 == nil {
		return nil
	}
	nextBlock := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(block1.Height + 1)

	var blockVo block.BlockVo
	blockVo.Height = block1.Height
	blockVo.ConfirmCount = BlockTool.GetTransactionCount(block1)
	blockVo.BlockSize = StringUtil.ValueOfUint64(SizeTool.CalculateBlockSize(block1)) + "字符"
	blockVo.TransactionCount = BlockTool.GetTransactionCount(block1)
	blockVo.Time = TimeUtil.FormatMillisecondTimestamp(block1.Timestamp)
	blockVo.MinerIncentiveValue = BlockTool.GetWritedIncentiveValue(block1)
	blockVo.Difficulty = BlockTool.FormatDifficulty(block1.Difficulty)
	blockVo.Nonce = block1.Nonce
	blockVo.Hash = block1.Hash
	blockVo.PreviousBlockHash = block1.PreviousHash
	if nextBlock == nil {
		//blockVo.NextBlockHash=nil
	} else {
		blockVo.NextBlockHash = nextBlock.Hash
	}
	blockVo.MerkleTreeRoot = block1.MerkleTreeRoot
	return &blockVo
}

func (b *BlockchainBrowserApplicationService) QueryUnconfirmedTransactionByTransactionHash(transactionHash string) *transaction.UnconfirmedTransactionVo {
	transactionDto := b.blockchainNetCore.GetBlockchainCore().QueryUnconfirmedTransactionByTransactionHash(transactionHash)
	if transactionDto == nil {
		return nil
	}
	transaction1 := core.TransactionDto2Transaction(b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase(), transactionDto)
	var transactionDtoVo transaction.UnconfirmedTransactionVo
	transactionDtoVo.TransactionHash = transaction1.TransactionHash

	var inputDtos []*transaction.TransactionInputVo2
	inputs := transaction1.Inputs
	if inputs != nil {
		for _, input := range inputs {
			var transactionInputVo transaction.TransactionInputVo2
			transactionInputVo.Address = input.UnspentTransactionOutput.Address
			transactionInputVo.TransactionHash = input.UnspentTransactionOutput.TransactionHash
			transactionInputVo.TransactionOutputIndex = input.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputVo.Value = input.UnspentTransactionOutput.Value
			inputDtos = append(inputDtos, &transactionInputVo)
		}
	}
	transactionDtoVo.Inputs = inputDtos

	var outputDtos []*transaction.TransactionOutputVo2
	outputs := transaction1.Outputs
	if outputs != nil {
		for _, output := range outputs {
			var transactionOutputVo transaction.TransactionOutputVo2
			transactionOutputVo.Address = output.Address
			transactionOutputVo.Value = output.Value
			outputDtos = append(outputDtos, &transactionOutputVo)
		}
	}
	transactionDtoVo.Outputs = outputDtos

	return &transactionDtoVo
}

func (b *BlockchainBrowserApplicationService) QueryTransactionByTransactionHash(transactionHash string) *transaction.TransactionVo {
	transaction1 := b.blockchainNetCore.GetBlockchainCore().QueryTransactionByTransactionHash(transactionHash)
	if transaction1 == nil {
		return nil
	}

	var transactionVo transaction.TransactionVo
	transactionVo.TransactionHash = transaction1.TransactionHash
	transactionVo.BlockHeight = transaction1.BlockHeight

	transactionVo.TransactionFee = TransactionTool.CalculateTransactionFee(transaction1)
	transactionVo.TransactionType = transaction1.TransactionType
	transactionVo.TransactionInputCount = TransactionTool.GetTransactionInputCount(transaction1)
	transactionVo.TransactionOutputCount = TransactionTool.GetTransactionOutputCount(transaction1)
	transactionVo.TransactionInputValues = TransactionTool.GetInputValue(transaction1)
	transactionVo.TransactionOutputValues = TransactionTool.GetOutputValue(transaction1)

	blockchainHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(transaction1.BlockHeight)
	transactionVo.ConfirmCount = blockchainHeight - block.Height + 1
	transactionVo.BlockTime = TimeUtil.FormatMillisecondTimestamp(block.Timestamp)
	transactionVo.BlockHash = block.Hash

	inputs := transaction1.Inputs
	var transactionInputVos []*transaction.TransactionInputVo
	if inputs != nil {
		for _, transactionInput := range inputs {
			var transactionInputVo transaction.TransactionInputVo
			transactionInputVo.Address = transactionInput.UnspentTransactionOutput.Address
			transactionInputVo.Value = transactionInput.UnspentTransactionOutput.Value
			transactionInputVo.InputScript = ScriptTool.StringInputScript(transactionInput.InputScript)
			transactionInputVo.TransactionHash = transactionInput.UnspentTransactionOutput.TransactionHash
			transactionInputVo.TransactionOutputIndex = transactionInput.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputVos = append(transactionInputVos, &transactionInputVo)
		}
	}
	transactionVo.TransactionInputs = transactionInputVos

	outputs := transaction1.Outputs
	var transactionOutputVos []*transaction.TransactionOutputVo
	if outputs != nil {
		for _, transactionOutput := range outputs {
			var transactionOutputVo transaction.TransactionOutputVo
			transactionOutputVo.Address = transactionOutput.Address
			transactionOutputVo.Value = transactionOutput.Value
			transactionOutputVo.OutputScript = ScriptTool.StringOutputScript(transactionOutput.OutputScript)
			transactionOutputVo.TransactionHash = transactionOutput.TransactionHash
			transactionOutputVo.TransactionOutputIndex = transactionOutput.TransactionOutputIndex
			transactionOutputVos = append(transactionOutputVos, &transactionOutputVo)
		}
	}
	transactionVo.TransactionOutputs = transactionOutputVos

	var inputScripts []string
	for _, transactionInputVo := range transactionInputVos {
		inputScripts = append(inputScripts, transactionInputVo.InputScript)
	}
	transactionVo.InputScripts = inputScripts

	var outputScripts []string
	for _, transactionOutputVo := range transactionOutputVos {
		outputScripts = append(outputScripts, transactionOutputVo.OutputScript)
	}
	transactionVo.OutputScripts = outputScripts

	return &transactionVo
}
