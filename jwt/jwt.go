// Package jwt 提供了JWT（JSON Web Token）的生成、验证和刷新功能
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义了JWT的载荷结构
// 包含用户ID、用户名和标准JWT声明
type Claims struct {
	UserID               string `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims        // 标准JWT声明（过期时间、签发时间等）
}

// JWTManager 是JWT管理器
// 负责JWT令牌的生成、验证和刷新
type JWTManager struct {
	secretKey     []byte        // 用于签名的密钥
	tokenDuration time.Duration // 令牌有效期
}

// NewJWTManager 创建一个新的JWT管理器
// secretKey: 用于签名的密钥
// duration: 令牌有效期
func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: duration,
	}
}

// GenerateToken 生成JWT令牌
// userID: 用户ID
// username: 用户名
// 返回生成的令牌字符串和可能的错误
func (m *JWTManager) GenerateToken(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)), // 设置过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                      // 设置签发时间
			Issuer:    "easygo",                                            // 设置签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// VerifyToken 验证JWT令牌
// tokenString: 要验证的令牌字符串
// 返回令牌的载荷和可能的错误
func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新JWT令牌
// tokenString: 要刷新的令牌字符串
// 返回新的令牌字符串和可能的错误
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(m.tokenDuration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}
