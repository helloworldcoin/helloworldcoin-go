package Ripemd160Util

import (
	"golang.org/x/crypto/ripemd160"
)

func Digest(input []byte) []byte {
	hasher := ripemd160.New()
	hasher.Write(input)
	hashBytes := hasher.Sum(nil)
	return hashBytes
}
