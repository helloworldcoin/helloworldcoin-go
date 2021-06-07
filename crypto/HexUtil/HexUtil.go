package HexUtil

import (
	"encoding/hex"
)

func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

func BytesToHexString(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	return hexString
}
