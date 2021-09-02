package ScriptTool

import (
	"helloworld-blockchain-go/core/Model/Script"
	"helloworld-blockchain-go/core/Model/Script/OperationCodeEnum"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

func CreatePayToPublicKeyHashOutputScript(address string) *Script.OutputScript {
	var script Script.OutputScript
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_DUP.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_HASH160.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_PUSHDATA.Code))
	publicKeyHash := AccountUtil.PublicKeyHashFromAddress(address)
	script = append(script, publicKeyHash)
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_EQUALVERIFY.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCodeEnum.OP_CHECKSIG.Code))
	return &script
}

/**
 * 构建完整脚本
 */
func CreateScript(inputScript *Script.InputScript, outputScript *Script.OutputScript) *Script.Script {
	var script Script.Script
	script = append(script, *inputScript...)
	script = append(script, *outputScript...)
	return &script
}

/**
 * 是否是P2PKH输入脚本
 */
func IsPayToPublicKeyHashInputScript(inputScript *Script.InputScript) bool {
	inputScriptDto := Model2DtoTool.InputScript2InputScriptDto(inputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashInputScript(inputScriptDto)
}

/**
 * 是否是P2PKH输出脚本
 */
func IsPayToPublicKeyHashOutputScript(outputScript *Script.OutputScript) bool {
	outputScriptDto := Model2DtoTool.OutputScript2OutputScriptDto(outputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashOutputScript(outputScriptDto)
}

//region 可视、可阅读的脚本，区块链浏览器使用
func StringInputScript(inputScript *Script.InputScript) string {
	return stringScript(inputScript)
}
func StringOutputScript(outputScript *Script.OutputScript) string {
	return stringScript(outputScript)
}
func stringScript(script *Script.Script) string {
	stringScript := ""
	for i := 0; i < len(*script); i++ {
		operationCode := (*script)[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.IsEquals(OperationCodeEnum.OP_DUP.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCodeEnum.OP_DUP.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_HASH160.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCodeEnum.OP_HASH160.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_EQUALVERIFY.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCodeEnum.OP_EQUALVERIFY.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_CHECKSIG.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCodeEnum.OP_CHECKSIG.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.IsEquals(OperationCodeEnum.OP_PUSHDATA.Code, bytesOperationCode) {
			i = i + 1
			operationData := (*script)[i]
			stringScript = StringUtil.Concatenate3(stringScript, OperationCodeEnum.OP_PUSHDATA.Name, StringUtil.BLANKSPACE)
			stringScript = StringUtil.Concatenate3(stringScript, operationData, StringUtil.BLANKSPACE)
		} else {
			panic("不能识别的指令")
		}
	}
	return stringScript
}

//endregion
