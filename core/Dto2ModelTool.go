package core

import (
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/Model/TransactionType"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/dto"
)

func BlockDto2Block(blockchainDataBase *BlockchainDatabase, blockDto *dto.BlockDto) *Model.Block {
	previousBlockHash := blockDto.PreviousHash
	previousBlock := blockchainDataBase.QueryBlockByBlockHash(previousBlockHash)
	block := new(Model.Block)
	block.Timestamp = blockDto.Timestamp
	block.PreviousBlockHash = previousBlockHash
	block.Nonce = blockDto.Nonce

	blockHeight := BlockTool.GetNextBlockHeight(previousBlock)
	block.Height = blockHeight
	transactionList := transactionDtos2Transactions(blockchainDataBase, blockDto.Transactions)
	block.Transactions = transactionList

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(block)
	block.MerkleTreeRoot = merkleTreeRoot

	blockHash := BlockTool.CalculateBlockHash(block)
	block.Hash = blockHash

	fillBlockProperty(blockchainDataBase, block)

	if !blockchainDataBase.Consensus.CheckConsensus(blockchainDataBase, block) {
		//throw new RuntimeException("区块预检失败。")
		return nil
	}
	return block
}
func transactionDtos2Transactions(blockchainDataBase *BlockchainDatabase, transactionDtoList []dto.TransactionDto) []Model.Transaction {
	var transactions []Model.Transaction
	if transactionDtoList != nil {
		for _, transactionDto := range transactionDtoList {
			transaction := transactionDto2Transaction(blockchainDataBase, &transactionDto)
			transactions = append(transactions, *transaction)
		}
	}
	return transactions
}
func transactionDto2Transaction(blockchainDataBase *BlockchainDatabase, transactionDto *dto.TransactionDto) *Model.Transaction {
	var inputs []Model.TransactionInput
	transactionInputDtos := transactionDto.Inputs
	if transactionInputDtos != nil {
		for _, transactionInputDto := range transactionInputDtos {
			unspentTransactionOutput := blockchainDataBase.QueryUnspentTransactionOutputByTransactionOutputId(transactionInputDto.TransactionHash, transactionInputDto.TransactionOutputIndex)
			if unspentTransactionOutput == nil {
				//throw new RuntimeException("非法交易。交易输入并不是一笔未花费交易输出。");
				return nil
			}
			var transactionInput Model.TransactionInput
			transactionInput.UnspentTransactionOutput = *unspentTransactionOutput
			transactionInput.InputScript = transactionInputDto.InputScript
			inputs = append(inputs, transactionInput)
		}
	}

	var outputs []Model.TransactionOutput
	transactionOutputDtos := transactionDto.Outputs
	if transactionOutputDtos != nil {
		for _, transactionOutputDto := range transactionOutputDtos {
			transactionOutput := transactionOutputDto2TransactionOutput(transactionOutputDto)
			outputs = append(outputs, *transactionOutput)
		}
	}

	transaction := new(Model.Transaction)
	transactionType := obtainTransactionDto(transactionDto)
	transaction.TransactionType = transactionType
	transaction.TransactionHash = TransactionDtoTool.CalculateTransactionHash(transactionDto)
	transaction.Inputs = inputs
	transaction.Outputs = outputs
	return transaction
}

func transactionOutputDto2TransactionOutput(transactionOutputDto dto.TransactionOutputDto) *Model.TransactionOutput {
	var transactionOutput Model.TransactionOutput
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
func fillBlockProperty(blockchainDataBase *BlockchainDatabase, block *Model.Block) {
	//TODO
}
