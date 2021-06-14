package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func myHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(con))
	w.Write([]byte("Hello World"))
}
func main() {
	http.HandleFunc("/", myHandle)
	http.ListenAndServe(":8888", nil)
}
