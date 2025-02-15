package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

// EncryptCTR using CTR mode.
// @param plainText: [string] The plain text.
// @param secretKey: [string] The secret key.
// @return The encrypted text and error. When an error occurs, the encrypted text will be empty string.
func EncryptCTR(plainText string, secretKey string) (string, error) {
	p := []byte(plainText)
	s := []byte(padLeft16Times(secretKey))
	block, err := aes.NewCipher(s)
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(p))
	iv := cipherText[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], p)
	result := hex.EncodeToString(cipherText)
	return result, nil
}

// DecryptCTR using CTR mode.
// @param cipherText: [string] The encrypted text.
// @param secretKey: [string] The secret key.
// @return The decrypted text and error. When an error occurs, the encrypted text will be empty string.
func DecryptCTR(cipherText string, secretKey string) (string, error) {
	decoded, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	c := decoded
	s := []byte(padLeft16Times(secretKey))
	block, err := aes.NewCipher(s)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(c[aes.BlockSize:]))
	stream := cipher.NewCTR(block, c[:aes.BlockSize])
	stream.XORKeyStream(decrypted, c[aes.BlockSize:])
	return string(decrypted), nil
}

// 0 Padding with 16 times from left side.
// @param text: [string] The text to pad with 0.
func padLeft16Times(text string) string {
	padCnt := aes.BlockSize - len(text)%aes.BlockSize
	if padCnt%aes.BlockSize == 0 {
		return text
	} else {
		return strings.Repeat("0", padCnt) + text
	}
}
