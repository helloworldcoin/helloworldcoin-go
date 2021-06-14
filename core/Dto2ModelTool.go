package core

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/dto"
)

func BlockDto2Block(blockchainDatabase *BlockchainDatabase, blockDto *dto.BlockDto) *model.Block {
	previousHash := blockDto.PreviousHash
	previousBlock := blockchainDatabase.QueryBlockByBlockHash(previousHash)
	block := new(model.Block)
	block.Timestamp = blockDto.Timestamp
	block.PreviousHash = previousHash
	block.Nonce = blockDto.Nonce

	blockHeight := BlockTool.GetNextBlockHeight(previousBlock)
	block.Height = blockHeight
	transactionList := transactionDtos2Transactions(blockchainDatabase, blockDto.Transactions)
	block.Transactions = transactionList

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(block)
	block.MerkleTreeRoot = merkleTreeRoot

	blockHash := BlockTool.CalculateBlockHash(block)
	block.Hash = blockHash

	difficult := blockchainDatabase.Consensus.CalculateDifficult(blockchainDatabase, block)
	block.Difficulty = difficult

	fillBlockProperty(blockchainDatabase, block)

	if !blockchainDatabase.Consensus.CheckConsensus(blockchainDatabase, block) {
		//throw new RuntimeException("区块预检失败。")
		return nil
	}
	return block
}
func transactionDtos2Transactions(blockchainDatabase *BlockchainDatabase, transactionDtoList []dto.TransactionDto) []model.Transaction {
	var transactions []model.Transaction
	if transactionDtoList != nil {
		for _, transactionDto := range transactionDtoList {
			transaction := transactionDto2Transaction(blockchainDatabase, &transactionDto)
			transactions = append(transactions, *transaction)
		}
	}
	return transactions
}
func transactionDto2Transaction(blockchainDatabase *BlockchainDatabase, transactionDto *dto.TransactionDto) *model.Transaction {
	var inputs []model.TransactionInput
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			unspentTransactionOutput := blockchainDatabase.QueryUnspentTransactionOutputByTransactionOutputId(transactionInputDto.TransactionHash, transactionInputDto.TransactionOutputIndex)
			if unspentTransactionOutput == nil {
				//throw new RuntimeException("非法交易。交易输入并不是一笔未花费交易输出。");
				return nil
			}
			var transactionInput model.TransactionInput
			transactionInput.UnspentTransactionOutput = *unspentTransactionOutput
			transactionInput.InputScript = transactionInputDto.InputScript
			inputs = append(inputs, transactionInput)
		}
	}

	var outputs []model.TransactionOutput
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			transactionOutput := transactionOutputDto2TransactionOutput(transactionOutputDto)
			outputs = append(outputs, *transactionOutput)
		}
	}

	transaction := new(model.Transaction)
	transactionType := obtainTransactionDto(transactionDto)
	transaction.TransactionType = transactionType
	transaction.TransactionHash = TransactionDtoTool.CalculateTransactionHash(transactionDto)
	transaction.Inputs = inputs
	transaction.Outputs = outputs
	return transaction
}

func transactionOutputDto2TransactionOutput(transactionOutputDto dto.TransactionOutputDto) *model.TransactionOutput {
	var transactionOutput model.TransactionOutput
	publicKeyHash := ScriptTool.GetPublicKeyHashByPayToPublicKeyHashOutputScript(transactionOutputDto.OutputScript)
	address := AccountUtil.AddressFromStringPublicKeyHash(publicKeyHash)
	transactionOutput.Address = address
	transactionOutput.Value = transactionOutputDto.Value
	transactionOutput.OutputScript = transactionOutputDto.OutputScript
	return &transactionOutput
}
func obtainTransactionDto(transactionDto *dto.TransactionDto) TransactionType.TransactionType {
	if transactionDto.Inputs == nil || len(transactionDto.Inputs) == 0 {
		return TransactionType.GENESIS_TRANSACTION
	}
	return TransactionType.STANDARD_TRANSACTION
}
func fillBlockProperty(blockchainDatabase *BlockchainDatabase, block *model.Block) {
	transactionIndex := uint64(0)
	transactionHeight := blockchainDatabase.QueryBlockchainTransactionHeight()
	transactionOutputHeight := blockchainDatabase.QueryBlockchainTransactionOutputHeight()
	blockHeight := block.Height
	blockHash := block.Hash
	transactions := block.Transactions
	transactionCount := BlockTool.GetTransactionCount(block)
	block.TransactionCount = transactionCount
	block.PreviousTransactionHeight = transactionHeight
	if transactions != nil {
		for _, transaction := range transactions {
			transactionIndex := transactionIndex + 1
			transactionHeight = transactionHeight + 1
			transaction.BlockHeight = blockHeight
			transaction.TransactionIndex = transactionIndex
			transaction.TransactionHeight = transactionHeight

			outputs := transaction.Outputs
			if outputs != nil {
				for i := 0; i < len(outputs); i = i + 1 {
					transactionOutputHeight := transactionOutputHeight + 1
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
