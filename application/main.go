package main

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/application/controller"
	"helloworldcoin-go/application/interceptor"
	"helloworldcoin-go/application/service"
	"helloworldcoin-go/application/vo/BlockchainBrowserApplicationApi"
	"helloworldcoin-go/application/vo/NodeConsoleApplicationApi"
	"helloworldcoin-go/application/vo/WalletApplicationApi"
	"helloworldcoin-go/netcore"
	"helloworldcoin-go/util/SystemUtil"
	"io"
	"net/http"
)

func main() {
	SystemUtil.CallDefaultBrowser(`http://localhost/`)

	blockchainNetCore := netcore.NewDefaultBlockchainNetCore()
	blockchainNetCore.Start()

	walletApplicationService := service.NewWalletApplicationService(blockchainNetCore)
	blockchainBrowserApplicationService := service.NewBlockchainBrowserApplicationService(blockchainNetCore)

	blockchainBrowserApplicationController := controller.NewBlockchainBrowserApplicationController(blockchainNetCore, blockchainBrowserApplicationService)
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
	apiMux.HandleFunc(BlockchainBrowserApplicationApi.QUERY_LATEST_10_BLOCKS, blockchainBrowserApplicationController.QueryLatest10Blocks)

	apiMux.HandleFunc(WalletApplicationApi.CREATE_ACCOUNT, walletApplicationController.CreateAccount)
	apiMux.HandleFunc(WalletApplicationApi.CREATE_AND_SAVE_ACCOUNT, walletApplicationController.CreateAndSaveAccount)
	apiMux.HandleFunc(WalletApplicationApi.SAVE_ACCOUNT, walletApplicationController.SaveAccount)
	apiMux.HandleFunc(WalletApplicationApi.DELETE_ACCOUNT, walletApplicationController.DeleteAccount)
	apiMux.HandleFunc(WalletApplicationApi.QUERY_ALL_ACCOUNTS, walletApplicationController.QueryAllAccounts)
	apiMux.HandleFunc(WalletApplicationApi.AUTOMATIC_BUILD_TRANSACTION, walletApplicationController.AutomaticBuildTransaction)
	apiMux.HandleFunc(WalletApplicationApi.SUBMIT_TRANSACTION_TO_BLOCKCHIAIN_NEWWORK, walletApplicationController.SubmitTransactionToBlockchainNetwork)

	apiMux.HandleFunc(NodeConsoleApplicationApi.IS_MINER_ACTIVE, nodeConsoleApplicationController.IsMinerActive)
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
	apiMux.HandleFunc(NodeConsoleApplicationApi.GET_MINER_MINE_MAX_BLOCK_HEIGHT, nodeConsoleApplicationController.GetMinerMineMaxBlockHeight)
	apiMux.HandleFunc(NodeConsoleApplicationApi.SET_MINER_MINE_MAX_BLOCK_HEIGHT, nodeConsoleApplicationController.SetMinerMineMaxBlockHeight)

	apiMux.Handle("/", http.FileServer(http.Dir(SystemUtil.SystemRootDirectory()+"/application/resources/static")))

	ipInterceptorServeMux := http.NewServeMux()
	ipInterceptorServeMux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if interceptor.IsIpAllow(req) {
			apiMux.ServeHTTP(rw, req)
			return
		}
		s := "{\"status\":\"fail\",\"message\":\"service_unauthorized\",\"data\":null" + "}"
		rw.Header().Set("content-type", "text/json")
		io.WriteString(rw, s)
	}))

	http.ListenAndServe(":80", ipInterceptorServeMux)
}
