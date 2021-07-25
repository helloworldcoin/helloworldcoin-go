package StructureTool

import (
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/core/Model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/setting/BlockSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

func CheckBlockStructure(block *Model.Block) bool {
	transactions := block.Transactions
	if transactions == nil || len(transactions) == 0 {
		LogUtil.Debug("区块数据异常：区块中的交易数量为0。区块必须有一笔创世的交易。")
		return false
	}
	//校验区块中交易的数量
	transactionCount := BlockTool.GetTransactionCount(block)
	if transactionCount > BlockSetting.BLOCK_MAX_TRANSACTION_COUNT {
		LogUtil.Debug("区块包含的交易数量是[" + StringUtil.ValueOfUint64(transactionCount) + "]，超过了限制[" + StringUtil.ValueOfUint64(BlockSetting.BLOCK_MAX_TRANSACTION_COUNT) + "]。")
		return false
	}
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]
		if i == 0 {
			if transaction.TransactionType != TransactionType.GENESIS_TRANSACTION {
				LogUtil.Debug("区块数据异常：区块第一笔交易必须是创世交易。")
				return false
			}
		} else {
			if transaction.TransactionType != TransactionType.STANDARD_TRANSACTION {
				LogUtil.Debug("区块数据异常：区块非第一笔交易必须是标准交易。")
				return false
			}
		}
	}
	//校验交易的结构
	for _, transaction := range transactions {
		if !CheckTransactionStructure(&transaction) {
			LogUtil.Debug("交易数据异常：交易结构异常。")
			return false
		}
	}
	return true
}

func CheckTransactionStructure(transaction *Model.Transaction) bool {
	if transaction.TransactionType == TransactionType.GENESIS_TRANSACTION {
		inputs := transaction.Inputs
		if inputs != nil && len(inputs) != 0 {
			LogUtil.Debug("交易数据异常：创世交易不能有交易输入。")
			return false
		}
		outputs := transaction.Outputs
		if outputs == nil || len(outputs) != 1 {
			LogUtil.Debug("交易数据异常：创世交易有且只能有一笔交易输出。")
			return false
		}
	} else if transaction.TransactionType == TransactionType.STANDARD_TRANSACTION {
		inputs := transaction.Inputs
		if inputs == nil || len(inputs) < 1 {
			LogUtil.Debug("交易数据异常：标准交易的交易输入数量至少是1。")
			return false
		}
		outputs := transaction.Outputs
		if outputs == nil || len(outputs) < 1 {
			LogUtil.Debug("交易数据异常：标准交易的交易输出数量至少是1。")
			return false
		}
	} else {
		LogUtil.Debug("交易数据异常：不能识别的交易的类型。")
		return false
	}
	//校验脚本结构
	//输入脚本不需要校验，如果输入脚本结构有误，则在业务[交易输入脚本解锁交易输出脚本]上就通不过。

	//校验输入脚本
	inputs := transaction.Inputs
	if inputs != nil {
		for _, input := range inputs {
			//这里采用严格校验，必须是P2PKH输入脚本。
			if !ScriptTool.IsPayToPublicKeyHashInputScript(input.InputScript) {
				LogUtil.Debug("交易数据异常：创世交易不能有交易输入。")
				return false
			}
		}
	}
	//校验输出脚本
	//校验输出脚本
	outputs := transaction.Outputs
	if outputs != nil {
		for _, transactionOutput := range outputs {
			if !ScriptTool.IsPayToPublicKeyHashOutputScript(transactionOutput.OutputScript) {
				return false
			}
		}
	}
	return true
}
