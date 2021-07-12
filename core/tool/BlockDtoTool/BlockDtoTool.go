package BlockDtoTool

import (
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/MerkleTreeUtil"
	"helloworldcoin-go/crypto/Sha256Util"
	"helloworldcoin-go/dto"
)

func CalculateBlockHash(block *dto.BlockDto) string {
	bytesTimestamp := ByteUtil.Uint64ToBytes(block.Timestamp)
	bytesPreviousBlockHash := ByteUtil.HexStringToBytes(block.PreviousHash)
	bytesMerkleTreeRoot := ByteUtil.HexStringToBytes(CalculateBlockMerkleTreeRoot(block))
	bytesNonce := ByteUtil.HexStringToBytes(block.Nonce)

	bytes := ByteUtil.Concat(bytesTimestamp, bytesPreviousBlockHash, bytesMerkleTreeRoot, bytesNonce)
	hash := Sha256Util.DoubleDigest(bytes)
	hexHash := ByteUtil.BytesToHexString(hash)
	return hexHash
}

func CalculateBlockMerkleTreeRoot(block *dto.BlockDto) string {
	transactions := block.Transactions
	var bytesTransactionHashs [][]byte
	for _, transaction := range transactions {
		transactionHash := TransactionDtoTool.CalculateTransactionHash(&transaction)
		bytesTransactionHash := ByteUtil.HexStringToBytes(transactionHash)
		bytesTransactionHashs = append(bytesTransactionHashs, bytesTransactionHash)
	}
	return ByteUtil.BytesToHexString(MerkleTreeUtil.CalculateMerkleTreeRoot(bytesTransactionHashs))
}
