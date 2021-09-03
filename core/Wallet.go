package core

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/tool/EncodeDecodeTool"
	"helloworld-blockchain-go/core/tool/ScriptDtoTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

const WALLET_DATABASE_NAME = "WalletDatabase"

type Wallet struct {
	CoreConfiguration  *CoreConfiguration
	BlockchainDatabase *BlockchainDatabase
}

func (w *Wallet) GetAllAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	bytesAccounts := KvDbUtil.Gets(w.getWalletDatabasePath(), 1, 100000000)
	for e := bytesAccounts.Front(); e != nil; e = e.Next() {
		account := EncodeDecodeTool.DecodeToAccount(e.Value.([]byte))
		accounts = append(accounts, account)
	}
	return accounts
}
func (w *Wallet) GetNonZeroBalanceAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	bytesAccounts := KvDbUtil.Gets(w.getWalletDatabasePath(), 1, 100000000)
	for e := bytesAccounts.Front(); e != nil; e = e.Next() {
		account := EncodeDecodeTool.DecodeToAccount(e.Value.([]byte))
		utxo := w.BlockchainDatabase.QueryUnspentTransactionOutputByAddress(account.Address)
		if utxo != nil && utxo.Value > 0 {
			accounts = append(accounts, account)
		}
	}
	return accounts
}
func (w *Wallet) CreateAccount() *AccountUtil.Account {
	return AccountUtil.RandomAccount()
}
func (w *Wallet) CreateAndSaveAccount() *AccountUtil.Account {
	account := AccountUtil.RandomAccount()
	w.SaveAccount(account)
	return account
}
func (w *Wallet) SaveAccount(account *AccountUtil.Account) {
	KvDbUtil.Put(w.getWalletDatabasePath(), w.getKeyByAccount(account), EncodeDecodeTool.EncodeAccount(account))
}
func (w *Wallet) DeleteAccountByAddress(address string) {
	KvDbUtil.Delete(w.getWalletDatabasePath(), w.getKeyByAddress(address))
}
func (w *Wallet) GetBalanceByAddress(address string) uint64 {
	utxo := w.BlockchainDatabase.QueryUnspentTransactionOutputByAddress(address)
	if utxo != nil {
		return utxo.Value
	}
	return uint64(0)
}
func (w *Wallet) getWalletDatabasePath() string {
	return FileUtil.NewPath(w.CoreConfiguration.getCorePath(), WALLET_DATABASE_NAME)
}
func (w *Wallet) getKeyByAddress(address string) []byte {
	return ByteUtil.StringToUtf8Bytes(address)
}
func (w *Wallet) getKeyByAccount(account *AccountUtil.Account) []byte {
	return ByteUtil.StringToUtf8Bytes(account.Address)
}

func (w *Wallet) AutoBuildTransaction(request *model.AutoBuildTransactionRequest) *model.AutoBuildTransactionResponse {
	//校验[非找零]收款方
	nonChangePayees := (*request).NonChangePayees
	if nonChangePayees == nil || len(nonChangePayees) == 0 {
		var response model.AutoBuildTransactionResponse
		response.BuildTransactionSuccess = false
		response.Message = "收款方不能为空。"
		return &response
	}
	for _, payee := range nonChangePayees {
		if StringUtil.IsNullOrEmpty(payee.Address) {
			var response model.AutoBuildTransactionResponse
			response.BuildTransactionSuccess = false
			response.Message = "收款方不能为空。"
			return &response
		}
		if payee.Value <= 0 {
			var response model.AutoBuildTransactionResponse
			response.BuildTransactionSuccess = false
			response.Message = "收款方不能为空。"
			return &response
		}
	}
	//创建付款方
	var payers []model.Payer
	//遍历钱包里的账户,用钱包里的账户付款
	allAccounts := w.GetNonZeroBalanceAccounts()
	if allAccounts != nil {
		for _, account := range allAccounts {
			utxo := w.BlockchainDatabase.QueryUnspentTransactionOutputByAddress(account.Address)
			//构建一个新的付款方
			var payer model.Payer
			payer.PrivateKey = account.PrivateKey
			payer.Address = account.Address
			payer.TransactionHash = utxo.TransactionHash
			payer.TransactionOutputIndex = utxo.TransactionOutputIndex
			payer.Value = utxo.Value
			payers = append(payers, payer)
			//设置默认手续费
			fee := uint64(0)
			haveEnoughMoneyToPay := w.haveEnoughMoneyToPay(payers, nonChangePayees, fee)
			if haveEnoughMoneyToPay {
				//创建一个找零账户，并将找零账户保存在钱包里。
				changeAccount := w.CreateAndSaveAccount()
				//创建一个找零收款方
				changePayee := w.createChangePayee(payers, nonChangePayees, changeAccount.Address, fee)
				//创建收款方(收款方=[非找零]收款方+[找零]收款方)
				var payees []model.Payee
				payees = append(payees, nonChangePayees...)
				if changePayee != nil {
					payees = append(payees, *changePayee)
				}
				//构造交易
				var transactionDto dto.TransactionDto
				var response model.AutoBuildTransactionResponse
				response.BuildTransactionSuccess = true
				response.Message = "构建交易成功"
				response.Transaction = transactionDto
				response.TransactionHash = TransactionDtoTool.CalculateTransactionHash(&transactionDto)
				response.Fee = fee
				response.Payers = payers
				response.NonChangePayees = nonChangePayees
				response.ChangePayee = *changePayee
				response.Payees = payees
				return &response
			}
		}
	}
	var response model.AutoBuildTransactionResponse
	response.Message = "没有足够的金额去支付！"
	response.BuildTransactionSuccess = false
	return &response
}

