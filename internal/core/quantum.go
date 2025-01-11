package core

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

var DefaultLastSig = "00"

const (
	// QuantumTypeInformation specifies the quantum to post information. can be omitted
	QuantumTypeInformation = 0

	// QuantumTypeIntegration specifies the quantum to update user's (signer's) profile.
	QuantumTypeIntegration = 1

	//...
)

type QContent struct {
	Data   interface{} `json:"data,omitempty"` // string, int
	Format string      `json:"fmt"`            // number, txt, json, base64
}

type UnsignedQuantum struct {
	// Contents contain all data in this quantum
	Contents []*QContent `json:"cs,omitempty"`
	// Last contains the signature in last quantum
	Last string `json:"last"`
	// Nonce specifies the nonce of this quantum
	Nonce int `json:"nonce"`
	// References contains all references in this quantum
	References []string `json:"refs"`
	// Type specifies the type of this quantum
	Type int `json:"type,omitempty"`
}

type SignedQuantum struct {
	UnsignedQuantum
	Signature string `json:"sig,omitempty"`
	Signer    string `json:"signer,omitempty"`
}

func NewUnsignedQuantum(contents []*QContent, last string, nonce int, references []string) *UnsignedQuantum {
	return &UnsignedQuantum{
		Contents:   contents,
		Last:       last,
		Nonce:      nonce,
		References: references,
	}
}

func GenerateSignedJSON(privateKey *ecdsa.PrivateKey, unsignedQuantum UnsignedQuantum) ([]byte, error) {
	// （1）先把“待签名部分”序列化
	data, err := json.Marshal(unsignedQuantum)
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

	// （5）将原始字段 + 签名合并到 SignedQuantum 结构里
	signedQuantum := SignedQuantum{
		UnsignedQuantum: unsignedQuantum,
		Signature:       signatureHex,
	}

	// （6）再把完整结构序列化成 JSON，对外发送
	finalJSON, err := json.Marshal(signedQuantum)
	if err != nil {
		return nil, fmt.Errorf("final json.Marshal error: %v", err)
	}

	return finalJSON, nil
}

// ---------------------------------------------------
// 4) 验证从 JSON 中获取的数据
// ---------------------------------------------------
func VerifySignedJSON(jsonBytes []byte) (bool, string, error) {
	// （1）先解析出 UnsignedQuantum
	var signed SignedQuantum
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
	unsigned := UnsignedQuantum{
		Contents:   signed.Contents,
		Last:       signed.Last,
		Nonce:      signed.Nonce,
		References: signed.References,
		Type:       signed.Type,
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
