package Sha256Util

/*
 @author x.king xdotking@gmail.com
*/

import (
	"crypto/sha256"
)

func Digest(input []byte) []byte {
	sha256 := sha256.New()
	sha256.Write(input)
	sha256Digest := sha256.Sum(nil)
	return sha256Digest
}

func DoubleDigest(input []byte) []byte {
	return Digest(Digest(input))
}
