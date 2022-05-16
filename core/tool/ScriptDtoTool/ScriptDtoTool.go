package ScriptDtoTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/model/script/OperationCode"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/StringUtil"
)

//region Serialization and Deserialization
func InputScript2Bytes(inputScript *dto.InputScriptDto) []byte {
	return Script2Bytes(inputScript)
}
func OutputScript2Bytes(outputScript *dto.OutputScriptDto) []byte {
	return Script2Bytes(outputScript)
}
func Script2Bytes(script *dto.ScriptDto) []byte {
	var bytesScript []byte
	for i := 0; i < len(*script); i++ {
		operationCode := (*script)[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCode.OP_DUP.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCode.OP_HASH160.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCode.OP_EQUALVERIFY.Code, bytesOperationCode) ||
			ByteUtil.Equals(OperationCode.OP_CHECKSIG.Code, bytesOperationCode) {
			bytesScript = ByteUtil.Concatenate(bytesScript, ByteUtil.ConcatenateLength(bytesOperationCode))
		} else if ByteUtil.Equals(OperationCode.OP_PUSHDATA.Code, bytesOperationCode) {
			i = i + 1
			operationData := (*script)[i]
			bytesOperationData := ByteUtil.HexStringToBytes(operationData)
			bytesScript = ByteUtil.Concatenate3(bytesScript, ByteUtil.ConcatenateLength(bytesOperationCode), ByteUtil.ConcatenateLength(bytesOperationData))
		} else {
			panic("Unrecognized OperationCode.")
		}
	}
	return bytesScript
}

//endregion

func CreatePayToPublicKeyHashInputScript(sign string, publicKey string) *dto.InputScriptDto {
	var script dto.InputScriptDto

	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	script = append(script, sign)
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	script = append(script, publicKey)
	return &script
}
func CreatePayToPublicKeyHashOutputScript(address string) *dto.OutputScriptDto {
	var script dto.OutputScriptDto
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_DUP.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_HASH160.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code))
	publicKeyHash := AccountUtil.PublicKeyHashFromAddress(address)
	script = append(script, publicKeyHash)
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_EQUALVERIFY.Code))
	script = append(script, ByteUtil.BytesToHexString(OperationCode.OP_CHECKSIG.Code))
	return &script
}
func IsPayToPublicKeyHashInputScript(inputScriptDto *dto.InputScriptDto) bool {
	return (len(*inputScriptDto) == 4) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code), (*inputScriptDto)[0])) &&
		(136 <= StringUtil.Length((*inputScriptDto)[1]) && 144 >= StringUtil.Length((*inputScriptDto)[1])) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code), (*inputScriptDto)[2])) &&
		(66 == StringUtil.Length((*inputScriptDto)[3]))
}
func IsPayToPublicKeyHashOutputScript(outputScriptDto *dto.OutputScriptDto) bool {
	return (len(*outputScriptDto) == 6) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_DUP.Code), (*outputScriptDto)[0])) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_HASH160.Code), (*outputScriptDto)[1])) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_PUSHDATA.Code), (*outputScriptDto)[2])) &&
		(40 == StringUtil.Length((*outputScriptDto)[3])) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_EQUALVERIFY.Code), (*outputScriptDto)[4])) &&
		(StringUtil.Equals(ByteUtil.BytesToHexString(OperationCode.OP_CHECKSIG.Code), (*outputScriptDto)[5]))
}
func GetPublicKeyHashFromPayToPublicKeyHashOutputScript(outputScript *dto.OutputScriptDto) string {
	return (*outputScript)[3]
}
