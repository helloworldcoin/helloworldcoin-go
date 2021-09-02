package framwork

import "helloworld-blockchain-go/util/JsonUtil"

type PageCondition struct {
	From uint64
	Size uint64
}

func CreateSuccessResponse(message string, data interface{}) string {
	return "{\"status\":\"SUCCESS\",\"message\":\"message\",\"data\":" + JsonUtil.ToString(data) + "}"
}
