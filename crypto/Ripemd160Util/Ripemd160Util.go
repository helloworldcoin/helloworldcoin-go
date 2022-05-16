package Ripemd160Util

/*
 @author x.king xdotking@gmail.com
*/

import (
	"golang.org/x/crypto/ripemd160"
)

func Digest(input []byte) []byte {
	ripemd160 := ripemd160.New()
	ripemd160.Write(input)
	ripeMD160Digest := ripemd160.Sum(nil)
	return ripeMD160Digest
}
