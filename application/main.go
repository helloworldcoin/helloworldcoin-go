package main

import (
	"helloworld-blockchain-go/application/controller"
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/application/vo/BlockchainBrowserApplicationApi"
	"helloworld-blockchain-go/application/vo/NodeConsoleApplicationApi"
	"helloworld-blockchain-go/application/vo/WalletApplicationApi"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/SystemUtil"
	"net/http"
)

func main() {
	SystemUtil.CallDefaultBrowser(`http://localhost:8080/`)

	blockchainNetCore := netcore.CreateDefaultBlockchainNetCore()
	blockchainNetCore.Start()

	walletApplicationService := service.NewWalletApplicationService(blockchainNetCore)

	blockchainBrowserApplicationController := controller.NewBlockchainBrowserApplicationController(blockchainNetCore)
	nodeConsoleApplicationController := controller.NewNodeConsoleApplicationController(blockchainNetCore)
	walletApplicationController := controller.NewWalletApplicationController(blockchainNetCore, walletApplicationService)
	apiMux := http.NewServeMux()

	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_TRANSACTION_BY_TRANSACTION_HASH, blockchainBrowserApplicationController.QueryTransactionByTransactionHash)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_TRANSACTIONS_BY_BLOCK_HASH_TRANSACTION_HEIGHT, blockchainBrowserApplicationController.QueryTransactionsByBlockHashTransactionHeight)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_TRANSACTION_OUTPUT_BY_ADDRESS, blockchainBrowserApplicationController.QueryTransactionOutputByAddress)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_TRANSACTION_OUTPUT_BY_TRANSACTION_OUTPUT_ID, blockchainBrowserApplicationController.QueryTransactionOutputByTransactionOutputId)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_BLOCKCHAIN_HEIGHT, blockchainBrowserApplicationController.QueryBlockchainHeight)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_UNCONFIRMED_TRANSACTION_BY_TRANSACTION_HASH, blockchainBrowserApplicationController.QueryUnconfirmedTransactionByTransactionHash)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_UNCONFIRMED_TRANSACTIONS, blockchainBrowserApplicationController.QueryUnconfirmedTransactions)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_BLOCK_BY_BLOCK_HEIGHT, blockchainBrowserApplicationController.QueryBlockByBlockHeight)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_BLOCK_BY_BLOCK_HASH, blockchainBrowserApplicationController.QueryBlockByBlockHash)
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_TOP10_BLOCKS, blockchainBrowserApplicationController.QueryTop10Blocks)

	apiMux.HandleFunc(WalletApplicationApi.CREATE_ACCOUNT, walletApplicationController.CreateAccount)
	apiMux.HandleFunc(WalletApplicationApi.CREATE_AND_SAVE_ACCOUNT, walletApplicationController.CreateAndSaveAccount)
	apiMux.HandleFunc(WalletApplicationApi.SAVE_ACCOUNT, walletApplicationController.SaveAccount)
	apiMux.HandleFunc(WalletApplicationApi.DELETE_ACCOUNT, walletApplicationController.DeleteAccount)
	apiMux.HandleFunc(WalletApplicationApi.QUERY_ALL_ACCOUNTS, walletApplicationController.QueryAllAccounts)
	apiMux.HandleFunc(WalletApplicationApi.AUTO_BUILD_TRANSACTION, walletApplicationController.AutoBuildTransaction)
	apiMux.HandleFunc(WalletApplicationApi.SUBMIT_TRANSACTION_TO_BLOCKCHIAIN_NEWWORK, walletApplicationController.SubmitTransactionToBlockchainNetwork)

	apiMux.HandleFunc(NodeConsoleApplicationApi.IS_MINER_ACTIVE, nodeConsoleApplicationController.IsMineActive)
	apiMux.HandleFunc(NodeConsoleApplicationApi.ACTIVE_MINER, nodeConsoleApplicationController.ActiveMiner)
	apiMux.HandleFunc(NodeConsoleApplicationApi.DEACTIVE_MINER, nodeConsoleApplicationController.DeactiveMiner)

	apiMux.HandleFunc(NodeConsoleApplicationApi.IS_AUTO_SEARCH_BLOCK, nodeConsoleApplicationController.IsAutoSearchBlock)
	apiMux.HandleFunc(NodeConsoleApplicationApi.ACTIVE_AUTO_SEARCH_BLOCK, nodeConsoleApplicationController.ActiveAutoSearchBlock)
	apiMux.HandleFunc(NodeConsoleApplicationApi.DEACTIVE_AUTO_SEARCH_BLOCK, nodeConsoleApplicationController.DeactiveAutoSearchBlock)

	apiMux.HandleFunc(NodeConsoleApplicationApi.ADD_NODE, nodeConsoleApplicationController.AddNode)
	apiMux.HandleFunc(NodeConsoleApplicationApi.UPDATE_NODE, nodeConsoleApplicationController.UpdateNode)
	apiMux.HandleFunc(NodeConsoleApplicationApi.DELETE_NODE, nodeConsoleApplicationController.DeleteNode)
	apiMux.HandleFunc(NodeConsoleApplicationApi.QUERY_ALL_NODES, nodeConsoleApplicationController.QueryAllNodes)

	apiMux.HandleFunc(NodeConsoleApplicationApi.IS_AUTO_SEARCH_NODE, nodeConsoleApplicationController.IsAutoSearchNode)
	apiMux.HandleFunc(NodeConsoleApplicationApi.ACTIVE_AUTO_SEARCH_NODE, nodeConsoleApplicationController.ActiveAutoSearchNode)
	apiMux.HandleFunc(NodeConsoleApplicationApi.DEACTIVE_AUTO_SEARCH_NODE, nodeConsoleApplicationController.DeactiveAutoSearchNode)

	apiMux.HandleFunc(NodeConsoleApplicationApi.DELETE_BLOCKS, nodeConsoleApplicationController.DeleteBlocks)
	apiMux.HandleFunc(NodeConsoleApplicationApi.GET_MAX_BLOCK_HEIGHT, nodeConsoleApplicationController.GetMaxBlockHeight)
	apiMux.HandleFunc(NodeConsoleApplicationApi.SET_MAX_BLOCK_HEIGHT, nodeConsoleApplicationController.SetMaxBlockHeight)

	apiMux.Handle("/", http.FileServer(http.Dir(SystemUtil.SystemRootDirectory()+"\\application\\resources\\static")))

	ipInterceptorServeMux := http.NewServeMux()
	ipInterceptorServeMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host != "localhost:8080" {
			http.Error(w, "Blocked", 401)
			return
		}
		apiMux.ServeHTTP(w, req)
	}))

	http.ListenAndServe(":8080", ipInterceptorServeMux)
}
