package ByteUtil

import (
	"bytes"
	"encoding/binary"
	"helloworldcoin-go/crypto/HexUtil"
)

func StringToUtf8Bytes(stringValue string) []byte {
	//TODO is utf8?
	return []byte(stringValue)
}
func Utf8BytesToString(bytesValue []byte) string {
	//TODO is utf8?
	return string(bytesValue)
}

func Concat(arrays ...[]byte) []byte {
	return bytes.Join(arrays, []byte(""))
}

func ConcatLength(value []byte) []byte {
	return Concat(Long8ToByte8(uint64(len(value))), value)
}

func Long8ToByte8(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, number)
	return bytes
}

func Byte8ToLong8(bytes []byte) uint64 {
	number := uint64(binary.BigEndian.Uint64(bytes))
	return number
}

func Equals(a []byte, b []byte) bool {
	return true
}
func Flat(arraysarrays [][]byte) []byte {
	concatBytes := []byte{}
	for _, value := range arraysarrays {
		concatBytes = Concat(concatBytes, value)
	}
	return concatBytes
}
func FlatAndConcatLength(arraysarrays [][]byte) []byte {
	flatBytes := Flat(arraysarrays)
	return ConcatLength(flatBytes)
}
func Copy(src []byte, srcPos int, destPos int) []byte {
	length := destPos - srcPos
	return src[srcPos:length]
}
func CopyTo(src []byte, srcPos int, dest []byte, destPos int, length int) {
	for len(dest) < destPos+length {
		dest = append(dest, byte(0x00))
	}
	for i := 0; i < length; i = i + 1 {
		dest[destPos+i] = src[srcPos+i]
	}
}
func Long8ToHexString8(number uint64) string {
	return HexUtil.BytesToHexString(Long8ToByte8(number))
}
