package EncodeDecodeTool

import (
	"bytes"

	"encoding/gob"
	"helloworldcoin-go/crypto/AccountUtil"
)

func EncodeAccount(account *AccountUtil.Account) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&account)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToAccount(bytesAccount []byte) AccountUtil.Account {
	decoder := gob.NewDecoder(bytes.NewReader(bytesAccount))
	var account AccountUtil.Account
	err := decoder.Decode(&account)
	if err != nil {
		panic(err)
	}
	return account
}
