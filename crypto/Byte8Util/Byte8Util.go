package Byte8Util

import (
	"encoding/binary"
	"helloworldcoin-go/crypto/HexUtil"
)

func Uint64ToByte8(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, number)
	return bytes
}

func Uint64ToHexString64(number uint64) string {
	return HexUtil.BytesToHexString(Uint64ToByte8(number))
}