func (w *Wallet) haveEnoughMoneyToPay(payers []model.Payer, payees []model.Payee, fee uint64) bool {
	//计算找零金额
	changeValue := w.changeValue(payers, payees, fee)
	//判断是否有足够的金额去支付
	haveEnoughMoneyToPay := changeValue >= 0
	return haveEnoughMoneyToPay
}
func (w *Wallet) createChangePayee(payers []model.Payer, payees []model.Payee, changeAddress string, fee uint64) *model.Payee {
	//计算找零金额
	changeValue := w.changeValue(payers, payees, fee)
	if changeValue > 0 {
		//构造找零收款方
		var changePayee model.Payee
		changePayee.Address = changeAddress
		changePayee.Value = changeValue
		return &changePayee
	}
	return nil
}

func (w *Wallet) changeValue(payers []model.Payer, payees []model.Payee, fee uint64) uint64 {
	//交易输入总金额
	transactionInputValues := uint64(0)
	for _, payer := range payers {
		transactionInputValues += payer.Value
	}
	//收款方收款总金额
	payeeValues := uint64(0)
	if payees != nil {
		for _, payee := range payees {
			payeeValues += payee.Value
		}
	}
	//计算找零金额，找零金额=交易输入金额-收款方交易输出金额-交易手续费
	changeValue := transactionInputValues - payeeValues - fee
	return changeValue
}
func (w *Wallet) buildTransaction(payers []model.Payer, payees []model.Payee) dto.TransactionDto {
	//构建交易输入
	var transactionInputs []*dto.TransactionInputDto
	for _, payer := range payers {
		var transactionInput *dto.TransactionInputDto
		transactionInput.TransactionHash = payer.TransactionHash
		transactionInput.TransactionOutputIndex = payer.TransactionOutputIndex
		transactionInputs = append(transactionInputs, transactionInput)
	}
	//构建交易输出
	var transactionOutputs []*dto.TransactionOutputDto
	//构造收款方交易输出
	if payees != nil {
		for _, payee := range payees {
			var transactionOutput *dto.TransactionOutputDto
			outputScript := ScriptDtoTool.CreatePayToPublicKeyHashOutputScript(payee.Address)
			transactionOutput.Value = payee.Value
			transactionOutput.OutputScript = outputScript
			transactionOutputs = append(transactionOutputs, transactionOutput)
		}
	}
	//构造交易
	var transaction dto.TransactionDto
	transaction.Inputs = transactionInputs
	transaction.Outputs = transactionOutputs
	//签名
	for i, transactionInput := range transactionInputs {
		account := AccountUtil.AccountFromPrivateKey(payers[i].PrivateKey)
		signature := TransactionDtoTool.Signature(account.PrivateKey, &transaction)
		inputScript := ScriptDtoTool.CreatePayToPublicKeyHashInputScript(signature, account.PublicKey)
		transactionInput.InputScript = inputScript
	}
	return transaction
}
