package main

import (
	"helloworld-blockchain-go/application/controller"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/SystemUtil"
	"io"
	"net/http"
)

func main() {

	blockchainNetCore := netcore.CreateDefaultBlockchainNetCore()
	blockchainBrowserApplicationController := controller.NewBlockchainBrowserApplicationController(blockchainNetCore, blockchainNetCore.GetBlockchainCore())

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/111", blockchainBrowserApplicationController.Get111)
	apiMux.HandleFunc("/Api/BlockchainBrowserApplication/QueryTop10Blocks", blockchainBrowserApplicationController.QueryTop10Blocks)
	apiMux.HandleFunc("/1/11", get11)
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
func get1(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, "1")
	println(1)
}
func get11(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, "11")
	println(11)
}
