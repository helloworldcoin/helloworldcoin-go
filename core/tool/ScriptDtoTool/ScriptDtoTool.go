package ScriptDtoTool

import (
	"helloworldcoin-go/core/Model/Script/OperationCodeEnum"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/util/StringUtil"
)

//region 序列化与反序列化
func BytesInputScript(inputScript dto.InputScriptDto) []byte {
	return BytesScript(inputScript)
}
func BytesOutputScript(outputScript dto.OutputScriptDto) []byte {
	return BytesScript(outputScript)
}

/**
 * 字节型脚本：将脚本序列化，要求序列化后的脚本可以反序列化。
 */
func BytesScript(script dto.ScriptDto) []byte {
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

//endregion

func GetPublicKeyHashFromPayToPublicKeyHashOutputScript(outputScript dto.OutputScriptDto) string {
	return outputScript[3]
}

/**
 * 是否是P2PKH输入脚本
 */
func IsPayToPublicKeyHashInputScript(inputScriptDto dto.InputScriptDto) bool {
	return (len(inputScriptDto) == 4) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_PUSHDATA.Code), inputScriptDto[0])) &&
		(136 <= StringUtil.UtfCharacterCount(inputScriptDto[1]) && 144 >= StringUtil.UtfCharacterCount(inputScriptDto[1])) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_PUSHDATA.Code), inputScriptDto[2])) &&
		(66 == StringUtil.UtfCharacterCount(inputScriptDto[3]))
}

/**
 * 是否是P2PKH输出脚本
 */
func IsPayToPublicKeyHashOutputScript(outputScriptDto dto.OutputScriptDto) bool {
	return (len(outputScriptDto) == 6) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_DUP.Code), outputScriptDto[0])) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_HASH160.Code), outputScriptDto[1])) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_PUSHDATA.Code), outputScriptDto[2])) &&
		(40 == StringUtil.UtfCharacterCount(outputScriptDto[3])) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_EQUALVERIFY.Code), outputScriptDto[4])) &&
		(StringUtil.IsEquals(ByteUtil.BytesToHexString(OperationCodeEnum.OP_CHECKSIG.Code), outputScriptDto[5]))
}
