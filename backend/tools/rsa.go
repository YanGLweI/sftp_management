package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sftpbackend/config"
	"sync"
)

var (
	raMutex sync.RWMutex
)

// Encrypt RSA 加密（用于前端传输数据到后端）
// 注意：使用公钥加密，确保只有持有私钥的后端才能解密
func Encrypt(plaintext string) (string, error) {
	raMutex.RLock()
	defer raMutex.RUnlock()

	publicKey := config.GlobalConfig.System.RSAPublicKey
	if publicKey == nil {
		return "", fmt.Errorf("RSA 公钥未加载")
	}

	// RSA-OAEP 加密
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("RSA-OAEP 加密失败：%w", err)
	}

	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return "ENC[" + encoded + "]", nil
}

// Decrypt RSA 解密（用于后端解密存储的数据）
// 注意：使用私钥解密，确保只有持有私钥的系统才能读取明文
func Decrypt(ciphertext string) (string, error) {
	raMutex.RLock()
	defer raMutex.RUnlock()

	privateKey := config.GlobalConfig.System.RSAPrivateKey
	if privateKey == nil {
		return ciphertext, fmt.Errorf("RSA 私钥未加载")
	}

	// 去除 ENC[] 前缀
	if len(ciphertext) < 7 || ciphertext[:4] != "ENC[" || ciphertext[len(ciphertext)-1] != ']' {
		return ciphertext, nil // 如果未加密，直接返回原值
	}

	encrypted := ciphertext[4 : len(ciphertext)-1]

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	// RSA-OAEP 解密
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
	if err != nil {
		return "", fmt.Errorf("RSA-OAEP decrypt error: %w", err)
	}

	return string(plaintext), nil
}

// DecryptPassword 解密前端传来的 RSA 加密密码（Base64 编码，PKCS1v15 模式）
// 注意：使用私钥解密，确保只有持有私钥的系统才能读取明文
func DecryptPassword(encryptedPassword string) (string, error) {
	raMutex.RLock()
	defer raMutex.RUnlock()

	privateKey := config.GlobalConfig.System.RSAPrivateKey
	if privateKey == nil {
		return "", fmt.Errorf("RSA 私钥未加载")
	}

	// Base64 解码（前端加密后是 Base64 字符串，无 ENC[] 包装）
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	// RSA-PKCS1v15 解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", fmt.Errorf("RSA decrypt error: %w", err)
	}

	return string(plaintext), nil
}
