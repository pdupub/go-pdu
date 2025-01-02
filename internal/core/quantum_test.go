package core

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewUnsignedQuantum(t *testing.T) {
	type args struct {
		contents   []*QContent
		last       string
		nonce      int
		references []string
	}
	tests := []struct {
		name string
		args args
		want *UnsignedQuantum
	}{
		{
			name: "test",
			args: args{
				contents:   []*QContent{},
				last:       DefaultLastSig,
				nonce:      1,
				references: []string{},
			},
			want: &UnsignedQuantum{
				Contents:   []*QContent{},
				Last:       DefaultLastSig,
				Nonce:      1,
				References: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnsignedQuantum(tt.args.contents, tt.args.last, tt.args.nonce, tt.args.references); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnsignedQuantum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	// 生成一个私钥（secp256k1）
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("GenerateKey error: %v", err)
	}

	// 私钥 -> 地址
	ownerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	t.Log("原始私钥对应的地址：", ownerAddress.Hex())

	// 构造要签名的“quantum”
	quantum := UnsignedQuantum{
		Contents: []*QContent{
			{
				Data:   "hello world",
				Format: "string",
			},
			{
				Data:   []byte{0x01, 0x02, 0x03},
				Format: "binary",
			},
			{
				Data:   123,
				Format: "number",
			},
		},
		Last:       DefaultLastSig,
		Nonce:      1,
		References: []string{DefaultLastSig},
	}

	// 生成带签名的 JSON
	signedJSON, err := GenerateSignedJSON(privateKey, quantum)
	if err != nil {
		t.Errorf("GenerateSignedJSON error: %v", err)
	}

	t.Log("\n==== 生成的 JSON（包括签名）====")
	t.Log(string(signedJSON))

	// 模拟“接收方”收到这份 JSON，开始验证
	isValid, recoveredAddr, err := VerifySignedJSON(signedJSON)
	if err != nil {
		t.Errorf("VerifySignedJSON error: %v", err)
	}

	t.Log("\n==== 验证结果 ====")
	t.Log("签名是否有效？", isValid)
	t.Log("从签名恢复出的地址：", recoveredAddr)

	if recoveredAddr == ownerAddress.Hex() {
		t.Log("地址匹配，验证通过！")
	} else {
		t.Log("地址不一致，验证失败。")
	}
}
