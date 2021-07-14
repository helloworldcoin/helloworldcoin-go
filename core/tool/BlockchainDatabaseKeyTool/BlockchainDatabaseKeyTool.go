package BlockchainDatabaseKeyTool

import (
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/util/StringUtil"
)

const (
	BLOCKCHAIN_HEIGHT_KEY                                                  = "A"
	BLOCKCHAIN_TRANSACTION_HEIGHT_KEY                                      = "B"
	BLOCKCHAIN_TRANSACTION_OUTPUT_HEIGHT_KEY                               = "C"
	HASH_PREFIX_FLAG                                                       = "D"
	BLOCK_HEIGHT_TO_BLOCK_PREFIX_FLAG                                      = "E"
	BLOCK_HASH_TO_BLOCK_HEIGHT_PREFIX_FLAG                                 = "F"
	TRANSACTION_HEIGHT_TO_TRANSACTION_PREFIX_FLAG                          = "G"
	TRANSACTION_HASH_TO_TRANSACTION_HEIGHT_PREFIX_FLAG                     = "H"
	TRANSACTION_OUTPUT_HEIGHT_TO_TRANSACTION_OUTPUT_PREFIX_FLAG            = "I"
	TRANSACTION_OUTPUT_ID_TO_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG         = "J"
	TRANSACTION_OUTPUT_ID_TO_UNSPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG = "K"
	TRANSACTION_OUTPUT_ID_TO_SPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG   = "L"
	TRANSACTION_OUTPUT_ID_TO_SOURCE_TRANSACTION_HEIGHT_PREFIX_FLAG         = "M"
	TRANSACTION_OUTPUT_ID_TO_DESTINATION_TRANSACTION_HEIGHT_PREFIX_FLAG    = "N"
	ADDRESS_PREFIX_FLAG                                                    = "O"
	ADDRESS_TO_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG                       = "P"
	ADDRESS_TO_UNSPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG               = "Q"
	ADDRESS_TO_SPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG                 = "R"
	END_FLAG                                                               = "#"
	VERTICAL_LINE_FLAG                                                     = "|"
)

//拼装数据库Key的值
func BuildBlockchainHeightKey() []byte {
	stringKey := BLOCKCHAIN_HEIGHT_KEY + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildHashKey(hash string) []byte {
	stringKey := HASH_PREFIX_FLAG + hash + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildAddressKey(address string) []byte {
	stringKey := ADDRESS_PREFIX_FLAG + address + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildBlockHeightToBlockKey(blockHeight uint64) []byte {
	stringKey := BLOCK_HEIGHT_TO_BLOCK_PREFIX_FLAG + StringUtil.ValueOfUint64(blockHeight) + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildBlockHashToBlockHeightKey(blockHash string) []byte {
	stringKey := BLOCK_HASH_TO_BLOCK_HEIGHT_PREFIX_FLAG + blockHash + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionHashToTransactionHeightKey(transactionHash string) []byte {
	stringKey := TRANSACTION_HASH_TO_TRANSACTION_HEIGHT_PREFIX_FLAG + transactionHash + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputHeightToTransactionOutputKey(transactionOutputHeight uint64) []byte {
	stringKey := TRANSACTION_OUTPUT_HEIGHT_TO_TRANSACTION_OUTPUT_PREFIX_FLAG + StringUtil.ValueOfUint64(transactionOutputHeight) + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputIdToTransactionOutputHeightKey(transactionHash string, transactionOutputIndex uint64) []byte {
	transactionOutputId := BuildTransactionOutputId(transactionHash, transactionOutputIndex)
	stringKey := TRANSACTION_OUTPUT_ID_TO_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + transactionOutputId + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputIdToUnspentTransactionOutputHeightKey(transactionHash string, transactionOutputIndex uint64) []byte {
	transactionOutputId := BuildTransactionOutputId(transactionHash, transactionOutputIndex)
	stringKey := TRANSACTION_OUTPUT_ID_TO_UNSPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + transactionOutputId + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputIdToSourceTransactionHeightKey(transactionHash string, transactionOutputIndex uint64) []byte {
	transactionOutputId := BuildTransactionOutputId(transactionHash, transactionOutputIndex)
	stringKey := TRANSACTION_OUTPUT_ID_TO_SOURCE_TRANSACTION_HEIGHT_PREFIX_FLAG + transactionOutputId + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputIdToDestinationTransactionHeightKey(transactionHash string, transactionOutputIndex uint64) []byte {
	transactionOutputId := BuildTransactionOutputId(transactionHash, transactionOutputIndex)
	stringKey := TRANSACTION_OUTPUT_ID_TO_DESTINATION_TRANSACTION_HEIGHT_PREFIX_FLAG + transactionOutputId + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildAddressToTransactionOutputHeightKey(address string) []byte {
	stringKey := ADDRESS_TO_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + address + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildAddressToUnspentTransactionOutputHeightKey(address string) []byte {
	stringKey := ADDRESS_TO_UNSPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + address + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildAddressToSpentTransactionOutputHeightKey(address string) []byte {
	stringKey := ADDRESS_TO_SPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + address + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildBlockchainTransactionHeightKey() []byte {
	stringKey := BLOCKCHAIN_TRANSACTION_HEIGHT_KEY + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildBlockchainTransactionOutputHeightKey() []byte {
	stringKey := BLOCKCHAIN_TRANSACTION_OUTPUT_HEIGHT_KEY + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionHeightToTransactionKey(transactionHeight uint64) []byte {
	stringKey := TRANSACTION_HEIGHT_TO_TRANSACTION_PREFIX_FLAG + StringUtil.ValueOfUint64(transactionHeight) + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}

func BuildTransactionOutputIdToSpentTransactionOutputHeightKey(transactionHash string, transactionOutputIndex uint64) []byte {
	transactionOutputId := BuildTransactionOutputId(transactionHash, transactionOutputIndex)
	stringKey := TRANSACTION_OUTPUT_ID_TO_SPENT_TRANSACTION_OUTPUT_HEIGHT_PREFIX_FLAG + transactionOutputId + END_FLAG
	return ByteUtil.StringToUtf8Bytes(stringKey)
}
func BuildTransactionOutputId(transactionHash string, transactionOutputIndex uint64) string {
	transactionOutputId := StringUtil.Concatenate3(transactionHash, VERTICAL_LINE_FLAG, StringUtil.ValueOfUint64(transactionOutputIndex))
	return transactionOutputId
}
