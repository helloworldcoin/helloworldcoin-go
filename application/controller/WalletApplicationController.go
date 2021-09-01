package controller

import (
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/netcore"
)

type WalletApplicationController struct {
	blockchainNetCore        *netcore.BlockchainNetCore
	walletApplicationService *service.WalletApplicationService
}

func NewWalletApplicationController(blockchainNetCore *netcore.BlockchainNetCore) *WalletApplicationController {
	var b WalletApplicationController
	b.blockchainNetCore = blockchainNetCore
	return &b
}
