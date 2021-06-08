package KvDbUtil

import (
	"container/list"

	"github.com/syndtr/goleveldb/leveldb"
)

var dbMap = make(map[string]*leveldb.DB)

func getDb(dbPath string) *leveldb.DB {
	db := dbMap[dbPath]
	if db == nil {
		dbnew, _ := leveldb.OpenFile(dbPath, nil)
		dbMap[dbPath] = dbnew
		return dbnew
	}
	return db
}
func Put(dbPath string, bytesKey []byte, bytesValue []byte) {
	getDb(dbPath).Put(bytesKey, bytesValue, nil)
}
func Delete(dbPath string, bytesKey []byte) {
	getDb(dbPath).Delete(bytesKey, nil)
}
func Get(dbPath string, bytesKey []byte) []byte {
	bytesValue, _ := getDb(dbPath).Get(bytesKey, nil)
	return bytesValue
}
func Gets(dbPath string, from int, size int) *list.List {
	bytes := list.New()
	iter := getDb(dbPath).NewIterator(nil, nil)
	i := 1
	for iter.Next() {
		if i < from {
			continue
		}
		if i > from+size {
			break
		}
		bytes.PushBack(iter.Value())
		i = i + 1
	}
	// iter.Release()
	// err := iter.Error()
	// if err != nil {
	// 	panic(err)
	// }
	return bytes
}
