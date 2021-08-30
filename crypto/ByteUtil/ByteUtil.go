package ByteUtil

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"strconv"
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

func BooleanToUtf8Bytes(booleanValue bool) []byte {

	return StringToUtf8Bytes(strconv.FormatBool(booleanValue))
}

func Utf8BytesToBoolean(bytesValue []byte) bool {
	b, _ := strconv.ParseBool(string(bytesValue))
	return b
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

func IsEquals(bytes1 []byte, bytes2 []byte) bool {
	return bytes.Equal(bytes1, bytes2)
}

func Copy(sourceBytes []byte, startPosition int, length int) []byte {
	return sourceBytes[startPosition : startPosition+length]
}

func CopyTo(sourceBytes []byte, sourceStartPosition int, length int, destinationBytes *[]byte, destinationStartPosition int) {
	for len(*destinationBytes) < destinationStartPosition+length {
		*destinationBytes = append(*destinationBytes, byte(0x00))
	}
	for i := 0; i < length; i = i + 1 {
		(*destinationBytes)[destinationStartPosition+i] = sourceBytes[sourceStartPosition+i]
	}
}

func Random32Bytes() []byte {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		// Handle err
	}
	return token
}
