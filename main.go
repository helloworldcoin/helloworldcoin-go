package main

import "helloworld-blockchain-go/netcore"

func main() {
	blockchainNetCore := netcore.CreateDefaultBlockchainNetCore()
	blockchainNetCore.Start()
}
