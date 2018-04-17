package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

// AES256EncryptWithCBC returns an encrypted text using the given key
// by way of cipher block chaining and AES 256
func AES256EncryptWithCBC(key, iv, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	msg := pad(text)
	ciphertext := make([]byte, aes.BlockSize+len(msg))

	// Use CBC https://golang.org/src/crypto/cipher/cbc.go?s=1152:1202#L33
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], []byte(msg))
	return ciphertext, nil
}

// AES256DecryptWithCBC returns a plaintext using the given key by
// way of cipher block chaining and AES 256
func AES256DecryptWithCBC(key, iv []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if (len(text) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	msg := []byte(text)[aes.BlockSize:]

	cfb := cipher.NewCBCDecrypter(block, iv)
	cfb.CryptBlocks(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}
