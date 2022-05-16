package ResourceTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/SystemUtil"
)

func GetDataRootPath() string {
	var dataRootPath string
	if SystemUtil.IsWindowsOperateSystem() {
		dataRootPath = "C:\\helloworldcoin-go\\"
	} else if SystemUtil.IsMacOperateSystem() {
		dataRootPath = "/tmp/helloworldcoin-go/"
	} else if SystemUtil.IsLinuxOperateSystem() {
		dataRootPath = "/tmp/helloworldcoin-go/"
	} else {
		dataRootPath = "/tmp/helloworldcoin-go/"
	}
	FileUtil.MakeDirectory(dataRootPath)
	return dataRootPath
}

func GetTestDataRootPath() string {
	var dataRootPath string
	if SystemUtil.IsWindowsOperateSystem() {
		dataRootPath = "C:\\helloworldcoin-go-test\\"
	} else if SystemUtil.IsMacOperateSystem() {
		dataRootPath = "/tmp/helloworldcoin-go-test/"
	} else if SystemUtil.IsLinuxOperateSystem() {
		dataRootPath = "/tmp/helloworldcoin-go-test/"
	} else {
		dataRootPath = "/tmp/helloworldcoin-go-test/"
	}
	FileUtil.MakeDirectory(dataRootPath)
	return dataRootPath
}
