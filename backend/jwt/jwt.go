package jwt

import (
	"sftpbackend/config"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 读取jwt配置文件
var jwt_config = config.GlobalConfig.JWT

// 自定义密钥
var mySecret = []byte(jwt_config.Secret)

// 定义过期时间: X小时
var expirationTime = time.Hour * time.Duration(jwt_config.Expire)

// 用于存储已失效的Token，使用互斥锁来保证并发安全
var InvalidTokens sync.Map

// ! GenerateToken 生成JWT
func GenerateToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(), // 过期时间
			Issuer:    jwt_config.Issuer,                     // 签发人
		},
	}
	// 创建token：使用HS256算法，传入结构体实例
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 生成token字符串：传入签名密钥
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ! ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token,
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		// 返回解析后的claims
		return claims, nil
	}
	return nil, err
}

// ! 标记Token过期
func MarkTokenExpired(tokenString string) {
	// 将该Token标记为失效，存入invalidTokens集合中
	InvalidTokens.Store(tokenString, true)
}
