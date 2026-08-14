package jwt

import (
	"sftpbackend/config"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims 基础JWT声明
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CustomClaims 扩展JWT声明，包含登录类型、角色和路由权限
type CustomClaims struct {
	Username  string   `json:"username"`
	LoginType string   `json:"login_type"`
	RoleID    *uint    `json:"role_id"`
	Routes    []string `json:"routes"`
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

// GenerateToken 生成JWT（兼容旧接口）
func GenerateToken(username string) (string, error) {
	c := MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
			Issuer:    jwt_config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateTokenWithClaims 生成包含完整声明的JWT
func GenerateTokenWithClaims(claims *CustomClaims) (string, error) {
	c := CustomClaims{
		Username:  claims.Username,
		LoginType: claims.LoginType,
		RoleID:    claims.RoleID,
		Routes:    claims.Routes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
			Issuer:    jwt_config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateLimitedToken 生成受限Token（仅可用于改密）
func GenerateLimitedToken(username string) (string, error) {
	c := CustomClaims{
		Username:  username,
		LoginType: "local",
		Routes:    []string{"ChangePasswordOnly"},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(), // 30分钟过期
			Issuer:    jwt_config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析JWT（支持 MyClaims 和 CustomClaims）
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 先尝试解析为 CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}

	// 再尝试解析为 MyClaims（兼容旧Token）
	token2, err2 := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err2 == nil {
		if claims, ok := token2.Claims.(*MyClaims); ok && token2.Valid {
			return &CustomClaims{
				Username:  claims.Username,
				LoginType: "",
				RoleID:    nil,
				Routes:    []string{},
			}, nil
		}
	}

	return nil, err
}

// MarkTokenExpired 标记Token过期
func MarkTokenExpired(tokenString string) {
	InvalidTokens.Store(tokenString, true)
}