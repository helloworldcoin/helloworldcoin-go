package EncodeDecodeTool

import (
	"bytes"
	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/dto"

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

func DecodeToAccount(bytesAccount []byte) *AccountUtil.Account {
	decoder := gob.NewDecoder(bytes.NewReader(bytesAccount))
	var account AccountUtil.Account
	err := decoder.Decode(&account)
	if err != nil {
		panic(err)
	}
	return &account
}
func EncodeBlock(block *model.Block) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToBlock(bytesBlock []byte) *model.Block {
	decoder := gob.NewDecoder(bytes.NewReader(bytesBlock))
	var block *model.Block
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return block
}

func EncodeTransaction(transaction *model.Transaction) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&transaction)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToTransaction(bytesTransaction []byte) *model.Transaction {
	decoder := gob.NewDecoder(bytes.NewReader(bytesTransaction))
	var transaction *model.Transaction
	err := decoder.Decode(&transaction)
	if err != nil {
		panic(err)
	}
	return transaction
}

func EncodeTransactionOutput(transactionOutput *model.TransactionOutput) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&transactionOutput)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToTransactionOutput(bytesTransactionOutput []byte) *model.TransactionOutput {
	decoder := gob.NewDecoder(bytes.NewReader(bytesTransactionOutput))
	var transactionOutput *model.TransactionOutput
	err := decoder.Decode(&transactionOutput)
	if err != nil {
		panic(err)
	}
	return transactionOutput
}

func EncodeTransactionDto(transactionDto *dto.TransactionDto) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&transactionDto)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToTransactionDto(bytesTransactionDto []byte) *dto.TransactionDto {
	decoder := gob.NewDecoder(bytes.NewReader(bytesTransactionDto))
	var transactionDto *dto.TransactionDto
	err := decoder.Decode(&transactionDto)
	if err != nil {
		panic(err)
	}
	return transactionDto
}
func EncodeBlockDto(blockDto *dto.BlockDto) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&blockDto)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func DecodeToBlockDto(bytesBlockDto []byte) *dto.BlockDto {
	decoder := gob.NewDecoder(bytes.NewReader(bytesBlockDto))
	var blockDto *dto.BlockDto
	err := decoder.Decode(blockDto)
	if err != nil {
		panic(err)
	}
	return blockDto
}
