package FileUtil

/*
 @author king 409060350@qq.com
*/

import (
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
			panic(err)
		}
	}
}
func DeleteDirectory(path string) {
	error := os.RemoveAll(path)
	if error != nil {
		panic(error)
	}
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

func Read(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	bytesData, _ := ioutil.ReadAll(f)
	return string(bytesData)
}
