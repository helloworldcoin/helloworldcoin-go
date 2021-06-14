package main

import (
	"fmt"
	"helloworldcoin-go/core"
	"helloworldcoin-go/core/tool/EncodeDecodeTool"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/util/JsonUtil"
)

func main3() {

	transactionOutputDto := dto.TransactionOutputDto{OutputScript: []string{"01", "02", "00", "82f46bdbb4550d3c552f1764b53fd0005c81ad3d", "03", "04"}, Value: uint64(5000000000)}
	outputs := []dto.TransactionOutputDto{}
	outputs = append(outputs, transactionOutputDto)

	transactionDto := dto.TransactionDto{Inputs: nil, Outputs: outputs}
	transactions := []dto.TransactionDto{}
	transactions = append(transactions, transactionDto)

	blockDto := &dto.BlockDto{Timestamp: uint64(1623645017179), PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000", Transactions: transactions, Nonce: "c21afb034d3be7f2f233b72aa4136dfcc4ee258af213b91a616bcca1ab780f5b"}

	fmt.Println(*blockDto)
	fmt.Println(EncodeDecodeTool.EncodeBlockDto(blockDto))
	fmt.Println(JsonUtil.ToJson(blockDto))

	consensus := &core.Consensus{}
	incentive := &core.Incentive{}
	coreConfiguration := &core.CoreConfiguration{CorePath: "d:"}
	blockchainDatabase := core.BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration}
	fmt.Println(blockchainDatabase)
	blockchainDatabase.AddBlockDto(blockDto)
	fmt.Println("1---------------------------------------------------------------")
	block := blockchainDatabase.QueryBlockByBlockHeight(uint64(1))
	fmt.Println(JsonUtil.ToJsonStringBlock(block))
	fmt.Println("2---------------------------------------------------------------")

}
