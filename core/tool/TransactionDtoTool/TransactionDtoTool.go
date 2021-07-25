package TransactionDtoTool

import (
	"helloworld-blockchain-go/core/tool/ScriptDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/crypto/Sha256Util"
	"helloworld-blockchain-go/dto"
)

func CalculateTransactionHash(transaction *dto.TransactionDto) string {
	bytesTransaction := BytesTransaction(transaction, false)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return ByteUtil.BytesToHexString(sha256Digest)
}

func BytesTransaction(transaction *dto.TransactionDto, omitInputScript bool) []byte {
	var bytesUnspentTransactionOutputs [][]byte
	inputs := transaction.Inputs
	for _, input := range inputs {
		bytesTransactionHash := ByteUtil.HexStringToBytes(input.TransactionHash)
		bytesTransactionOutputIndex := ByteUtil.Uint64ToBytes(input.TransactionOutputIndex)
		bytesUnspentTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesTransactionHash),
			ByteUtil.ConcatenateLength(bytesTransactionOutputIndex))
		if !omitInputScript {
			bytesInputScript := ScriptDtoTool.BytesInputScript(input.InputScript)
			bytesUnspentTransactionOutput = ByteUtil.Concatenate(bytesUnspentTransactionOutput, ByteUtil.ConcatenateLength(bytesInputScript))
		}
		bytesUnspentTransactionOutputs = append(bytesUnspentTransactionOutputs, ByteUtil.ConcatenateLength(bytesUnspentTransactionOutput))
	}

	var bytesTransactionOutputs [][]byte
	outputs := transaction.Outputs
	for _, output := range outputs {
		bytesOutputScript := ScriptDtoTool.BytesOutputScript(output.OutputScript)
		bytesValue := ByteUtil.Uint64ToBytes(output.Value)
		bytesTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesOutputScript), ByteUtil.ConcatenateLength(bytesValue))
		bytesTransactionOutputs = append(bytesTransactionOutputs, ByteUtil.ConcatenateLength(bytesTransactionOutput))
	}

	data := ByteUtil.Concatenate(ByteUtil.FlatAndConcatenateLength(bytesUnspentTransactionOutputs),
		ByteUtil.FlatAndConcatenateLength(bytesTransactionOutputs))
	return data
}

/**
 * 获取待签名数据
 */
func SignatureHashAll(transactionDto *dto.TransactionDto) string {
	bytesTransaction := BytesTransaction(transactionDto, true)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return ByteUtil.BytesToHexString(sha256Digest)
}

/**
 * 验证签名
 */
func VerifySignature(transaction *dto.TransactionDto, publicKey string, bytesSignature []byte) bool {
	message := SignatureHashAll(transaction)
	bytesMessage := ByteUtil.HexStringToBytes(message)
	return AccountUtil.VerifySignature(publicKey, bytesMessage, bytesSignature)
}
