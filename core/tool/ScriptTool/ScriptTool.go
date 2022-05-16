package ScriptTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/script/OperationCode"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/ScriptDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/StringUtil"
)

func InputScript2String(inputScript *model.InputScript) string {
	return Script2String(inputScript)
}
func OutputScript2String(outputScript *model.OutputScript) string {
	return Script2String(outputScript)
}
func Script2String(script *model.Script) string {
	stringScript := ""
	for i := 0; i < len(*script); i++ {
		operationCode := (*script)[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCode.OP_DUP.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCode.OP_DUP.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.Equals(OperationCode.OP_HASH160.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCode.OP_HASH160.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.Equals(OperationCode.OP_EQUALVERIFY.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCode.OP_EQUALVERIFY.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.Equals(OperationCode.OP_CHECKSIG.Code, bytesOperationCode) {
			stringScript = StringUtil.Concatenate3(stringScript, OperationCode.OP_CHECKSIG.Name, StringUtil.BLANKSPACE)
		} else if ByteUtil.Equals(OperationCode.OP_PUSHDATA.Code, bytesOperationCode) {
			i = i + 1
			operationData := (*script)[i]
			stringScript = StringUtil.Concatenate3(stringScript, OperationCode.OP_PUSHDATA.Name, StringUtil.BLANKSPACE)
			stringScript = StringUtil.Concatenate3(stringScript, operationData, StringUtil.BLANKSPACE)
		} else {
			panic("Unrecognized OperationCode.")
		}
	}
	return stringScript
}

func CreateScript(inputScript *model.InputScript, outputScript *model.OutputScript) *model.Script {
	var script model.Script
	script = append(script, *inputScript...)
	script = append(script, *outputScript...)
	return &script
}
func CreatePayToPublicKeyHashInputScript(sign string, publicKey string) *model.InputScript {
	var script model.InputScript
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	script = append(script, sign)
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	script = append(script, publicKey)
	return &script
}
func CreatePayToPublicKeyHashOutputScript(address string) *model.OutputScript {
	var script model.OutputScript
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_DUP.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_HASH160.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	publicKeyHash := AccountUtil.PublicKeyHashFromAddress(address)
	script = append(script, publicKeyHash)
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_EQUALVERIFY.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_CHECKSIG.Code))
	return &script
}

func IsPayToPublicKeyHashInputScript(inputScript *model.InputScript) bool {
	inputScriptDto := Model2DtoTool.InputScript2InputScriptDto(inputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashInputScript(inputScriptDto)
}
func IsPayToPublicKeyHashOutputScript(outputScript *model.OutputScript) bool {
	outputScriptDto := Model2DtoTool.OutputScript2OutputScriptDto(outputScript)
	return ScriptDtoTool.IsPayToPublicKeyHashOutputScript(outputScriptDto)
}
