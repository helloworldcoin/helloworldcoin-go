package ScriptTool

import (
	"helloworldcoin-go/core/model/OperationCodeEnum"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/HexUtil"
)

func BytesScript(script []string) []byte {
	bytesScript := []byte{}
	for i := 0; i < len(script); i++ {
		operationCode := script[i]
		bytesOperationCode := HexUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCodeEnum.OP_DUP.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_HASH160.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_EQUALVERIFY.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_CHECKSIG.Code, bytesOperationCode) {
			bytesScript = ByteUtil.Concat(bytesScript, ByteUtil.ConcatLength(bytesOperationCode))

		} else if ByteUtil.Equals(OperationCodeEnum.OP_PUSHDATA.Code, bytesOperationCode) {
			i = i + 1
			operationData := script[i]
			bytesOperationData := HexUtil.HexStringToBytes(operationData)
			bytesScript = ByteUtil.Concat(bytesScript, ByteUtil.ConcatLength(bytesOperationCode), ByteUtil.ConcatLength(bytesOperationData))

		} else {

		}

	}

	return bytesScript
}
func GetPublicKeyHashByPayToPublicKeyHashOutputScript(outputScript []string) string {
	return outputScript[3]
}
