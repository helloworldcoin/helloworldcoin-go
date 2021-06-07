package Sha256Util

import (
	"crypto/sha256"
)

func Digest(message []byte) []byte {
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	return bytes
}

func DoubleDigest(message []byte) []byte {
	return Digest(Digest(message))
}
