package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func generateAddress(publicKey []byte) string {
	hash := sha256.Sum256(publicKey)
	return hex.EncodeToString(hash[len(hash)-20:])
}

func main() {
	// 生成私钥
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// 公钥（非压缩格式）
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// 地址生成
	address := generateAddress(publicKey)

	fmt.Printf("Private Key: %x\n", privateKey.D.Bytes())
	fmt.Printf("Public Key: %x\n", publicKey)
	fmt.Printf("Address: 0x%s\n", address)

	// 签名消息
	message := []byte("Hello, Compatibility!")
	hash := sha256.Sum256(message)
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("Signature: %x\n", signature)

	// 验证签名
	isValid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Printf("Is signature valid? %v\n", isValid)
}
