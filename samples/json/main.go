package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 1. 自定义的结构体，嵌入 jwt.RegisteredClaims
type MyCustomClaims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`

	// 通过嵌入 RegisteredClaims 来处理 exp, iat, nbf, iss, aud, sub 等标准字段
	jwt.RegisteredClaims
}

// GenerateToken 生成 Token（HS256 签名示例）
func GenerateToken(claims MyCustomClaims, secret []byte) (string, error) {
	// 设置常见的 JWT 字段（也可以在外部设置后传进来）
	claims.RegisteredClaims = jwt.RegisteredClaims{
		// 过期时间：2 小时后
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		// 签发时间
		IssuedAt: jwt.NewNumericDate(time.Now()),
		// 也可以设置别的：NotBefore, Subject, Issuer, Audience 等
	}

	// 创建 Token，并传入指定的签名算法和自定义 Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用 secret 进行签名，得到最终的 Token 字符串
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析并验证 Token
func ParseToken(tokenString string, secret []byte) (*MyCustomClaims, error) {
	// 用自定义结构体承接解析结果
	claims := &MyCustomClaims{}

	// 解析时，把 claims 实例传进去，这样库会自动把 JSON 反序列化到结构体里
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 校验 token.Header["alg"] 是否为预期算法等，可以在此处加检查
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	// 检查 Token 是否有效
	if !parsedToken.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// 此时 claims 已经被填充好，可以直接返回
	return claims, nil
}

func main() {
	// 准备签名用的 Secret
	secretKey := []byte("MySecretKey")

	// 需要写进 Token 的信息
	originalClaims := MyCustomClaims{
		UserID: 123,
		Name:   "Alice",
		Role:   "admin",
	}

	// 1) 生成 Token
	tokenString, err := GenerateToken(originalClaims, secretKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", tokenString)

	// 2) 解析并验证 Token
	parsedClaims, err := ParseToken(tokenString, secretKey)
	if err != nil {
		panic(err)
	}

	// 3) 使用解析后的数据
	fmt.Println("Recovered Claims:")
	fmt.Printf("  UserID: %d\n", parsedClaims.UserID)
	fmt.Printf("  Name:   %s\n", parsedClaims.Name)
	fmt.Printf("  Role:   %s\n", parsedClaims.Role)

	// 你也可以查看或校验标准字段
	fmt.Printf("  ExpiresAt: %v\n", parsedClaims.ExpiresAt)
	fmt.Printf("  IssuedAt:  %v\n", parsedClaims.IssuedAt)
}
