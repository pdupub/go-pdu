package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

// ---------------------------------------------------
// 1) 定义要签名的字段（不含签名）
// ---------------------------------------------------
type UnsignedClaims struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Platform  string `json:"platform"`
}

// ---------------------------------------------------
// 2) 定义最终发送/接收的数据结构（含签名）
// ---------------------------------------------------
type SignedClaims struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Platform  string `json:"platform"`

	// 这里多了一个字段用来放签名（hex 格式）
	Signature string `json:"signature"`
}

// ---------------------------------------------------
// 3) 生成带签名的 JSON
// ---------------------------------------------------
func GenerateSignedJSON(privateKey *ecdsa.PrivateKey, claims UnsignedClaims) ([]byte, error) {
	// （1）先把“待签名部分”序列化
	data, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error: %v", err)
	}

	// （2）Keccak256 哈希
	hash := crypto.Keccak256Hash(data)

	// （3）对哈希做 secp256k1 签名
	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign error: %v", err)
	}

	// （4）把签名转成 hex 字符串
	signatureHex := hex.EncodeToString(signatureBytes)

	// （5）将原始字段 + 签名合并到 SignedClaims 结构里
	signedClaims := SignedClaims{
		Message:   claims.Message,
		Timestamp: claims.Timestamp,
		Platform:  claims.Platform,
		Signature: signatureHex,
	}

	// （6）再把完整结构序列化成 JSON，对外发送
	finalJSON, err := json.Marshal(signedClaims)
	if err != nil {
		return nil, fmt.Errorf("final json.Marshal error: %v", err)
	}

	return finalJSON, nil
}

// ---------------------------------------------------
// 4) 验证从 JSON 中获取的数据
// ---------------------------------------------------
func VerifySignedJSON(jsonBytes []byte) (bool, string, error) {
	// （1）先解析出 SignedClaims
	var signed SignedClaims
	err := json.Unmarshal(jsonBytes, &signed)
	if err != nil {
		return false, "", fmt.Errorf("json.Unmarshal error: %v", err)
	}

	// （2）取出 signature
	signatureBytes, err := hex.DecodeString(signed.Signature)
	if err != nil {
		return false, "", fmt.Errorf("DecodeString error: %v", err)
	}

	// （3）把“签名以外的字段”组装成 UnsignedClaims，并序列化
	unsigned := UnsignedClaims{
		Message:   signed.Message,
		Timestamp: signed.Timestamp,
		Platform:  signed.Platform,
	}
	data, err := json.Marshal(unsigned)
	if err != nil {
		return false, "", fmt.Errorf("json.Marshal error: %v", err)
	}

	// （4）再做一次哈希
	hash := crypto.Keccak256Hash(data)

	// （5）用 Ecrecover 恢复公钥
	recoveredPub, err := crypto.Ecrecover(hash.Bytes(), signatureBytes)
	if err != nil {
		return false, "", fmt.Errorf("Ecrecover error: %v", err)
	}

	pubKeyECDSA, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		return false, "", fmt.Errorf("UnmarshalPubkey error: %v", err)
	}

	// （6）推导出恢复地址
	recoveredAddress := crypto.PubkeyToAddress(*pubKeyECDSA).Hex()

	// 如果你希望进一步对比“是否与某个预期地址匹配”，在这里做对比
	return true, recoveredAddress, nil
}

// ---------------------------------------------------
// 5) 主函数：生成私钥 -> 生成带签名的 JSON -> 再验证
// ---------------------------------------------------
func main() {
	// 生成一个私钥（secp256k1）
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("GenerateKey error: %v", err)
	}

	// 私钥 -> 地址
	ownerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("原始私钥对应的地址：", ownerAddress.Hex())

	// 构造要签名的“claims”
	claims := UnsignedClaims{
		Message:   "Hello from secp256k1",
		Timestamp: 1234567890,
		Platform:  "Go + Ethereum Crypto",
	}

	// 生成带签名的 JSON
	signedJSON, err := GenerateSignedJSON(privateKey, claims)
	if err != nil {
		log.Fatalf("GenerateSignedJSON error: %v", err)
	}

	fmt.Println("\n==== 生成的 JSON（包括签名）====")
	fmt.Println(string(signedJSON))

	// 模拟“接收方”收到这份 JSON，开始验证
	isValid, recoveredAddr, err := VerifySignedJSON(signedJSON)
	if err != nil {
		log.Fatalf("VerifySignedJSON error: %v", err)
	}

	fmt.Println("\n==== 验证结果 ====")
	fmt.Println("签名是否有效？", isValid)
	fmt.Println("从签名恢复出的地址：", recoveredAddr)

	if recoveredAddr == ownerAddress.Hex() {
		fmt.Println("地址匹配，验证通过！")
	} else {
		fmt.Println("地址不一致，验证失败。")
	}
}
