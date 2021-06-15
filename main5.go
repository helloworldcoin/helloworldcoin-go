package main

import (
	"fmt"
	"helloworldcoin-go/core"
)

func main() {

	consensus := &core.Consensus{}
	incentive := &core.Incentive{}
	coreConfiguration := &core.CoreConfiguration{CorePath: "C:\\HelloworldBlockchainDataGo"}
	blockchainDatabase := &core.BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration}
	wallet := &core.Wallet{CoreConfiguration: coreConfiguration}
	unconfirmedTransactionDatabase := &core.UnconfirmedTransactionDatabase{CoreConfiguration: coreConfiguration}
	miner := core.Miner{CoreConfiguration: coreConfiguration, Wallet: wallet, BlockchainDatabase: blockchainDatabase, UnconfirmedTransactionDatabase: unconfirmedTransactionDatabase}
	fmt.Println(miner)
	miner.Start()
}
