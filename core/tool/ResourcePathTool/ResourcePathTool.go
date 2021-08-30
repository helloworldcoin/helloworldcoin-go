package ResourcePathTool

import (
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/OperateSystemUtil"
)

/**
 * 获取区块链数据存放目录
 */
func GetDataRootPath() string {
	var dataRootPath string
	if OperateSystemUtil.IsWindowsOperateSystem() {
		dataRootPath = "C:\\helloworld-blockchain-go\\"
	} else if OperateSystemUtil.IsMacOperateSystem() {
		dataRootPath = "/tmp/helloworld-blockchain-go/"
	} else {
		dataRootPath = "/tmp/helloworld-blockchain-go/"
	}
	FileUtil.MakeDirectory(dataRootPath)
	return dataRootPath
}

/**
 * 获取测试区块链数据存放目录
 */
func GetTestDataRootPath() string {
	var dataRootPath string
	if OperateSystemUtil.IsWindowsOperateSystem() {
		dataRootPath = "C:\\helloworld-blockchain-go-test\\"
	} else if OperateSystemUtil.IsMacOperateSystem() {
		dataRootPath = "/tmp/helloworld-blockchain-go-test/"
	} else {
		dataRootPath = "/tmp/helloworld-blockchain-go-test/"
	}
	FileUtil.MakeDirectory(dataRootPath)
	return dataRootPath
}
