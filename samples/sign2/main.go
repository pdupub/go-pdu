package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

// 示例 Claims，类似 jwt.MapClaims
type MyClaims struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Platform  string `json:"platform"`
}

func main() {
	// 1. 生成 secp256k1 私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("GenerateKey error: %v", err)
	}

	// 2. 私钥转 hex
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	fmt.Println("► 私钥 (hex):", privateKeyHex)

	// 2.1 由私钥推导出地址（即“钱包地址”）
	ownerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("► 由此私钥推导出的地址:", ownerAddress.Hex())

	// 3. 构造示例 Claims（实际上可替换为任何数据）
	claims := MyClaims{
		Message:   "Hello from secp256k1",
		Timestamp: 1234567890,
		Platform:  "Go + Ethereum Crypto",
	}

	// 4. 把 claims 序列化为 JSON，并计算 Keccak256 哈希
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		log.Fatalf("json.Marshal error: %v", err)
	}
	// 以太坊常用的哈希函数是 Keccak256
	hash := crypto.Keccak256Hash(claimsJSON)

	fmt.Println("► 原始 Claims JSON:", string(claimsJSON))
	fmt.Println("► Keccak256 hash :", hex.EncodeToString(hash.Bytes()))

	// 5. 使用私钥对 Keccak256 哈希进行签名
	//    签名结果为 65 字节：前 64 字节是 (R, S)，最后 1 字节是 “恢复标志” V。
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("Sign error: %v", err)
	}
	fmt.Println("► 签名结果 (hex)   :", hex.EncodeToString(signature))

	// 6. 从签名恢复公钥，再推导出地址。这样就能证明“谁签了这份数据”。
	//    Ecrecover 返回的是未压缩格式的公钥（65字节，开头 0x04）
	recoveredPub, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatalf("Ecrecover error: %v", err)
	}
	// 转成 ecdsa.PublicKey
	pubKeyECDSA, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		log.Fatalf("UnmarshalPubkey error: %v", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKeyECDSA)
	fmt.Println("► 由签名恢复出的地址:", recoveredAddress.Hex())

	// 7. 对比“原始私钥对应的地址”和“恢复出的地址”
	if ownerAddress == recoveredAddress {
		fmt.Println("==> 地址一致，验证成功！")
	} else {
		fmt.Println("==> 地址不匹配，验证失败！")
	}
}
