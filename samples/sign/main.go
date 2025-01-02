package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	jwt "github.com/golang-jwt/jwt/v4"
)

// 声明一个自定义的 secp256k1 SigningMethod
// 默认的 ES256/ES384/ES512 使用 P-256/P-384/P-521，需要我们手动扩展支持 secp256k1。
var SigningMethodES256K *signingMethodES256K

func init() {
	SigningMethodES256K = &signingMethodES256K{Name: "ES256K"}
	// 注册到 jwt 库的全局表，下文会用到 signingMethodES256K.Alg() 获取此对象
	jwt.RegisterSigningMethod(SigningMethodES256K.Alg(), func() jwt.SigningMethod {
		return SigningMethodES256K
	})
}

// 自定义签名方法结构体
type signingMethodES256K struct {
	Name string
}

// Alg 返回算法标识
func (m *signingMethodES256K) Alg() string {
	return m.Name
}

// Verify 验证签名
func (m *signingMethodES256K) Verify(signingString string, sig string, key interface{}) error {
	ecdsaPubKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("invalid public key type, expect ecdsa.PublicKey")
	}

	// 签名是 Base64URL 编码之后的，需要先解码
	signatureBytes, err := jwt.DecodeSegment(sig)
	if err != nil {
		return err
	}
	// secp256k1 的签名长度通常为 64字节 (R,S)
	if len(signatureBytes) != 64 {
		return errors.New("invalid signature length")
	}

	// 计算哈希后，使用 ecdsa.Verify() 验签
	hash := crypto.Keccak256Hash([]byte(signingString)) // 以太坊中常用 Keccak256
	r := signatureBytes[:32]
	s := signatureBytes[32:]

	// 转为大整数
	rInt := new(big.Int).SetBytes(r)
	sInt := new(big.Int).SetBytes(s)

	if ecdsa.Verify(ecdsaPubKey, hash.Bytes(), rInt, sInt) {
		return nil
	}
	return errors.New("signature verification failed")
}

// Sign 使用私钥对消息签名
func (m *signingMethodES256K) Sign(signingString string, key interface{}) (string, error) {
	ecdsaPrivKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return "", errors.New("invalid private key type, expect ecdsa.PrivateKey")
	}

	// 先哈希 (以太坊 / 大多数区块链中常用 Keccak256)
	hash := crypto.Keccak256Hash([]byte(signingString))
	sig, err := crypto.Sign(hash.Bytes(), ecdsaPrivKey)
	if err != nil {
		return "", err
	}

	// eth_sign 的返回 sig 有 65字节: (R, S, V)
	// 其中最后一字节 V 对 ECDSA 验签没用，标准 ECDSA 只有 (R, S)
	// 这里只保留 R, S (前64字节)
	sigRS := sig[:64]

	// 将签名做 Base64URL 编码得到字符串
	return jwt.EncodeSegment(sigRS), nil
}

func main() {
	// 1. 生成 secp256k1 临时私钥
	privKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("generate key failed: %v", err)
	}
	// 导出原始私钥字节（32字节）
	rawPrivKey := crypto.FromECDSA(privKey)
	// 显示私钥 (Hex)
	fmt.Printf("Ephemeral Private Key (hex): %s\n", hex.EncodeToString(rawPrivKey))

	// 2. 构建一个 JWT 的 claims
	claims := jwt.MapClaims{
		"message":   "Hello from secp256k1",
		"timestamp": 1234567890,
		"platform":  "Go + Ethereum Crypto",
	}

	// 3. 创建 token 并指定我们自定义的 ES256K 签名方法
	token := jwt.NewWithClaims(SigningMethodES256K, claims)
	// 签名
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		log.Fatalf("sign token failed: %v", err)
	}
	fmt.Println("Signed token:", tokenString)

	// 4. 用公钥验证 (以演示完整流程)
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 验证所用的 signing method 是否匹配
		if t.Method.Alg() != SigningMethodES256K.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return pubKey, nil
	})

	if err != nil {
		fmt.Println("验证失败:", err)
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Println("验证成功, payload:")
		for k, v := range claims {
			fmt.Printf("  %s: %v\n", k, v)
		}
	} else {
		fmt.Println("验证失败或无效 token")
	}
}
