package core

import (
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/JsonUtil"
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

	account1 := AccountUtil.RandomAccount()
	wallet.SaveAccount(account1)
	account2 := AccountUtil.RandomAccount()
	wallet.SaveAccount(account2)

	accounts := wallet.GetAllAccounts()

	if 2 != len(accounts) {
		t.Error("test failed")
	}

	if JsonUtil.ToString(account1) != JsonUtil.ToString(accounts[0]) {
		t.Error("test failed")
	}

	if JsonUtil.ToString(account2) != JsonUtil.ToString(accounts[1]) {
		t.Error("test failed")
	}
}
