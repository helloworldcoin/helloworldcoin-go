package Sha256Util

/*
 @author king 409060350@qq.com
*/

import (
	"crypto/sha256"
)

func Digest(message []byte) []byte {
	sha256 := sha256.New()
	sha256.Write(message)
	sha256Digest := sha256.Sum(nil)
	return sha256Digest
}

func DoubleDigest(message []byte) []byte {
	return Digest(Digest(message))
}
