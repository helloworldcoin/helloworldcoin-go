package controller

import (
	"helloworld-blockchain-go/util/JsonUtil"
	"io"
	"io/ioutil"
	"net/http"
)

func GetRequest(req *http.Request, object interface{}) interface{} {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), object)
	return request
}

func SuccessHttpResponse(rw http.ResponseWriter, message string, response interface{}) {
	s := "{\"status\":\"success\",\"message\":\"" + message + "\",\"data\":" + JsonUtil.ToString(response) + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func FailedHttpResponse(rw http.ResponseWriter, message string) {
	s := "{\"status\":\"fail\",\"message\":\"" + message + "\",\"data\":null" + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
