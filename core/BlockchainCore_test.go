package core

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/core/tool/ResourcePathTool"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/SystemUtil"
	"testing"
)

func TestBlockchainDataFormat(t *testing.T) {
	FileUtil.DeleteDirectory(ResourcePathTool.GetTestDataRootPath())

	stringBlock1 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block1.json")
	block1 := JsonUtil.ToObject(stringBlock1, dto.BlockDto{}).(*dto.BlockDto)
	stringBlock2 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block2.json")
	block2 := JsonUtil.ToObject(stringBlock2, dto.BlockDto{}).(*dto.BlockDto)
	stringBlock3 := FileUtil.Read(SystemUtil.SystemRootDirectory() + "\\core" + "\\test\\resources\\blocks\\block3.json")
	block3 := JsonUtil.ToObject(stringBlock3, dto.BlockDto{}).(*dto.BlockDto)
	block3Hash := "739f3554dae0a4d2b73142ae8be398fccc8971c9fac52baea1741f4205dc0315"

	blockchainCore := NewBlockchainCore(ResourcePathTool.GetTestDataRootPath())
	blockchainCore.AddBlockDto(block1)
	blockchainCore.AddBlockDto(block2)
	blockchainCore.AddBlockDto(block3)

	//若一切正常，此时区块链的最后一个区块就是我们传入的最后一个区块
	tailBlock := blockchainCore.QueryTailBlock()

	//校验区块哈希
	if block3Hash != tailBlock.Hash {
		t.Error("test failed")
	}
}
