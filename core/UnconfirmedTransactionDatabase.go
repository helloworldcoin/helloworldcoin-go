package core

import (
	"helloworld-blockchain-go/core/tool/EncodeDecodeTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
)

const UNCONFIRMED_TRANSACTION_DATABASE_NAME = "UnconfirmedTransactionDatabase"

type UnconfirmedTransactionDatabase struct {
	CoreConfiguration *CoreConfiguration
}

func (u *UnconfirmedTransactionDatabase) InsertTransaction(transactionDto *dto.TransactionDto) {
	transactionHash := TransactionDtoTool.CalculateTransactionHash(transactionDto)
	KvDbUtil.Put(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash), EncodeDecodeTool.EncodeTransactionDto(transactionDto))

}

func (u *UnconfirmedTransactionDatabase) SelectTransactions(from uint64, size uint64) []dto.TransactionDto {
	var transactionDtos []dto.TransactionDto
	bytesTransactionDtos := KvDbUtil.Gets(u.getUnconfirmedTransactionDatabasePath(), from, size)
	if bytesTransactionDtos != nil {
		for e := bytesTransactionDtos.Front(); e != nil; e = e.Next() {
			transactionDto := EncodeDecodeTool.DecodeToTransactionDto(e.Value.([]byte))
			transactionDtos = append(transactionDtos, *transactionDto)
		}
	}
	return transactionDtos
}

func (u *UnconfirmedTransactionDatabase) DeleteByTransactionHash(transactionHash string) {
	KvDbUtil.Delete(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash))
}

func (u *UnconfirmedTransactionDatabase) SelectTransactionByTransactionHash(transactionHash string) *dto.TransactionDto {
	byteTransactionDto := KvDbUtil.Get(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash))
	if byteTransactionDto == nil {
		return nil
	}
	return EncodeDecodeTool.DecodeToTransactionDto(byteTransactionDto)
}

func (u *UnconfirmedTransactionDatabase) getUnconfirmedTransactionDatabasePath() string {
	return FileUtil.NewPath(u.CoreConfiguration.CorePath, UNCONFIRMED_TRANSACTION_DATABASE_NAME)
}

func (u *UnconfirmedTransactionDatabase) getKey(transactionHash string) []byte {
	return ByteUtil.HexStringToBytes(transactionHash)
}
