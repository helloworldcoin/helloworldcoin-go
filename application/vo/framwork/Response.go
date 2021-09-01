package framwork

import "helloworld-blockchain-go/util/JsonUtil"

func CreateSuccessResponse(message string, data interface{}) string {
	return "{\"status\":\"SUCCESS\",\"message\":\"message\",\"data\":" + JsonUtil.ToString(data) + "}"
}
