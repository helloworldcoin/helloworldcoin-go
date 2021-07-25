package core

import (
	"helloworld-blockchain-go/core/tool/EncodeDecodeTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
)

const WALLET_DATABASE_NAME = "WalletDatabase"

type Wallet struct {
	CoreConfiguration  *CoreConfiguration
	BlockchainDatabase *BlockchainDatabase
}

func (w *Wallet) GetAllAccounts() []*AccountUtil.Account {
	var accounts []*AccountUtil.Account
	list := KvDbUtil.Gets(w.getWalletDatabasePath(), 1, 11)
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
