package main

import (
	"fmt"
	"helloworld-blockchain-go/util/SystemUtil"
	"io"
	"log"
	"net/http"
)

func main() {

	fmt.Println("SystemRootDirectory：" + SystemUtil.SystemRootDirectory())

	//blockchainNetCore := netcore.CreateDefaultBlockchainNetCore()
	//blockchainNetCore.Start()
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/1", get1)
	apiMux.HandleFunc("/1/11", get11)

	http.FileServer(http.Dir("C:\\Users\\xingkaichun\\IdeaProjects\\helloworld-blockchain-java\\helloworld-blockchain-application\\src\\main\\resources\\static"))
	err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host != "localhost:8080" {
			http.Error(w, "Blocked", 401)
			return
		}
		// pass the request to the mux
		apiMux.ServeHTTP(w, req)
		println("**************")
	}))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	print(4)
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
