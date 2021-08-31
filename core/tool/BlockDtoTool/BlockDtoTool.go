package BlockDtoTool

import (
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/crypto/MerkleTreeUtil"
	"helloworld-blockchain-go/crypto/Sha256Util"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/util/StringUtil"
)

func CalculateBlockHash(block *dto.BlockDto) string {
	bytesTimestamp := ByteUtil.Uint64ToBytes(block.Timestamp)
	bytesPreviousBlockHash := ByteUtil.HexStringToBytes(block.PreviousHash)
	bytesMerkleTreeRoot := ByteUtil.HexStringToBytes(CalculateBlockMerkleTreeRoot(block))
	bytesNonce := ByteUtil.HexStringToBytes(block.Nonce)

	bytes := ByteUtil.Concatenate4(bytesTimestamp, bytesPreviousBlockHash, bytesMerkleTreeRoot, bytesNonce)
	hash := Sha256Util.DoubleDigest(bytes)
	hexHash := ByteUtil.BytesToHexString(hash)
	return hexHash
}

func CalculateBlockMerkleTreeRoot(block *dto.BlockDto) string {
	transactions := block.Transactions
	var bytesTransactionHashs [][]byte
	for _, transaction := range transactions {
		transactionHash := TransactionDtoTool.CalculateTransactionHash(transaction)
		bytesTransactionHash := ByteUtil.HexStringToBytes(transactionHash)
		bytesTransactionHashs = append(bytesTransactionHashs, bytesTransactionHash)
	}
	return ByteUtil.BytesToHexString(MerkleTreeUtil.CalculateMerkleTreeRoot(bytesTransactionHashs))
}

/**
 * 简单的校验两个区块是否相等
 * 注意：这里没有严格校验,例如没有校验区块中的交易是否完全一样
 * ，所以即使这里认为两个区块相等，实际上这两个区块还是有可能不相等的。
 */
func IsBlockEquals(block1 *dto.BlockDto, block2 *dto.BlockDto) bool {
	//如果任一区块为为空，则认为两个区块不相等
	if block1 == nil || block2 == nil {
		return false
	}
	block1Hash := CalculateBlockHash(block1)
	block2Hash := CalculateBlockHash(block2)
	return StringUtil.IsEquals(block1Hash, block2Hash)
}
