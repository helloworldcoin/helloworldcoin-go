package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore/po"
)

func ToString(object interface{}) string {
	jsonString, _ := json.Marshal(object)
	return string(jsonString)
}

func ToObject(jsonString string, object interface{}) interface{} {
	o1, ok := object.(dto.GetBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o1)
		return &o1
	}
	o2, ok := object.(dto.PostBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o2)
		return &o2
	}
	o3, ok := object.(dto.GetUnconfirmedTransactionsRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o3)
		return &o3
	}
	o4, ok := object.(dto.GetBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o4)
		return &o4
	}
	o5, ok := object.(dto.PostTransactionRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o5)
		return &o5
	}
	o6, ok := object.(dto.PostBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o6)
		return &o6
	}
	o7, ok := object.(dto.PingRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o7)
		return &o7
	}
	o8, ok := object.(dto.GetNodesRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &o8)
		return &o8
	}

	o9, ok := object.(dto.BlockDto)
	if ok {
		json.Unmarshal([]byte(jsonString), &o9)
		return &o9
	}

	o10, ok := object.(po.NodePo)
	if ok {
		json.Unmarshal([]byte(jsonString), &o10)
		return &o10
	}
	panic("JsonUtil.ToObject can not recognize object type")
}
