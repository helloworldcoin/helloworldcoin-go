package JsonUtil

import (
	"encoding/json"
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/dto"
)

func ToJson(blockDto *dto.BlockDto) string {
	jsonString, _ := json.Marshal(blockDto)
	return string(jsonString)
}

func fromJson() {

}

func ToJsonStringBlock(block *model.Block) string {
	jsonString, _ := json.Marshal(block)
	return string(jsonString)
}
