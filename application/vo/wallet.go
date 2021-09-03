package vo

import "helloworld-blockchain-go/dto"

type Payer struct {
	/**
	 * 付款方私钥
	 */
	PrivateKey string `json:"privateKey"`

	/**
	 * 付款来源的交易哈希
	 */
	TransactionHash string `json:"transactionHash"`
	/**
	 * 付款来源的交易输出序列号。
	 */
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
	/**
	 * 付款来源的交易输出的金额
	 */
	Value uint64 `json:"value"`

	/**
	 * 付款方地址
	 */
	Address string `json:"address"`
}
type Payee struct {
	//交易输出的地址
	Address string `json:"address"`

	//交易输出的金额
	Value uint64 `json:"value"`
}
type AutoBuildTransactionRequest struct {
	NonChangePayees []Payee `json:"nonChangePayees"`
}
type AutoBuildTransactionResponse struct {
	//是否构建交易成功
	BuildTransactionSuccess bool `json:"buildTransactionSuccess"`
	//若失败，填写构建失败的原因
	Message string `json:"message"`

	//构建的交易哈希
	TransactionHash string `json:"transactionHash"`
	//交易手续费
	Fee uint64 `json:"fee"`
	//付款方
	Payers []*Payer `json:"payers"`
	//[非找零]收款方
	NonChangePayees []*Payee `json:"nonChangePayees"`
	//[找零]收款方
	ChangePayee Payee `json:"changePayee"`
	//收款方=[非找零]收款方+[找零]收款方
	Payees []Payee `json:"payees"`
	//构建的完整交易
	Transaction dto.TransactionDto `json:"transaction"`
}
