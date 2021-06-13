package BlockDtoTool

import (
	"helloworldcoin-go/core/tool/TransactionDtoTool"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/MerkleTreeUtil"
	"helloworldcoin-go/crypto/Sha256Util"
	"helloworldcoin-go/dto"
)

func CalculateBlockHash(block *dto.BlockDto) string {
	bytesTimestamp := ByteUtil.Long8ToByte8(block.Timestamp)
	bytesPreviousBlockHash := HexUtil.HexStringToBytes(block.PreviousHash)
	bytesMerkleTreeRoot := HexUtil.HexStringToBytes(CalculateBlockMerkleTreeRoot(block))
	bytesNonce := HexUtil.HexStringToBytes(block.Nonce)

	bytes := ByteUtil.Concat(bytesTimestamp, bytesPreviousBlockHash, bytesMerkleTreeRoot, bytesNonce)
	hash := Sha256Util.DoubleDigest(bytes)
	hexHash := HexUtil.BytesToHexString(hash)
	return hexHash
}

func CalculateBlockMerkleTreeRoot(block *dto.BlockDto) string {
	transactions := block.Transactions
	var bytesTransactionHashs [][]byte
	for _, transaction := range transactions {
		transactionHash := TransactionDtoTool.CalculateTransactionHash(&transaction)
		bytesTransactionHash := HexUtil.HexStringToBytes(transactionHash)
		bytesTransactionHashs = append(bytesTransactionHashs, bytesTransactionHash)
	}
	return HexUtil.BytesToHexString(MerkleTreeUtil.CalculateMerkleTreeRoot(bytesTransactionHashs))
}
