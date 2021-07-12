package Sha256Util

import (
	"crypto/sha256"
)

func Digest(message []byte) []byte {
	sha256 := sha256.New()
	sha256.Write(message)
	bytes := sha256.Sum(nil)
	return bytes
}

func DoubleDigest(message []byte) []byte {
	return Digest(Digest(message))
}
