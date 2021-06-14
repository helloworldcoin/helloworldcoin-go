package core

import (
	"helloworldcoin-go/core/tool/EncodeDecodeTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
)

const WALLET_DATABASE_NAME = "WalletDatabase"

type Wallet struct {
	CoreConfiguration *CoreConfiguration
}

func (w *Wallet) GetAllAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	list := KvDbUtil.Gets(w.GetWalletDatabasePath(), 1, 11)
	for e := list.Front(); e != nil; e = e.Next() {
		account := EncodeDecodeTool.DecodeToAccount(e.Value.([]byte))
		accounts = append(accounts, account)
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
	KvDbUtil.Put(w.GetWalletDatabasePath(), getKeyByAccount(account), EncodeDecodeTool.EncodeAccount(account))
}
func (w *Wallet) DeleteAccountByAddress(address string) {
	KvDbUtil.Delete(w.GetWalletDatabasePath(), getKeyByAddress(address))
}

func (w *Wallet) GetWalletDatabasePath() string {
	return FileUtil.NewPath(w.CoreConfiguration.getCorePath(), WALLET_DATABASE_NAME)
}
func getKeyByAddress(address string) []byte {
	return ByteUtil.StringToUtf8Bytes(address)
}
func getKeyByAccount(account *AccountUtil.Account) []byte {
	return ByteUtil.StringToUtf8Bytes(account.Address)
}
