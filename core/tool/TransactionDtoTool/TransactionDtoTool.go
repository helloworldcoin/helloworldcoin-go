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
	var bytesUnspentTransactionOutputs [][]byte
	inputs := transaction.Inputs
	for _, input := range inputs {
		bytesTransactionHash := ByteUtil.HexStringToBytes(input.TransactionHash)
		bytesTransactionOutputIndex := ByteUtil.Uint64ToBytes(input.TransactionOutputIndex)
		bytesUnspentTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesTransactionHash),
			ByteUtil.ConcatenateLength(bytesTransactionOutputIndex))
		if !omitInputScript {
			bytesInputScript := ScriptTool.BytesScript(input.InputScript)
			bytesUnspentTransactionOutput = ByteUtil.Concatenate(bytesUnspentTransactionOutput, ByteUtil.ConcatenateLength(bytesInputScript))
		}
		bytesUnspentTransactionOutputs = append(bytesUnspentTransactionOutputs, ByteUtil.ConcatenateLength(bytesUnspentTransactionOutput))
	}

	var bytesTransactionOutputs [][]byte
	outputs := transaction.Outputs
	for _, output := range outputs {
		bytesOutputScript := ScriptTool.BytesScript(output.OutputScript)
		bytesValue := ByteUtil.Uint64ToBytes(output.Value)
		bytesTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesOutputScript), ByteUtil.ConcatenateLength(bytesValue))
		bytesTransactionOutputs = append(bytesTransactionOutputs, ByteUtil.ConcatenateLength(bytesTransactionOutput))
	}

	data := ByteUtil.Concatenate(ByteUtil.FlatAndConcatenateLength(bytesUnspentTransactionOutputs),
		ByteUtil.FlatAndConcatenateLength(bytesTransactionOutputs))
	return data
}
