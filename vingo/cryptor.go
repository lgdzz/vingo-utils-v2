package vingo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/duke-git/lancet/v2/cryptor"
	"io"
)

// 生成secret
// 填0默认16字节
func GenerateRandomKey(size int) string {
	if size == 0 {
		size = 16
	}
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		panic(err) // 处理错误
	}
	return base64.StdEncoding.EncodeToString(key)
}

// 文本加密
func TextEncrypt(plaintext string, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 使用CTR模式进行加密
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plaintext))

	// 返回Base64编码的加密文本
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// 文本解密
func TextDecrypt(cipherText string, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 获取IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)

	// 解密
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

// base64编码
func TextBase64Encode(text string) string {
	return cryptor.Base64StdEncode(text)
}

// base64解码
func TextBase64Decode(text string) string {
	return cryptor.Base64StdDecode(text)
}
