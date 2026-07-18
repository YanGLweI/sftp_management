package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"sftpbackend/config"
)

// ! 解密密码
func DecryptPassword(encryptedPassword string) (string, error) {
	privateKey := config.GlobalConfig.System.RSAPrivateKey
	// 1. Base64解码
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	// 2. RSA解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", fmt.Errorf("RSA decrypt error: %w", err)
	}

	return string(plaintext), nil
}
