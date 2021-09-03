package main

import (
	"fmt"
	"helloworld-blockchain-go/crypto/ByteUtil"
)

func main() {
	bytes := ByteUtil.HexStringToBytes("00010203")
	fmt.Printf("%p", bytes)
	fmt.Println("")
	fmt.Printf("%p", ByteUtil.Copy(bytes, 0, 1))
	fmt.Println("")
	fmt.Printf("%p", ByteUtil.Copy(bytes, 1, 1))
	fmt.Println("")
	fmt.Printf("%p", ByteUtil.Copy(bytes, 2, 1))
}
