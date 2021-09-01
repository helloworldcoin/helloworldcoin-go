package framwork

import "helloworld-blockchain-go/util/JsonUtil"

type PageCondition struct {
	from uint64
	size uint64
}

func CreateSuccessResponse(message string, data interface{}) string {
	return "{\"status\":\"SUCCESS\",\"message\":\"message\",\"data\":" + JsonUtil.ToString(data) + "}"
}
