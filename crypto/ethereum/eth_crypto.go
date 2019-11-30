// Copyright 2019 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/crypto"
	"math/big"
)

const (
	// SourceName is name of this
	SourceName = "ETH"
)

func genKey() (*ecdsa.PrivateKey, error) {
	return eth.GenerateKey()
}

// GenKey generate the private and public key pair
func GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	if len(params) == 0 {
		return nil, nil, crypto.ErrSigTypeNotSupport
	}
	sigType := params[0].(string)
	switch sigType {
	case crypto.Signature2PublicKey:
		pk, err := genKey()
		if err != nil {
			return nil, nil, err
		}
		return &crypto.PrivateKey{Source: SourceName, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: SourceName, SigType: crypto.Signature2PublicKey, PubKey: pk.PublicKey}, nil

	case crypto.MultipleSignatures:
		if len(params) == 1 {
			return nil, nil, crypto.ErrParamsMissing
		}
		var privKeys, pubKeys []interface{}
		for i := 0; i < params[1].(int); i++ {
			pk, err := genKey()
			if err != nil {
				return nil, nil, err
			}
			privKeys = append(privKeys, pk)
			pubKeys = append(pubKeys, pk.PublicKey)
		}
		return &crypto.PrivateKey{Source: SourceName, SigType: crypto.MultipleSignatures, PriKey: privKeys}, &crypto.PublicKey{Source: SourceName, SigType: crypto.MultipleSignatures, PubKey: pubKeys}, nil
	default:
		return nil, nil, crypto.ErrSigTypeNotSupport
	}
}

// ParsePriKey parse the private key
func ParsePriKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	switch priKey.(type) {
	case *ecdsa.PrivateKey:
		pk = priKey.(*ecdsa.PrivateKey)
	case ecdsa.PrivateKey:
		*pk = priKey.(ecdsa.PrivateKey)
	case []byte:
		return eth.ToECDSA(priKey.([]byte))
	case *big.Int:
		return eth.ToECDSA(priKey.(*big.Int).Bytes())
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

// ParsePubKey parse the public key
func ParsePubKey(pubKey interface{}) (*ecdsa.PublicKey, error) {
	pk := new(ecdsa.PublicKey)
	switch pubKey.(type) {
	case *ecdsa.PublicKey:
		pk = pubKey.(*ecdsa.PublicKey)
	case ecdsa.PublicKey:
		*pk = pubKey.(ecdsa.PublicKey)
	case []byte:
		return eth.UnmarshalPubkey(pubKey.([]byte))
	case *big.Int:
		return eth.UnmarshalPubkey(pubKey.(*big.Int).Bytes())
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

// Sign is used to create signature of content by private key
func Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != SourceName {
		return nil, crypto.ErrSourceNotMatch
	}
	switch priKey.SigType {
	case crypto.Signature2PublicKey:
		pk, err := ParsePriKey(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		signature, err := eth.Sign(signHash(hash), pk)
		if err != nil {
			return nil, err
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: SourceName, SigType: priKey.SigType, PubKey: pk.PublicKey},
			Signature: signature,
		}, nil
	case crypto.MultipleSignatures:
		pks := priKey.PriKey.([]interface{})
		var pubKeys []interface{}
		var signatures []byte
		for _, item := range pks {
			pk, err := ParsePriKey(item)
			if err != nil {
				return nil, err
			}
			signature, err := eth.Sign(signHash(hash), pk)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, signature...)
			pubKeys = append(pubKeys, pk.PublicKey)
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: SourceName, SigType: priKey.SigType, PubKey: pubKeys},
			Signature: signatures,
		}, nil
	default:
		return nil, crypto.ErrSigTypeNotSupport
	}
}

// Verify is used to verify the signature
func Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	if sig.Source != SourceName {
		return false, crypto.ErrSourceNotMatch
	}
	switch sig.SigType {
	case crypto.Signature2PublicKey:
		pk, err := ParsePubKey(sig.PubKey)
		if err != nil {
			return false, err
		}

		recoveredPubKey, err := eth.SigToPub(signHash(hash), sig.Signature)
		if err != nil || recoveredPubKey == nil {
			return false, err
		}

		if recoveredPubKey.Y.Cmp(pk.Y) != 0 || recoveredPubKey.X.Cmp(pk.X) != 0 {
			return false, nil
		}
		return true, nil

	case crypto.MultipleSignatures:
		pks := sig.PubKey.([]interface{})
		if len(pks) != len(sig.Signature)/65 {
			return false, crypto.ErrSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			pk, err := ParsePubKey(pubkey)
			if err != nil {
				return false, err
			}

			recoveredPubKey, err := eth.SigToPub(signHash(hash), sig.Signature[i*65:(i+1)*65])
			if err != nil || recoveredPubKey == nil {
				return false, err
			}

			if recoveredPubKey.Y.Cmp(pk.Y) != 0 || recoveredPubKey.X.Cmp(pk.X) != 0 {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, crypto.ErrSigTypeNotSupport
	}
}

// UnmarshalJSON unmarshal public key from json
func UnmarshalJSON(input []byte) (*crypto.PublicKey, error) {
	p := crypto.PublicKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == SourceName {
		if p.SigType == crypto.Signature2PublicKey {
			pk, err := hex.DecodeString(aMap["pubKey"].(string))
			if err != nil {
				return nil, err
			}
			pubKey, err := eth.UnmarshalPubkey(pk)
			if err != nil {
				return nil, err
			}
			p.PubKey = *pubKey
		} else if p.SigType == crypto.MultipleSignatures {
			pks := aMap["pubKey"].([]interface{})
			var pubKeys []ecdsa.PublicKey
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				pubKey, err := eth.UnmarshalPubkey(pk)
				if err != nil {
					return nil, err
				}
				pubKeys = append(pubKeys, *pubKey)
			}
			p.PubKey = pubKeys
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}

	return &p, nil
}

// MarshalJSON marshal public key to json
func MarshalJSON(a crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == SourceName {
		if a.SigType == crypto.Signature2PublicKey {
			pk := a.PubKey.(ecdsa.PublicKey)
			aMap["pubKey"] = hex.EncodeToString(eth.FromECDSAPub(&pk))
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PubKey.(type) {
			case []ecdsa.PublicKey:
				pks := a.PubKey.([]ecdsa.PublicKey)
				pubKey := make([]string, len(pks))
				for i, pk := range pks {
					pubKey[i] = hex.EncodeToString(eth.FromECDSAPub(&pk))
				}
				aMap["pubKey"] = pubKey
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks))
				for i, v := range pks {
					pk := v.(ecdsa.PublicKey)
					pubKey[i] = hex.EncodeToString(eth.FromECDSAPub(&pk))
				}
				aMap["pubKey"] = pubKey
			}

		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return eth.Keccak256([]byte(msg))
}
