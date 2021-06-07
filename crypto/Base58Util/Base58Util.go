package Base58Util

import (
	"github.com/btcsuite/btcutil/base58"
)

func Encode(input []byte) string {
	return base58.Encode(input)
}
