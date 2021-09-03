package core

import (
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/util/FileUtil"
	"testing"
)

func TestGetAllAccounts(t *testing.T) {
	FileUtil.DeleteDirectory(ResourcePathTool.GetTestDataRootPath())
	coreConfiguration := &CoreConfiguration{corePath: ResourcePathTool.GetTestDataRootPath()}
	incentive := &Incentive{}
	consensus := &Consensus{}
	virtualMachine := &VirtualMachine{}
	blockchainDatabase := &BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration, VirtualMachine: virtualMachine}
	wallet := &Wallet{CoreConfiguration: coreConfiguration, BlockchainDatabase: blockchainDatabase}

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
	//TODO 校验三个账户
}