package controller

import (
	"helloworldcoin-go/util/JsonUtil"
	"io"
	"io/ioutil"
	"net/http"
)

func GetRequest(req *http.Request, object interface{}) interface{} {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), object)
	return request
}

func success(rw http.ResponseWriter, response interface{}) {
	s := "{\"status\":\"success\",\"message\":null,\"data\":" + JsonUtil.ToString(response) + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func fail(rw http.ResponseWriter, message string) {
	s := "{\"status\":\"fail\",\"message\":\"" + message + "\",\"data\":null" + "}"
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func serviceUnavailable(rw http.ResponseWriter) {
	fail(rw, "service_unavailable")
}
func ServiceUnauthorized(rw http.ResponseWriter) {
	fail(rw, "service_unauthorized")
}
