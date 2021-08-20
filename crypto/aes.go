package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

func AesCBCBase64EncryptString(text, key, iv string) (string, error) {

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Println("错误 -" + err.Error())
		return "", err
	}

	blockSize := block.BlockSize()
	originData := pad([]byte(text), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	crypt := make([]byte, len(originData))
	blockMode.CryptBlocks(crypt, originData)

	return base64.StdEncoding.EncodeToString(crypt), nil
}

func AesCBCBase64Encrypt(text, key, iv []byte) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" + err.Error())
		return "", err
	}

	blockSize := block.BlockSize()
	originData := pad(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypt := make([]byte, len(originData))
	blockMode.CryptBlocks(crypt, originData)

	return base64.StdEncoding.EncodeToString(crypt), nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	plaintext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, plaintext...)
}

func AesCBCBase64Decrypt(text string, key, iv []byte) (string, error) {
	decodeData, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)

	return string(unpacked(originData)), nil
}

func AesCBCBase64DecryptString(text, key, iv string) (string, error) {
	decodeData, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)

	return string(unpacked(originData)), nil
}

func unpacked(ciphertext []byte) []byte {
	length := len(ciphertext)
	unpacking := int(ciphertext[length-1])
	return ciphertext[:(length - unpacking)]
}
