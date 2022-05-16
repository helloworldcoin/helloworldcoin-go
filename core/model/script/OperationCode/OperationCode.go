package OperationCode

/*
 @author x.king xdotking@gmail.com
*/

type OperationCode struct {
	Code []byte
	Name string
}

var bytes00 []byte = []byte{0x00}
var OP_PUSHDATA = OperationCode{Code: bytes00, Name: "OP_PUSHDATA"}
var bytes01 []byte = []byte{0x01}
var OP_DUP = OperationCode{Code: bytes01, Name: "OP_DUP"}
var bytes02 []byte = []byte{0x02}
var OP_HASH160 = OperationCode{Code: bytes02, Name: "OP_HASH160"}
var bytes03 []byte = []byte{0x03}
var OP_EQUALVERIFY = OperationCode{Code: bytes03, Name: "OP_EQUALVERIFY"}
var bytes04 []byte = []byte{0x04}
var OP_CHECKSIG = OperationCode{Code: bytes04, Name: "OP_CHECKSIG"}
