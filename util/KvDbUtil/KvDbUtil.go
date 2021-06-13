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
func Gets(dbPath string, from uint64, size uint64) *list.List {
	bytes := list.New()
	iter := getDb(dbPath).NewIterator(nil, nil)
	i := 1
	for iter.Next() {
		if i < int(from) {
			continue
		}
		if i > int(from+size) {
			break
		}
		bytes.PushBack(iter.Value())
		i = i + 1
	}
	return bytes
}

func Write(dbPath string, kvWriteBatch *KvWriteBatch) {
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
