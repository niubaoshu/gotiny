package gotiny

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type aesConfigStruct struct {
	block     *cipher.Block
	gcm       cipher.AEAD
	nonceSize int
}

func newAESconfig(key []byte) *aesConfigStruct {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	return &aesConfigStruct{
		block:     &block,
		gcm:       gcm,
		nonceSize: gcm.NonceSize(),
	}
}

// create AES-256 config for encryption encoding key must be 32 byte
func NewAES256config(key [32]byte) *aesConfigStruct {
	return newAESconfig(key[:])
}

// create AES-192 config for encryption encoding key must be 24 byte
func NewAES192config(key [24]byte) *aesConfigStruct {
	return newAESconfig(key[:])
}

// create AES-128 config for encryption encoding key must be 16 byte
func NewAES128config(key [16]byte) *aesConfigStruct {
	return newAESconfig(key[:])
}

// Generate a random nonce
func (aesConfig *aesConfigStruct) generatorNonce() []byte {
	nonce := make([]byte, aesConfig.nonceSize)
	rand.Read(nonce)

	return nonce
}

// Encrypt the plaintext
func (aesConfig *aesConfigStruct) Encrypt(plaintext []byte) []byte {
	nonce := aesConfig.generatorNonce()
	return aesConfig.gcm.Seal(nonce, nonce, plaintext, nil)
}

// Decrypt the cipherData
func (aesConfig *aesConfigStruct) Decrypt(cipherData []byte) []byte {

	decrypted, _ := aesConfig.gcm.Open(nil,
		cipherData[:aesConfig.nonceSize],
		cipherData[aesConfig.nonceSize:],
		nil)
	return decrypted
}
