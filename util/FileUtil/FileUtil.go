package FileUtil

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/util/SystemUtil"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewPath(parent string, child string) string {
	return filepath.Join(parent, child)
}

func MakeDirectory(path string) {
	if !isExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			SystemUtil.ErrorExit("create directory failed. path is "+path+".", err)
		}
	}
}

func DeleteDirectory(path string) {
	error := os.RemoveAll(path)
	if error != nil {
		SystemUtil.ErrorExit("delete directory failed. path is "+path+".", nil)
	}
}

func Read(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	bytesData, _ := ioutil.ReadAll(f)
	return string(bytesData)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
