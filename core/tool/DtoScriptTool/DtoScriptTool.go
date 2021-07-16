package DtoScriptTool

import (
	"helloworldcoin-go/core/Model/script/OperationCodeEnum"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/util/StringUtil"
)

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
