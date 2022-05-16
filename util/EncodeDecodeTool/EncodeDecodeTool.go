package EncodeDecodeTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/util/ByteUtil"
	"helloworldcoin-go/util/JsonUtil"
)

func Encode(object interface{}) []byte {
	return ByteUtil.StringToUtf8Bytes(JsonUtil.ToString(object))
}
func Decode(bytesObject []byte, object interface{}) interface{} {
	return JsonUtil.ToObject(ByteUtil.Utf8BytesToString(bytesObject), object)
}
