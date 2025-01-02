package core

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
	Last string `json:"last,omitempty"`
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
}

func NewUnsignedQuantum(contents []*QContent, last string, nonce int, references []string) *UnsignedQuantum {
	return &UnsignedQuantum{
		Contents:   contents,
		Last:       last,
		Nonce:      nonce,
		References: references,
	}
}
