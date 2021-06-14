package AccountUtil

import (
	"fmt"
	"helloworldcoin-go/crypto/Base58Util"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/Ripemd160Util"
	"helloworldcoin-go/crypto/Sha256Util"

	"github.com/btcsuite/btcd/btcec"
)

func RandomAccount() *Account {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PubKey().SerializeCompressed()

	stringPrivateKey := encodePrivateKey(privateKey)
	stringPublicKey := encodePublicKey(publicKey)
	stringPublicKeyHash := publicKeyHashFromStringPublicKey(stringPublicKey)
	stringAddress := addressFromStringPublicKey(stringPublicKey)
	account := Account{PrivateKey: stringPrivateKey, PublicKey: stringPublicKey, PublicKeyHash: stringPublicKeyHash, Address: stringAddress}
	return &account
}
func encodePrivateKey(privateKey *btcec.PrivateKey) string {
	return HexUtil.BytesToHexString(privateKey.D.Bytes())
}
func encodePublicKey(publicKey []byte) string {
	return HexUtil.BytesToHexString(publicKey)
}
func publicKeyHashFromStringPublicKey(stringPublicKey string) string {
	publicKeyHash := Ripemd160Util.Digest(Sha256Util.Digest(HexUtil.HexStringToBytes(stringPublicKey)))
	return HexUtil.BytesToHexString(publicKeyHash)
}
func addressFromStringPublicKey(stringPublicKey string) string {
	bytesPublicKey := HexUtil.HexStringToBytes(stringPublicKey)
	return base58AddressFromPublicKey(bytesPublicKey)
}
func base58AddressFromPublicKey(bytesPublicKey []byte) string {
	publicKeyHash := publicKeyHashFromPublicKey(bytesPublicKey)
	return base58AddressFromBytesPublicKeyHash(publicKeyHash)
}
func publicKeyHashFromPublicKey(publicKey []byte) []byte {
	return Ripemd160Util.Digest(Sha256Util.Digest(publicKey))
}
func base58AddressFromBytesPublicKeyHash(bytesPublicKeyHash []byte) string {
	bytesCheckCode := ByteUtil.Copy(Sha256Util.DoubleDigest(append([]byte{0x00}, bytesPublicKeyHash...)), 0, 4)
	bytesAddress := []byte{}
	bytesAddress = append([]byte{0x00}, bytesPublicKeyHash...)
	bytesAddress = append(bytesAddress, bytesCheckCode...)
	base58Address := Base58Util.Encode(bytesAddress)
	return base58Address
}
func Signature(stringPrivateKey string, message string) string {
	privateKey := privateKeyFrom(stringPrivateKey)
	bytesMessage := HexUtil.HexStringToBytes(message)
	bytesSignature := signature0(privateKey, bytesMessage)
	return HexUtil.BytesToHexString(bytesSignature)
}
func privateKeyFrom(stringPrivateKey string) *btcec.PrivateKey {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), HexUtil.HexStringToBytes(stringPrivateKey))
	return privKey
}
func signature0(privateKey *btcec.PrivateKey, message []byte) []byte {
	signature, _ := privateKey.Sign(message)
	return signature.Serialize()
}
func VerifySignature(stringPublicKey string, stringMessage string, stringSignature string) bool {
	publicKey := publicKeyFrom(stringPublicKey)
	bytesMessage := HexUtil.HexStringToBytes(stringMessage)
	signature, _ := btcec.ParseDERSignature(HexUtil.HexStringToBytes(stringSignature), btcec.S256())
	return signature.Verify(bytesMessage, publicKey)
}
func publicKeyFrom(stringPublicKey string) *btcec.PublicKey {
	bytesPublicKey := HexUtil.HexStringToBytes(stringPublicKey)
	publicKey, _ := btcec.ParsePubKey(bytesPublicKey, btcec.S256())
	return publicKey
}

func AddressFromStringPublicKeyHash(stringPublicKeyHash string) string {
	bytesPublicKeyHash := HexUtil.HexStringToBytes(stringPublicKeyHash)
	return base58AddressFromBytesPublicKeyHash(bytesPublicKeyHash)
}
func PublicKeyHashFromStringAddress(stringAddress string) string {
	bytesAddress := Base58Util.Decode(stringAddress)
	fmt.Println(bytesAddress)
	var bytesPublicKeyHash []byte
	ByteUtil.CopyTo(bytesAddress, 1, bytesPublicKeyHash, 0, 20)
	return HexUtil.BytesToHexString(bytesPublicKeyHash)
}
