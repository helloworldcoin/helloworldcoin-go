package RandomUtil

import (
	"math/rand"
)

func Random32Bytes() []byte {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		// Handle err
	}
	return token
}
