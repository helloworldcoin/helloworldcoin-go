package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/tool/ScriptDtoTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/EncodeDecodeTool"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
)

const WALLET_DATABASE_NAME = "WalletDatabase"

type Wallet struct {
	coreConfiguration  *CoreConfiguration
	blockchainDatabase *BlockchainDatabase
}

func NewWallet(coreConfiguration *CoreConfiguration, blockchainDatabase *BlockchainDatabase) *Wallet {
	var wallet Wallet
	wallet.coreConfiguration = coreConfiguration
	wallet.blockchainDatabase = blockchainDatabase
	return &wallet
}

func (w *Wallet) GetAllAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	bytesAccounts := KvDbUtil.Gets(w.getWalletDatabasePath(), 1, 100000000)
	for e := bytesAccounts.Front(); e != nil; e = e.Next() {
		account := EncodeDecodeTool.Decode(e.Value.([]byte), AccountUtil.Account{}).(*AccountUtil.Account)
		accounts = append(accounts, account)
	}
	return accounts
}
func (w *Wallet) GetNonZeroBalanceAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	bytesAccounts := KvDbUtil.Gets(w.getWalletDatabasePath(), 1, 100000000)
	for e := bytesAccounts.Front(); e != nil; e = e.Next() {
		account := EncodeDecodeTool.Decode(e.Value.([]byte), AccountUtil.Account{}).(*AccountUtil.Account)
		utxo := w.blockchainDatabase.QueryUnspentTransactionOutputByAddress(account.Address)
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
	KvDbUtil.Put(w.getWalletDatabasePath(), w.getKeyByAccount(account), EncodeDecodeTool.Encode(account))
}
func (w *Wallet) DeleteAccountByAddress(address string) {
	KvDbUtil.Delete(w.getWalletDatabasePath(), w.getKeyByAddress(address))
}
func (w *Wallet) GetBalanceByAddress(address string) uint64 {
	utxo := w.blockchainDatabase.QueryUnspentTransactionOutputByAddress(address)
	if utxo != nil {
		return utxo.Value
	}
	return uint64(0)
}
func (w *Wallet) getWalletDatabasePath() string {
	return FileUtil.NewPath(w.coreConfiguration.getCorePath(), WALLET_DATABASE_NAME)
}
func (w *Wallet) getKeyByAddress(address string) []byte {
	return ByteUtil.StringToUtf8Bytes(address)
}
func (w *Wallet) getKeyByAccount(account *AccountUtil.Account) []byte {
	return ByteUtil.StringToUtf8Bytes(account.Address)
}

func (w *Wallet) AutoBuildTransaction(request *model.AutoBuildTransactionRequest) *model.AutoBuildTransactionResponse {
	nonChangePayees := (*request).NonChangePayees
	var payers []*model.Payer
	allAccounts := w.GetNonZeroBalanceAccounts()
	if allAccounts != nil {
		for _, account := range allAccounts {
			utxo := w.blockchainDatabase.QueryUnspentTransactionOutputByAddress(account.Address)
			var payer model.Payer
			payer.PrivateKey = account.PrivateKey
			payer.Address = account.Address
			payer.TransactionHash = utxo.TransactionHash
			payer.TransactionOutputIndex = utxo.TransactionOutputIndex
			payer.Value = utxo.Value
			payers = append(payers, &payer)
			fee := uint64(0)
			haveEnoughMoneyToPay := w.haveEnoughMoneyToPay(payers, nonChangePayees, fee)
			if haveEnoughMoneyToPay {
				changeAccount := w.CreateAndSaveAccount()
				changePayee := w.createChangePayee(payers, nonChangePayees, changeAccount.Address, fee)
				var payees []*model.Payee
				payees = append(payees, nonChangePayees...)
				if changePayee != nil {
					payees = append(payees, changePayee)
				}
				transactionDto := w.buildTransaction(payers, payees)
				var response model.AutoBuildTransactionResponse
				response.BuildTransactionSuccess = true
				response.Transaction = transactionDto
				response.TransactionHash = TransactionDtoTool.CalculateTransactionHash(transactionDto)
				response.Fee = fee
				response.Payers = payers
				response.NonChangePayees = nonChangePayees
				response.ChangePayee = changePayee
				response.Payees = payees
				return &response
			}
		}
	}
	var response model.AutoBuildTransactionResponse
	response.BuildTransactionSuccess = false
	return &response
}

func (w *Wallet) haveEnoughMoneyToPay(payers []*model.Payer, payees []*model.Payee, fee uint64) bool {
	changeValue := w.changeValue(payers, payees, fee)
	haveEnoughMoneyToPay := changeValue >= 0
	return haveEnoughMoneyToPay
}
func (w *Wallet) createChangePayee(payers []*model.Payer, payees []*model.Payee, changeAddress string, fee uint64) *model.Payee {
	changeValue := w.changeValue(payers, payees, fee)
	if changeValue > 0 {
		var changePayee model.Payee
		changePayee.Address = changeAddress
		changePayee.Value = changeValue
		return &changePayee
	}
	return nil
}

func (w *Wallet) changeValue(payers []*model.Payer, payees []*model.Payee, fee uint64) uint64 {
	transactionInputValues := uint64(0)
	for _, payer := range payers {
		transactionInputValues += payer.Value
	}
	payeeValues := uint64(0)
	if payees != nil {
		for _, payee := range payees {
			payeeValues += payee.Value
		}
	}
	changeValue := transactionInputValues - payeeValues - fee
	return changeValue
}
func (w *Wallet) buildTransaction(payers []*model.Payer, payees []*model.Payee) *dto.TransactionDto {
	var transactionInputs []*dto.TransactionInputDto
	for _, payer := range payers {
		var transactionInput dto.TransactionInputDto
		transactionInput.TransactionHash = payer.TransactionHash
		transactionInput.TransactionOutputIndex = payer.TransactionOutputIndex
		transactionInputs = append(transactionInputs, &transactionInput)
	}
	var transactionOutputs []*dto.TransactionOutputDto
	if payees != nil {
		for _, payee := range payees {
			var transactionOutput dto.TransactionOutputDto
			outputScript := ScriptDtoTool.CreatePayToPublicKeyHashOutputScript(payee.Address)
			transactionOutput.Value = payee.Value
			transactionOutput.OutputScript = outputScript
			transactionOutputs = append(transactionOutputs, &transactionOutput)
		}
	}
	var transaction dto.TransactionDto
	transaction.Inputs = transactionInputs
	transaction.Outputs = transactionOutputs
	for i, transactionInput := range transactionInputs {
		account := AccountUtil.AccountFromPrivateKey(payers[i].PrivateKey)
		signature := TransactionDtoTool.Signature(account.PrivateKey, &transaction)
		inputScript := ScriptDtoTool.CreatePayToPublicKeyHashInputScript(signature, account.PublicKey)
		transactionInput.InputScript = inputScript
	}
	return &transaction
}
