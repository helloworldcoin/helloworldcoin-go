package main

import (
	"helloworld-blockchain-go/application/controller"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/SystemUtil"
	"net/http"
)

func main() {
	SystemUtil.CallDefaultBrowser(`http://localhost:8080/`)

	blockchainNetCore := netcore.CreateDefaultBlockchainNetCore()
	blockchainNetCore.Start()

	blockchainBrowserApplicationController := controller.NewBlockchainBrowserApplicationController(blockchainNetCore, blockchainNetCore.GetBlockchainCore())

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/Api/BlockchainBrowserApplication/QueryTop10Blocks", blockchainBrowserApplicationController.QueryTop10Blocks)
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
