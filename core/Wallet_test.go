package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/tool/ResourceTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/util/FileUtil"
	"testing"
)

func TestGetAllAccounts(t *testing.T) {
	FileUtil.DeleteDirectory(ResourceTool.GetTestDataRootPath())
	coreConfiguration := &CoreConfiguration{corePath: ResourceTool.GetTestDataRootPath()}
	incentive := &Incentive{}
	consensus := &Consensus{}
	virtualMachine := &VirtualMachine{}
	blockchainDatabase := NewBlockchainDatabase(consensus, incentive, virtualMachine, coreConfiguration)
	wallet := NewWallet(coreConfiguration, blockchainDatabase)

	account0 := AccountUtil.RandomAccount()
	wallet.SaveAccount(account0)
	account1 := AccountUtil.RandomAccount()
	wallet.SaveAccount(account1)
	account2 := AccountUtil.RandomAccount()
	wallet.SaveAccount(account2)

	accounts := wallet.GetAllAccounts()

	if 3 != len(accounts) {
		t.Error("test failed")
	}
	//TODO test
}
