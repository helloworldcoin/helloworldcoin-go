package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/core/Model"
)

func ToString(emptyStruct interface{}) string {
	jsonString, _ := json.Marshal(emptyStruct)
	return string(jsonString)
}

func toObject() {

}

func ToStringBlock(block *Model.Block) string {
	jsonString, _ := json.Marshal(block)
	return string(jsonString)
}
