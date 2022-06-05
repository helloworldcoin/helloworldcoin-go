package service

/*
 @author x.king xdotking@gmail.com
*/
import (
	"helloworldcoin-go/application/vo"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/SizeTool"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/netcore"
	"helloworldcoin-go/util/StringUtil"
	"helloworldcoin-go/util/TimeUtil"
)

type BlockchainBrowserApplicationService struct {
	blockchainNetCore *netcore.BlockchainNetCore
}

func NewBlockchainBrowserApplicationService(blockchainNetCore *netcore.BlockchainNetCore) *BlockchainBrowserApplicationService {
	var b BlockchainBrowserApplicationService
	b.blockchainNetCore = blockchainNetCore
	return &b
}
func (b *BlockchainBrowserApplicationService) QueryTransactionOutputByTransactionOutputId(transactionHash string, transactionOutputIndex uint64) *vo.TransactionOutputVo3 {
	transactionOutput := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryTransactionOutputByTransactionOutputId(transactionHash, transactionOutputIndex)
	if transactionOutput == nil {
		return nil
	}

	var transactionOutputVo3 vo.TransactionOutputVo3
	transactionOutputVo3.FromBlockHeight = transactionOutput.BlockHeight
	transactionOutputVo3.FromBlockHash = transactionOutput.BlockHash
	transactionOutputVo3.FromTransactionHash = transactionOutput.TransactionHash
	transactionOutputVo3.Value = transactionOutput.Value
	transactionOutputVo3.FromOutputScript = ScriptTool.OutputScript2String(transactionOutput.OutputScript)
	transactionOutputVo3.FromTransactionOutputIndex = transactionOutput.TransactionOutputIndex

	transactionOutputTemp := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryUnspentTransactionOutputByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
	transactionOutputVo3.UnspentTransactionOutput = transactionOutputTemp != nil

	inputTransactionVo := b.QueryTransactionByTransactionHash(transactionOutput.TransactionHash)
	transactionOutputVo3.InputTransaction = inputTransactionVo
	transactionOutputVo3.TransactionType = inputTransactionVo.TransactionType

	var outputTransactionVo *vo.TransactionVo
	if transactionOutputTemp == nil {
		destinationTransaction := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().QueryDestinationTransactionByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
		inputs := destinationTransaction.Inputs
		if inputs != nil {
			for inputIndex := 0; inputIndex < len(inputs); inputIndex++ {
				transactionInput := inputs[inputIndex]
				unspentTransactionOutput := transactionInput.UnspentTransactionOutput
				if StringUtil.Equals(transactionOutput.TransactionHash, unspentTransactionOutput.TransactionHash) &&
					transactionOutput.TransactionOutputIndex == unspentTransactionOutput.TransactionOutputIndex {
					transactionOutputVo3.ToTransactionInputIndex = uint64(inputIndex) + uint64(1)
					transactionOutputVo3.ToInputScript = ScriptTool.InputScript2String(transactionInput.InputScript)
					break
				}
			}
		}
		outputTransactionVo = b.QueryTransactionByTransactionHash(destinationTransaction.TransactionHash)
		transactionOutputVo3.ToBlockHeight = outputTransactionVo.BlockHeight
		transactionOutputVo3.ToBlockHash = outputTransactionVo.BlockHash
		transactionOutputVo3.ToTransactionHash = outputTransactionVo.TransactionHash
		transactionOutputVo3.OutputTransaction = outputTransactionVo
	}
	return &transactionOutputVo3

}

func (b *BlockchainBrowserApplicationService) QueryTransactionOutputByAddress(address string) *vo.TransactionOutputVo3 {
	transactionOutput := b.blockchainNetCore.GetBlockchainCore().QueryTransactionOutputByAddress(address)
	if transactionOutput == nil {
		return nil
	}
	transactionOutputVo3 := b.QueryTransactionOutputByTransactionOutputId(transactionOutput.TransactionHash, transactionOutput.TransactionOutputIndex)
	return transactionOutputVo3
}

