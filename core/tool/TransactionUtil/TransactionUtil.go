package TransactionUtil

import (
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/Sha256Util"

	"helloworldcoin-go/core/Model"
	"helloworldcoin-go/core/tool/ScriptTool"
)

func CalculateTransactionHash(transaction Model.Transaction) string {
	bytesTransaction := BytesTransaction(transaction, false)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return HexUtil.BytesToHexString(sha256Digest)
}

func BytesTransaction(transaction Model.Transaction, omitInputScript bool) []byte {

	bytesUnspentTransactionOutputs := [][]byte{}
	inputs := transaction.Inputs
	for _, input := range inputs {
		bytesTransactionHash := HexUtil.HexStringToBytes(input.TransactionHash)
		bytesTransactionOutputIndex := ByteUtil.Long8ToByte8(input.TransactionOutputIndex)
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
		bytesValue := ByteUtil.Long8ToByte8(output.Value)
		bytesTransactionOutput := ByteUtil.Concat(ByteUtil.ConcatLength(bytesOutputScript), ByteUtil.ConcatLength(bytesValue))
		bytesTransactionOutputs = append(bytesTransactionOutputs, ByteUtil.ConcatLength(bytesTransactionOutput))
	}

	data := ByteUtil.Concat(ByteUtil.FlatAndConcatLength(bytesUnspentTransactionOutputs),
		ByteUtil.FlatAndConcatLength(bytesTransactionOutputs))
	return data
}
