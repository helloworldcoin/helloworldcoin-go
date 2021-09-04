package TransactionSettingTool

/*
 @author king 409060350@qq.com
*/

/**
 * 校验交易金额是否是一个合法的交易金额：这里用于限制交易金额的最大值、最小值、小数保留位等
 */
func CheckTransactionValue(transactionAmount uint64) bool {
	//交易金额不能小于等于0
	if transactionAmount <= 0 {
		return false
	}
	//最大值是2^64
	//小数保留位是0位
	return true
}
