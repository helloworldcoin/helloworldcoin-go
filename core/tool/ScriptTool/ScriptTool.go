package ScriptTool

import (
	"helloworldcoin-go/core/model/OperationCodeEnum"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/ByteUtil"
)

func BytesScript(script []string) []byte {
	var bytesScript []byte
	for i := 0; i < len(script); i++ {
		operationCode := script[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCodeEnum.OP_DUP.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_HASH160.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_EQUALVERIFY.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCodeEnum.OP_CHECKSIG.Code, bytesOperationCode) {
			bytesScript = ByteUtil.Concatenate(bytesScript, ByteUtil.ConcatenateLength(bytesOperationCode))
		} else if ByteUtil.Equals(OperationCodeEnum.OP_PUSHDATA.Code, bytesOperationCode) {
			i = i + 1
			operationData := script[i]
			bytesOperationData := ByteUtil.HexStringToBytes(operationData)
			bytesScript = ByteUtil.Concatenate3(bytesScript, ByteUtil.ConcatenateLength(bytesOperationCode), ByteUtil.ConcatenateLength(bytesOperationData))
		} else {

		}

	}

	return bytesScript
}
func GetPublicKeyHashByPayToPublicKeyHashOutputScript(outputScript []string) string {
	return outputScript[3]
}

func CreatePayToPublicKeyHashOutputScript(address string) []string {
	var script []string
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_DUP.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_HASH160.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_PUSHDATA.Code))
	publicKeyHash := AccountUtil.PublicKeyHashFromAddress(address)
	script = append(script, publicKeyHash)
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_EQUALVERIFY.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_CHECKSIG.Code))
	return script
}
