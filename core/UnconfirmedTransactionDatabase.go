package core

import (
	"helloworldcoin-go/core/tool/EncodeDecodeTool"
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/util/FileUtil"
	"helloworldcoin-go/util/KvDbUtil"
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
	var transactionDtoList []dto.TransactionDto
	bytesTransactionDtos := KvDbUtil.Gets(u.getUnconfirmedTransactionDatabasePath(), from, size)
	if bytesTransactionDtos != nil {
		for e := bytesTransactionDtos.Front(); e != nil; e = e.Next() {
			transactionDto := EncodeDecodeTool.DecodeToTransactionDto(e.Value.([]byte))
			transactionDtoList = append(transactionDtoList, *transactionDto)
		}
	}
	return transactionDtoList
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
	return HexUtil.HexStringToBytes(transactionHash)
}
