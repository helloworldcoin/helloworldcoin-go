package main

import (
	"fmt"
	"helloworldcoin-go/core"
)

func main() {
	coreConfiguration := core.CoreConfiguration{CorePath: "d://abcd"}
	fmt.Println("11111111111111111")
	fmt.Println(coreConfiguration)
	wallet := core.Wallet{CoreConfiguration: coreConfiguration}
	fmt.Println("22222222222222")
	fmt.Println(wallet)
	accout := wallet.CreateAccount()
	fmt.Println("33333333333333")
	fmt.Println(accout)
	wallet.DeleteAccountByAddress(accout.Address)
	accouts := wallet.GetAllAccounts()
	fmt.Println("4444444444444444")
	fmt.Println(accouts)
	wallet.SaveAccount(&accout)
	accouts = wallet.GetAllAccounts()
	fmt.Println("55555555555555")
	fmt.Println(accouts)
}
