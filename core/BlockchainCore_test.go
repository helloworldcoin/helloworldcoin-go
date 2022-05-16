package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/tool/ResourceTool"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/JsonUtil"
	"helloworldcoin-go/util/SystemUtil"
	"testing"
)

func TestBlockchainDataFormat(t *testing.T) {
	FileUtil.DeleteDirectory(ResourceTool.GetTestDataRootPath())

	stringBlock1 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block1.json")
	blockDto1 := JsonUtil.ToObject(stringBlock1, dto.BlockDto{}).(*dto.BlockDto)
	stringBlock2 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block2.json")
	blockDto2 := JsonUtil.ToObject(stringBlock2, dto.BlockDto{}).(*dto.BlockDto)
	stringBlock3 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block3.json")
	blockDto3 := JsonUtil.ToObject(stringBlock3, dto.BlockDto{}).(*dto.BlockDto)

	block1Hash := "e213eaae8259e1aca2044f35036ec5fc3c4370a33fa28353a749e8257e1d2e9e"
	block2Hash := "8759b498f57e3b359759b7723850a99968f6e8b4bd8143e2ea41b3dbfbb59942"
	block3Hash := "739f3554dae0a4d2b73142ae8be398fccc8971c9fac52baea1741f4205dc0315"

	blockchainCore := NewBlockchainCore(ResourceTool.GetTestDataRootPath())
	blockchainCore.AddBlockDto(blockDto1)
	blockchainCore.AddBlockDto(blockDto2)
	blockchainCore.AddBlockDto(blockDto3)

	block1 := blockchainCore.QueryBlockByBlockHeight(1)
	if block1Hash != block1.Hash {
		t.Error("test failed")
	}
	block2 := blockchainCore.QueryBlockByBlockHeight(2)
	if block2Hash != block2.Hash {
		t.Error("test failed")
	}
	block3 := blockchainCore.QueryBlockByBlockHeight(3)
	if block3Hash != block3.Hash {
		t.Error("test failed")
	}
}
