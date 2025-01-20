package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[len(hash)-20:])
}

// 从字符串构造私钥
func privateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	// 解码十六进制私钥
	privateKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	// 确保私钥长度正确
	if len(privateKeyBytes) != 32 {
		return nil, fmt.Errorf("invalid private key length: got %d bytes, expected 32", len(privateKeyBytes))
	}

	// 创建私钥
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256() // 使用 P-256 椭圆曲线
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

func main() {
	// 生成私钥
	// privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKey, _ := privateKeyFromHex("0416847558fae607eeaa3746e01e79b3753c7feae2bb6587109651905d08995b")

	// 公钥（非压缩格式）
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// 地址生成
	address := generateAddress(publicKey)

	fmt.Printf("Private Key: %x\n", privateKey.D.Bytes())
	fmt.Printf("Public Key: %x\n", publicKey)
	fmt.Printf("Address: 0x%s\n", address)

	// 签名消息
	message := []byte("Hello, Swift Crypto!")
	hash := sha256.Sum256(message)
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("Signature: %x\n", signature)

	// 验证签名
	isValid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Printf("Is signature valid? %v\n", isValid)
}
