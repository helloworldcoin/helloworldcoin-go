package core

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/ScriptDtoTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/dto"
)

func (b *BlockchainDatabase) BlockDto2Block(blockDto *dto.BlockDto) *model.Block {
	previousHash := blockDto.PreviousHash
	previousBlock := b.QueryBlockByBlockHash(previousHash)
	block := &model.Block{}
	block.Timestamp = blockDto.Timestamp
	block.PreviousHash = previousHash
	block.Nonce = blockDto.Nonce

	blockHeight := BlockTool.GetNextBlockHeight(previousBlock)
	block.Height = blockHeight
	transactions := b.transactionDtos2Transactions(blockDto.Transactions)
	block.Transactions = transactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(block)
	block.MerkleTreeRoot = merkleTreeRoot

	blockHash := BlockTool.CalculateBlockHash(block)
	block.Hash = blockHash

	difficult := b.consensus.CalculateDifficult(b, block)
	block.Difficulty = difficult

	b.fillBlockProperty(block)

	if !b.consensus.CheckConsensus(b, block) {
		//throw new RuntimeException("区块预检失败。")
		return nil
	}
	return block
}
func (b *BlockchainDatabase) transactionDtos2Transactions(transactionDtos []*dto.TransactionDto) []*model.Transaction {
	var transactions []*model.Transaction
	if transactionDtos != nil {
		for _, transactionDto := range transactionDtos {
			transaction := b.TransactionDto2Transaction(transactionDto)
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}
func (b *BlockchainDatabase) TransactionDto2Transaction(transactionDto *dto.TransactionDto) *model.Transaction {
	var inputs []*model.TransactionInput
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			unspentTransactionOutput := b.QueryUnspentTransactionOutputByTransactionOutputId(transactionInputDto.TransactionHash, transactionInputDto.TransactionOutputIndex)
			if unspentTransactionOutput == nil {
				//throw new RuntimeException("非法交易。交易输入并不是一笔未花费交易输出。");
				return nil
			}
			var transactionInput model.TransactionInput
			transactionInput.UnspentTransactionOutput = unspentTransactionOutput
			transactionInput.InputScript = b.InputScriptDto2InputScript(transactionInputDto.InputScript)
			inputs = append(inputs, &transactionInput)
		}
	}
	var outputs []*model.TransactionOutput
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			transactionOutput := b.transactionOutputDto2TransactionOutput(transactionOutputDto)
			outputs = append(outputs, transactionOutput)
		}
	}
	transaction := new(model.Transaction)
	transactionType := b.obtainTransactionDto(transactionDto)
	transaction.TransactionType = transactionType
	transaction.TransactionHash = TransactionDtoTool.CalculateTransactionHash(transactionDto)
	transaction.Inputs = inputs
	transaction.Outputs = outputs
	return transaction
}
func (b *BlockchainDatabase) transactionOutputDto2TransactionOutput(transactionOutputDto *dto.TransactionOutputDto) *model.TransactionOutput {
	var transactionOutput model.TransactionOutput
	publicKeyHash := ScriptDtoTool.GetPublicKeyHashFromPayToPublicKeyHashOutputScript(transactionOutputDto.OutputScript)
	address := AccountUtil.AddressFromPublicKeyHash(publicKeyHash)
	transactionOutput.Address = address
	transactionOutput.Value = transactionOutputDto.Value
	transactionOutput.OutputScript = b.OutputScriptDto2OutputScript(transactionOutputDto.OutputScript)
	return &transactionOutput
}
func (b *BlockchainDatabase) obtainTransactionDto(transactionDto *dto.TransactionDto) TransactionType.TransactionType {
	if transactionDto.Inputs == nil || len(transactionDto.Inputs) == 0 {
		return TransactionType.GENESIS_TRANSACTION
	}
	return TransactionType.STANDARD_TRANSACTION
}
func (b *BlockchainDatabase) fillBlockProperty(block *model.Block) {
	transactionIndex := uint64(0)
	transactionHeight := b.QueryBlockchainTransactionHeight()
	transactionOutputHeight := b.QueryBlockchainTransactionOutputHeight()
	blockHeight := block.Height
	blockHash := block.Hash
	transactions := block.Transactions
	transactionCount := BlockTool.GetTransactionCount(block)
	block.TransactionCount = transactionCount
	block.PreviousTransactionHeight = transactionHeight
	if transactions != nil {
		for _, transaction := range transactions {
			transactionIndex = transactionIndex + 1
			transactionHeight = transactionHeight + 1
			transaction.BlockHeight = blockHeight
			transaction.TransactionIndex = transactionIndex
			transaction.TransactionHeight = transactionHeight

			outputs := transaction.Outputs
			if outputs != nil {
				for i := 0; i < len(outputs); i = i + 1 {
					transactionOutputHeight = transactionOutputHeight + 1
					output := outputs[i]
					output.BlockHeight = blockHeight
					output.BlockHash = blockHash
					output.TransactionHeight = transactionHeight
					output.TransactionHash = transaction.TransactionHash
					output.TransactionOutputIndex = uint64(i) + uint64(1)
					output.TransactionIndex = transaction.TransactionIndex
					output.TransactionOutputHeight = transactionOutputHeight
				}
			}
		}
	}
}
func (b *BlockchainDatabase) OutputScriptDto2OutputScript(outputScriptDto *dto.OutputScriptDto) *model.OutputScript {
	var outputScript model.OutputScript
	outputScript = append(outputScript, *outputScriptDto...)
	return &outputScript
}
func (b *BlockchainDatabase) InputScriptDto2InputScript(inputScriptDto *dto.InputScriptDto) *model.InputScript {
	var inputScript model.InputScript
	inputScript = append(inputScript, *inputScriptDto...)
	return &inputScript
}
