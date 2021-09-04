package KvDbUtil

/*
 @author king 409060350@qq.com
*/

import (
	"container/list"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var dbMap = make(map[string]*leveldb.DB)
var Mutex = sync.Mutex{}

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
	Mutex.Lock()
	defer Mutex.Unlock()
	getDb(dbPath).Put(bytesKey, bytesValue, nil)
}
func Delete(dbPath string, bytesKey []byte) {
	Mutex.Lock()
	defer Mutex.Unlock()
	getDb(dbPath).Delete(bytesKey, nil)
}
func Get(dbPath string, bytesKey []byte) []byte {
	Mutex.Lock()
	defer Mutex.Unlock()
	bytesValue, _ := getDb(dbPath).Get(bytesKey, nil)
	return bytesValue
}
func Gets(dbPath string, from uint64, size uint64) *list.List {
	Mutex.Lock()
	defer Mutex.Unlock()
	bytes := list.New()
	iterator := getDb(dbPath).NewIterator(nil, nil)
	i := 1
	for iterator.Next() {
		if i < int(from) {
			continue
		}
		if i > int(from+size) {
			break
		}
		//deep clone value
		bytes.PushBack([]byte(string(iterator.Value())))
		i = i + 1
	}
	return bytes
}

func Write(dbPath string, kvWriteBatch *KvWriteBatch) {
	Mutex.Lock()
	defer Mutex.Unlock()
	batch := new(leveldb.Batch)
	kvWrites := kvWriteBatch.GetKvWrites()
	for _, kvWrite := range kvWrites {
		if kvWrite.kvWriteActionEnum == ADD {
			batch.Put(kvWrite.key, kvWrite.value)
		} else if kvWrite.kvWriteActionEnum == ADD {
			batch.Delete(kvWrite.key)
		}
	}
	getDb(dbPath).Write(batch, nil)
}

type KvWriteActionEnum = bool

const (
	ADD    KvWriteActionEnum = true
	DELETE KvWriteActionEnum = false
)

type KvWrite struct {
	kvWriteActionEnum KvWriteActionEnum
	key               []byte
	value             []byte
}
type KvWriteBatch struct {
	kvWrites []KvWrite
	key      []byte
	value    []byte
}

func (kvWriteBatch *KvWriteBatch) GetKvWrites() []KvWrite {
	return kvWriteBatch.kvWrites
}
func (kvWriteBatch *KvWriteBatch) SetKvWrites(kvWrites []KvWrite) {
	kvWriteBatch.kvWrites = kvWrites
}
func (kvWriteBatch *KvWriteBatch) Put(key []byte, value []byte) {
	kvWrite := KvWrite{kvWriteActionEnum: ADD, key: key, value: value}
	kvWriteBatch.kvWrites = append(kvWriteBatch.kvWrites, kvWrite)
}
func (kvWriteBatch *KvWriteBatch) Delete(key []byte) {
	kvWrite := KvWrite{kvWriteActionEnum: DELETE, key: key, value: nil}
	kvWriteBatch.kvWrites = append(kvWriteBatch.kvWrites, kvWrite)
}
