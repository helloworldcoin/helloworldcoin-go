package interceptor

import (
	"flag"
	"helloworldcoin-go/util/StringsUtil"
	"net/http"
)

const ALL_IP = "*"

var DEFAULT_ALLOW_IPS []string = []string{"localhost", "127.0.0.1", "0:0:0:0:0:0:0:1"}

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

func getAllowIps() []string {
	var allowIps = flag.String(ALLOW_IPS_KEY, "", "allowIps")
	flag.Parse()
	return StringsUtil.Split(*allowIps, ALLOW_IPS_VALUE_SEPARATOR)
}
