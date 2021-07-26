package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/dto"
)

func ToString(emptyInterface interface{}) string {
	jsonString, _ := json.Marshal(emptyInterface)
	return string(jsonString)
}

func ToObject(value string, object interface{}) interface{} {
	o1, ok := object.(dto.GetBlockRequest)
	if ok {
		json.Unmarshal([]byte(value), &o1)
		return &o1
	}
	o2, ok := object.(dto.PostBlockRequest)
	if ok {
		json.Unmarshal([]byte(value), &o2)
		return &o2
	}
	return nil
}
