package AccountUtil

import (
	"helloworldcoin-go/crypto/Base58Util"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/crypto/HexUtil"
	"helloworldcoin-go/crypto/Ripemd160Util"
	"helloworldcoin-go/crypto/Sha256Util"

	"github.com/btcsuite/btcd/btcec"
)

func RandomAccount() *Account {
	privateKey0, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	publicKey0 := privateKey0.PubKey().SerializeCompressed()

	privateKey := encodePrivateKey0(privateKey0)
	publicKey := encodePublicKey0(publicKey0)
	publicKeyHash := PublicKeyHashFromPublicKey(publicKey)
	address := AddressFromPublicKey(publicKey)
	account := Account{PrivateKey: privateKey, PublicKey: publicKey, PublicKeyHash: publicKeyHash, Address: address}
	return &account
}
func PublicKeyHashFromPublicKey(publicKey string) string {
	publicKeyHash := Ripemd160Util.Digest(Sha256Util.Digest(HexUtil.HexStringToBytes(publicKey)))
	return HexUtil.BytesToHexString(publicKeyHash)
}
func AddressFromPublicKey(publicKey string) string {
	bytesPublicKey := HexUtil.HexStringToBytes(publicKey)
	return base58AddressFromPublicKey0(bytesPublicKey)
}
func Signature(privateKey string, message string) string {
	privateKey0 := decodePrivateKey0(privateKey)
	bytesMessage := HexUtil.HexStringToBytes(message)
	bytesSignature := signature0(privateKey0, bytesMessage)
	return HexUtil.BytesToHexString(bytesSignature)
}
func VerifySignature(publicKey string, message string, signature string) bool {
	publicKey0 := decodePublicKey0(publicKey)
	bytesMessage := HexUtil.HexStringToBytes(message)
	signature0, _ := btcec.ParseDERSignature(HexUtil.HexStringToBytes(signature), btcec.S256())
	return signature0.Verify(bytesMessage, publicKey0)
}
func AddressFromPublicKeyHash(publicKeyHash string) string {
	bytesPublicKeyHash := HexUtil.HexStringToBytes(publicKeyHash)
	return base58AddressFromPublicKeyHash0(bytesPublicKeyHash)
}
func PublicKeyHashFromAddress(address string) string {
	bytesAddress := Base58Util.Decode(address)
	var bytesPublicKeyHash []byte
	ByteUtil.CopyTo(bytesAddress, 1, &bytesPublicKeyHash, 0, 20)
	return HexUtil.BytesToHexString(bytesPublicKeyHash)
}

func encodePrivateKey0(privateKey0 *btcec.PrivateKey) string {
	return HexUtil.BytesToHexString(privateKey0.D.Bytes())
}
func encodePublicKey0(publicKey []byte) string {
	return HexUtil.BytesToHexString(publicKey)
}
func decodePublicKey0(stringPublicKey string) *btcec.PublicKey {
	bytesPublicKey := HexUtil.HexStringToBytes(stringPublicKey)
	publicKey, _ := btcec.ParsePubKey(bytesPublicKey, btcec.S256())
	return publicKey
}
func decodePrivateKey0(privateKey string) *btcec.PrivateKey {
	privateKey0, _ := btcec.PrivKeyFromBytes(btcec.S256(), HexUtil.HexStringToBytes(privateKey))
	return privateKey0
}
func signature0(privateKey *btcec.PrivateKey, message []byte) []byte {
	signature, _ := privateKey.Sign(message)
	return signature.Serialize()
}
func base58AddressFromPublicKey0(bytesPublicKey []byte) string {
	publicKeyHash := publicKeyHashFromPublicKey0(bytesPublicKey)
	return base58AddressFromPublicKeyHash0(publicKeyHash)
}
func publicKeyHashFromPublicKey0(publicKey []byte) []byte {
	return Ripemd160Util.Digest(Sha256Util.Digest(publicKey))
}
func base58AddressFromPublicKeyHash0(bytesPublicKeyHash []byte) string {
	bytesCheckCode := ByteUtil.Copy(Sha256Util.DoubleDigest(append([]byte{0x00}, bytesPublicKeyHash...)), 0, 4)
	bytesAddress := []byte{}
	bytesAddress = append([]byte{0x00}, bytesPublicKeyHash...)
	bytesAddress = append(bytesAddress, bytesCheckCode...)
	base58Address := Base58Util.Encode(bytesAddress)
	return base58Address
}
