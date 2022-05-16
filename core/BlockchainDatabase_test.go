package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/tool/ResourceTool"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/util/FileUtil"
	"testing"
)

func TestQueryBlockByBlockHeight(t *testing.T) {
	FileUtil.DeleteDirectory(ResourceTool.GetTestDataRootPath())

	transactionOutputDto := dto.TransactionOutputDto{OutputScript: &[]string{"01", "02", "00", "82f46bdbb4550d3c552f1764b53fd0005c81ad3d", "03", "04"}, Value: uint64(5000000000)}
	var outputs []*dto.TransactionOutputDto
	outputs = append(outputs, &transactionOutputDto)

	transactionDto := dto.TransactionDto{Inputs: nil, Outputs: outputs}
	var transactions []*dto.TransactionDto
	transactions = append(transactions, &transactionDto)

	blockDto := &dto.BlockDto{Timestamp: uint64(1623645017179), PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000", Transactions: transactions, Nonce: "c21afb034d3be7f2f233b72aa4136dfcc4ee258af213b91a616bcca1ab780f5b"}

	consensus := &Consensus{}
	incentive := &Incentive{}
	virtualMachine := &VirtualMachine{}
	coreConfiguration := &CoreConfiguration{corePath: ResourceTool.GetTestDataRootPath()}
	blockchainDatabase := NewBlockchainDatabase(consensus, incentive, virtualMachine, coreConfiguration)
	blockchainDatabase.AddBlockDto(blockDto)

	block := blockchainDatabase.QueryBlockByBlockHeight(uint64(1))

	if "80da32d24607e952599eb6dc2b550319ed2052f15009262dc7b5a84a3ca063e0" == block.Hash {
		t.Log("test pass")
	} else {
		t.Error("test failed")
	}
}
