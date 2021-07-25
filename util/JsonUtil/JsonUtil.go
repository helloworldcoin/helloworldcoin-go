package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/core/Model"
	"helloworld-blockchain-go/dto"
)

func ToString(blockDto *dto.BlockDto) string {
	jsonString, _ := json.Marshal(blockDto)
	return string(jsonString)
}

func toObject() {

}

func ToStringBlock(block *Model.Block) string {
	jsonString, _ := json.Marshal(block)
	return string(jsonString)
}
