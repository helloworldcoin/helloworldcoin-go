package StructureTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/setting/BlockSetting"
	"helloworldcoin-go/util/LogUtil"
)

/**
 * Check Block Structure
 */
func CheckBlockStructure(block *model.Block) bool {
	transactions := block.Transactions
	if transactions == nil || len(transactions) == 0 {
		LogUtil.Debug("Block data error: The number of transactions in the block is 0. A block must have a coinbase transaction.")
		return false
	}
	//Check the number of transactions in the block
	transactionCount := BlockTool.GetTransactionCount(block)
	if transactionCount > BlockSetting.BLOCK_MAX_TRANSACTION_COUNT {
		LogUtil.Debug("Block data error: The number of transactions in the block exceeds the limit.")
		return false
	}
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		if i == 0 {
			if transaction.TransactionType != TransactionType.COINBASE_TRANSACTION {
				LogUtil.Debug("Block data error: The first transaction of the block must be a coinbase transaction.")
				return false
			}
		} else {
			if transaction.TransactionType != TransactionType.STANDARD_TRANSACTION {
				LogUtil.Debug("Block data error: The non-first transaction of the block must be a standard transaction.")
				return false
			}
		}
	}
	//Check the structure of the transaction
	for _, transaction := range transactions {
		if !CheckTransactionStructure(transaction) {
			LogUtil.Debug("Transaction data error: The transaction structure is abnormal.")
			return false
		}
	}
	return true
}

/**
 * Check Transaction Structure
 */
func CheckTransactionStructure(transaction *model.Transaction) bool {
	if transaction.TransactionType == TransactionType.COINBASE_TRANSACTION {
		inputs := transaction.Inputs
		if inputs != nil && len(inputs) != 0 {
			LogUtil.Debug("Transaction data error: The coinbase transaction cannot have transaction input.")
			return false
		}
		outputs := transaction.Outputs
		if outputs == nil || len(outputs) != 1 {
			LogUtil.Debug("Transaction data error: The coinbase transaction has one and only one transaction output.")
			return false
		}
	} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		inputs := transaction.Inputs
		if inputs == nil || len(inputs) < 1 {
			LogUtil.Debug("Transaction data error: The number of transaction inputs for a standard transaction is at least 1.")
			return false
		}
		outputs := transaction.Outputs
		if outputs == nil || len(outputs) < 1 {
			LogUtil.Debug("Transaction data error: The transaction output number of a standard transaction is at least 1.")
			return false
		}
	} else {
		panic(nil)
	}
	//Check Transaction Input Script
	inputs := transaction.Inputs
	if inputs != nil {
		for _, input := range inputs {
			//Strict Check: must be a P2PKH input script.
			if !ScriptTool.IsPayToPublicKeyHashInputScript(input.InputScript) {
				LogUtil.Debug("Transaction data error: The transaction input script is not a P2PKH input script.")
				return false
			}
		}
	}
	//Check Transaction Output Script
	outputs := transaction.Outputs
	if outputs != nil {
		for _, transactionOutput := range outputs {
			//Strict Check: must be a P2PKH output script.
			if !ScriptTool.IsPayToPublicKeyHashOutputScript(transactionOutput.OutputScript) {
				LogUtil.Debug("Transaction data error: The transaction output script is not a P2PKH output script.")
				return false
			}
		}
	}
	return true
}