func (b *BlockchainBrowserApplicationService) QueryTransactionListByBlockHashTransactionHeight(blockHash string, from uint64, size uint64) []*vo.TransactionVo {
	block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHash(blockHash)
	if block == nil {
		return nil
	}
	var transactionVos []*vo.TransactionVo
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

func (b *BlockchainBrowserApplicationService) QueryBlockViewByBlockHeight(blockHeight uint64) *vo.BlockVo {
	block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(blockHeight)
	if block == nil {
		return nil
	}

	blockchainHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	nextBlock := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(block.Height + 1)

	var blockVo vo.BlockVo
	blockVo.Height = block.Height
	blockVo.BlockConfirmations = blockchainHeight - block.Height + 1
	blockVo.BlockSize = SizeTool.CalculateBlockSize(block)
	blockVo.TransactionCount = BlockTool.GetTransactionCount(block)
	blockVo.Time = TimeUtil.FormatMillisecondTimestamp(block.Timestamp)
	blockVo.MinerIncentiveValue = BlockTool.GetWritedIncentiveValue(block)
	blockVo.Difficulty = BlockTool.FormatDifficulty(block.Difficulty)
	blockVo.Nonce = block.Nonce
	blockVo.Hash = block.Hash
	blockVo.PreviousBlockHash = block.PreviousHash
	if nextBlock == nil {
		//blockVo.NextBlockHash=nil
	} else {
		blockVo.NextBlockHash = nextBlock.Hash
	}
	blockVo.MerkleTreeRoot = block.MerkleTreeRoot
	return &blockVo
}

func (b *BlockchainBrowserApplicationService) QueryUnconfirmedTransactionByTransactionHash(transactionHash string) *vo.UnconfirmedTransactionVo {
	transactionDto := b.blockchainNetCore.GetBlockchainCore().QueryUnconfirmedTransactionByTransactionHash(transactionHash)
	if transactionDto == nil {
		return nil
	}
	transaction := b.blockchainNetCore.GetBlockchainCore().GetBlockchainDatabase().TransactionDto2Transaction(transactionDto)
	var transactionDtoVo vo.UnconfirmedTransactionVo
	transactionDtoVo.TransactionHash = transaction.TransactionHash

	var inputDtos []*vo.TransactionInputVo2
	inputs := transaction.Inputs
	if inputs != nil {
		for _, input := range inputs {
			var transactionInputVo vo.TransactionInputVo2
			transactionInputVo.Address = input.UnspentTransactionOutput.Address
			transactionInputVo.TransactionHash = input.UnspentTransactionOutput.TransactionHash
			transactionInputVo.TransactionOutputIndex = input.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputVo.Value = input.UnspentTransactionOutput.Value
			inputDtos = append(inputDtos, &transactionInputVo)
		}
	}
	transactionDtoVo.TransactionInputs = inputDtos

	var outputDtos []*vo.TransactionOutputVo2
	outputs := transaction.Outputs
	if outputs != nil {
		for _, output := range outputs {
			var transactionOutputVo vo.TransactionOutputVo2
			transactionOutputVo.Address = output.Address
			transactionOutputVo.Value = output.Value
			outputDtos = append(outputDtos, &transactionOutputVo)
		}
	}
	transactionDtoVo.TransactionOutputs = outputDtos

	return &transactionDtoVo
}

func (b *BlockchainBrowserApplicationService) QueryTransactionByTransactionHash(transactionHash string) *vo.TransactionVo {
	transaction := b.blockchainNetCore.GetBlockchainCore().QueryTransactionByTransactionHash(transactionHash)
	if transaction == nil {
		return nil
	}

	var transactionVo vo.TransactionVo
	transactionVo.TransactionHash = transaction.TransactionHash
	transactionVo.BlockHeight = transaction.BlockHeight

	transactionVo.TransactionFee = TransactionTool.GetTransactionFee(transaction)
	transactionVo.TransactionType = transaction.TransactionType
	transactionVo.TransactionInputCount = TransactionTool.GetTransactionInputCount(transaction)
	transactionVo.TransactionOutputCount = TransactionTool.GetTransactionOutputCount(transaction)
	transactionVo.TransactionInputValues = TransactionTool.GetInputValue(transaction)
	transactionVo.TransactionOutputValues = TransactionTool.GetOutputValue(transaction)

	blockchainHeight := b.blockchainNetCore.GetBlockchainCore().QueryBlockchainHeight()
	block := b.blockchainNetCore.GetBlockchainCore().QueryBlockByBlockHeight(transaction.BlockHeight)
	transactionVo.BlockConfirmations = blockchainHeight - block.Height + 1
	transactionVo.BlockTime = TimeUtil.FormatMillisecondTimestamp(block.Timestamp)
	transactionVo.BlockHash = block.Hash

	inputs := transaction.Inputs
	var transactionInputVos []*vo.TransactionInputVo
	if inputs != nil {
		for _, transactionInput := range inputs {
			var transactionInputVo vo.TransactionInputVo
			transactionInputVo.Address = transactionInput.UnspentTransactionOutput.Address
			transactionInputVo.Value = transactionInput.UnspentTransactionOutput.Value
			transactionInputVo.InputScript = ScriptTool.InputScript2String(transactionInput.InputScript)
			transactionInputVo.TransactionHash = transactionInput.UnspentTransactionOutput.TransactionHash
			transactionInputVo.TransactionOutputIndex = transactionInput.UnspentTransactionOutput.TransactionOutputIndex
			transactionInputVos = append(transactionInputVos, &transactionInputVo)
		}
	}
	transactionVo.TransactionInputs = transactionInputVos

	outputs := transaction.Outputs
	var transactionOutputVos []*vo.TransactionOutputVo
	if outputs != nil {
		for _, transactionOutput := range outputs {
			var transactionOutputVo vo.TransactionOutputVo
			transactionOutputVo.Address = transactionOutput.Address
			transactionOutputVo.Value = transactionOutput.Value
			transactionOutputVo.OutputScript = ScriptTool.OutputScript2String(transactionOutput.OutputScript)
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
