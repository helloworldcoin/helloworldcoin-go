package ScriptTool

import (
	"helloworld-blockchain-go/core/Model/Script"
	"helloworld-blockchain-go/core/Model/Script/OperationCodeEnum"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptDtoTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
)

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
func CreateScript(inputScript Script.InputScript, outputScript Script.OutputScript) []string {
	var script Script.Script
	script = append(script, inputScript...)
	script = append(script, outputScript...)
	return script
}

/**
 * 是否是P2PKH输入脚本
 */
func IsPayToPublicKeyHashInputScript(inputScript Script.InputScript) bool {
	inputScriptDto := Model2DtoTool.InputScript2InputScriptDto(inputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashInputScript(inputScriptDto)
}

/**
 * 是否是P2PKH输出脚本
 */
func IsPayToPublicKeyHashOutputScript(outputScript Script.OutputScript) bool {
	outputScriptDto := Model2DtoTool.OutputScript2OutputScriptDto(outputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashOutputScript(outputScriptDto)
}
