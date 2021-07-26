package main

import (
	"fmt"
	"helloworld-blockchain-go/core"
	"helloworld-blockchain-go/core/Model/ModelWallet"
	"helloworld-blockchain-go/util/JsonUtil"
)

func main() {
	consensus := &core.Consensus{}
	incentive := &core.Incentive{}
	coreConfiguration := &core.CoreConfiguration{CorePath: "C:\\HelloworldBlockchainDataGo"}
	blockchainDatabase := &core.BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration}
	wallet := &core.Wallet{CoreConfiguration: coreConfiguration, BlockchainDatabase: blockchainDatabase}
	unconfirmedTransactionDatabase := &core.UnconfirmedTransactionDatabase{CoreConfiguration: coreConfiguration}
	miner := &core.Miner{CoreConfiguration: coreConfiguration, Wallet: wallet, BlockchainDatabase: blockchainDatabase, UnconfirmedTransactionDatabase: unconfirmedTransactionDatabase}
	//miner.Start()

	var request ModelWallet.AutoBuildTransactionRequest
	var nonChangePayees []ModelWallet.Payee
	payee := ModelWallet.Payee{Address: "1FJXNFnyErEgHm5kyADKSoTFxVnaAUoQHq", Value: 8888888}
	nonChangePayees = append(nonChangePayees, payee)
	request.NonChangePayees = nonChangePayees
	response := wallet.AutoBuildTransaction(&request)
	fmt.Println(JsonUtil.ToString(response))

	unconfirmedTransactionDatabase.InsertTransaction(&response.Transaction)
	miner.Start()
}
