package ScriptTool

import (
	"helloworldcoin-go/core/Model/script"
	"helloworldcoin-go/core/Model/script/OperationCodeEnum"
	"helloworldcoin-go/core/tool/DtoScriptTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/ByteUtil"
)

func BytesScript(script []string) []byte {
	var bytesScript []byte
	for i := 0; i < len(script); i++ {
		operationCode := script[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.IsEquals(OperationCodeEnum.OP_DUP.Code, bytesOperationCode) ||
			ByteUtil.IsEquals(OperationCodeEnum.OP_HASH160.Code, bytesOperationCode) ||
			ByteUtil.IsEquals(OperationCodeEnum.OP_EQUALVERIFY.Code, bytesOperationCode) ||
			ByteUtil.IsEquals(OperationCodeEnum.OP_CHECKSIG.Code, bytesOperationCode) {
			bytesScript = ByteUtil.Concatenate(bytesScript, ByteUtil.ConcatenateLength(bytesOperationCode))
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_PUSHDATA.Code, bytesOperationCode) {
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

/**
 * 构建完整脚本
 */
func CreateScript(inputScript []string, outputScript []string) []string {
	var script []string
	script = append(script, inputScript...)
	script = append(script, outputScript...)
	return script
}

/**
 * 是否是P2PKH输入脚本
 */
func IsPayToPublicKeyHashInputScript(inputScript script.InputScript) bool {
	inputScriptDto := Model2DtoTool.InputScript2InputScriptDto(inputScript)
	return DtoScriptTool.IsPayToPublicKeyHashInputScript(inputScriptDto)
}

/**
 * 是否是P2PKH输出脚本
 */
func IsPayToPublicKeyHashOutputScript(outputScript script.OutputScript) bool {
	outputScriptDto := Model2DtoTool.OutputScript2OutputScriptDto(outputScript)
	return DtoScriptTool.IsPayToPublicKeyHashOutputScript(outputScriptDto)
}
