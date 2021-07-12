package TransactionDtoTool

import (
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/Sha256Util"
	"helloworldcoin-go/dto"

	"helloworldcoin-go/core/tool/ScriptTool"
)

func CalculateTransactionHash(transaction *dto.TransactionDto) string {
	bytesTransaction := BytesTransaction(transaction, false)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return ByteUtil.BytesToHexString(sha256Digest)
}

func BytesTransaction(transaction *dto.TransactionDto, omitInputScript bool) []byte {
	bytesUnspentTransactionOutputs := [][]byte{}
	inputs := transaction.Inputs
	for _, input := range inputs {
		bytesTransactionHash := ByteUtil.HexStringToBytes(input.TransactionHash)
		bytesTransactionOutputIndex := ByteUtil.Uint64ToBytes(input.TransactionOutputIndex)
		bytesUnspentTransactionOutput := ByteUtil.Concat(ByteUtil.ConcatLength(bytesTransactionHash),
			ByteUtil.ConcatLength(bytesTransactionOutputIndex))
		if !omitInputScript {
			bytesInputScript := ScriptTool.BytesScript(input.InputScript)
			bytesUnspentTransactionOutput = ByteUtil.Concat(bytesUnspentTransactionOutput, ByteUtil.ConcatLength(bytesInputScript))
		}
		bytesUnspentTransactionOutputs = append(bytesUnspentTransactionOutputs, ByteUtil.ConcatLength(bytesUnspentTransactionOutput))
	}

	bytesTransactionOutputs := [][]byte{}
	outputs := transaction.Outputs
	for _, output := range outputs {
		bytesOutputScript := ScriptTool.BytesScript(output.OutputScript)
		bytesValue := ByteUtil.Uint64ToBytes(output.Value)
		bytesTransactionOutput := ByteUtil.Concat(ByteUtil.ConcatLength(bytesOutputScript), ByteUtil.ConcatLength(bytesValue))
		bytesTransactionOutputs = append(bytesTransactionOutputs, ByteUtil.ConcatLength(bytesTransactionOutput))
	}

	data := ByteUtil.Concat(ByteUtil.FlatAndConcatLength(bytesUnspentTransactionOutputs),
		ByteUtil.FlatAndConcatLength(bytesTransactionOutputs))
	return data
}
