package ByteUtil

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
)

func BytesToHexString(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	return hexString
}
func HexStringToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

func Uint64ToBytes(number uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, number)
	return bytes
}
func BytesToUint64(bytes []byte) uint64 {
	number := uint64(binary.BigEndian.Uint64(bytes))
	return number
}

func StringToUtf8Bytes(stringValue string) []byte {
	return []byte(stringValue)
}
func Utf8BytesToString(bytesValue []byte) string {
	return string(bytesValue)
}

func Concatenate(bytes1 []byte, bytes2 []byte) []byte {
	return bytes.Join([][]byte{bytes1, bytes2}, []byte(""))
}
func Concatenate3(bytes1 []byte, bytes2 []byte, bytes3 []byte) []byte {
	return bytes.Join([][]byte{bytes1, bytes2, bytes3}, []byte(""))
}
func Concatenate4(bytes1 []byte, bytes2 []byte, bytes3 []byte, bytes4 []byte) []byte {
	return bytes.Join([][]byte{bytes1, bytes2, bytes3, bytes4}, []byte(""))
}

func ConcatenateLength(value []byte) []byte {
	return Concatenate(Uint64ToBytes(uint64(len(value))), value)
}

func Flat(values [][]byte) []byte {
	var concatBytes []byte
	for _, value := range values {
		concatBytes = Concatenate(concatBytes, value)
	}
	return concatBytes
}
func FlatAndConcatenateLength(values [][]byte) []byte {
	flatBytes := Flat(values)
	return ConcatenateLength(flatBytes)
}

func Equals(bytes1 []byte, bytes2 []byte) bool {
	return bytes.Equal(bytes1, bytes2)
}

func Copy(src []byte, srcPos int, destPos int) []byte {
	length := destPos - srcPos
	return src[srcPos:length]
}

func CopyTo(src []byte, srcPos int, dest *[]byte, destPos int, length int) {
	for len(*dest) < destPos+length {
		*dest = append(*dest, byte(0x00))
	}
	for i := 0; i < length; i = i + 1 {
		(*dest)[destPos+i] = src[srcPos+i]
	}
}

func Random32Bytes() []byte {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		// Handle err
	}
	return token
}
