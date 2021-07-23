package JsonUtil

import (
	"encoding/json"
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/dto"
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
