package Base58Util

/*
 @author king 409060350@qq.com
*/

import (
	"github.com/btcsuite/btcutil/base58"
)

func Encode(input []byte) string {
	return base58.Encode(input)
}
func Decode(input string) []byte {
	return base58.Decode(input)
}
