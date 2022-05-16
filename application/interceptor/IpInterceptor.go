package interceptor

import (
	"flag"
	"helloworldcoin-go/util/StringsUtil"
	"net/http"
)

//*代表允许所有ip访问。
const ALL_IP = "*"

//默认允许访问的ip列表。
var DEFAULT_ALLOW_IPS []string = []string{"localhost", "127.0.0.1", "0:0:0:0:0:0:0:1"}

//允许的ip列表，多个ip之间以分隔符逗号(,)进行分割分隔。
const ALLOW_IPS_KEY = "allowIps"
const ALLOW_IPS_VALUE_SEPARATOR = ","

func IsIpAllow(req *http.Request) bool {
	remoteHost := req.Host
	if StringsUtil.Contains(&DEFAULT_ALLOW_IPS, remoteHost) {
		return true
	}
	allowIps := getAllowIps()
	if allowIps != nil && len(allowIps) != 0 {
		if StringsUtil.Contains(&allowIps, ALL_IP) {
			return true
		}
		if StringsUtil.Contains(&allowIps, remoteHost) {
			return true
		}
	}
	return false
}

//获取允许的ip列表
func getAllowIps() []string {
	var allowIps = flag.String(ALLOW_IPS_KEY, "", "allowIps")
	flag.Parse()
	return StringsUtil.Split(*allowIps, ALLOW_IPS_VALUE_SEPARATOR)
}
