package core

import (
	"fmt"
	"testing"
)

func TestWallet(t *testing.T) {
	//TODO
	coreConfiguration := CoreConfiguration{CorePath: "d://abcd"}
	fmt.Println(coreConfiguration)
	wallet := Wallet{CoreConfiguration: &coreConfiguration}
	fmt.Println(wallet)
	account := wallet.CreateAccount()
	fmt.Println(account)
	wallet.DeleteAccountByAddress(account.Address)
	accounts := wallet.GetAllAccounts()
	fmt.Println(accounts)
	wallet.SaveAccount(account)
	accounts = wallet.GetAllAccounts()
	fmt.Println(accounts)
}
