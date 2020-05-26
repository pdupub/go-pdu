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

package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"

	btc "github.com/btcsuite/btcd/btcec"
	"github.com/pdupub/go-pdu/crypto"

	eth "github.com/ethereum/go-ethereum/crypto"
)

// BEngine is the engine of BTC
type BEngine struct {
	name string
}

// New create new BEngine
func New() *BEngine {
	return &BEngine{name: crypto.BTC}
}

// Name return the name of this engine (BTC)
func (e BEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e BEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.GenKey(e.name, genKey, params...)
}

// parsePriKey parse the private key
func parseKey(privKey interface{}) (interface{}, interface{}, error) {
	pk, err := parsePriKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	return pk, &pk.PublicKey, nil
}

// parsePriKey parse the private key
func parsePriKey(priKey interface{}) (*btc.PrivateKey, error) {
	pk := new(btc.PrivateKey)
	switch priKey.(type) {
	case *btc.PrivateKey:
		pk = priKey.(*btc.PrivateKey)
	case btc.PrivateKey:
		*pk = priKey.(btc.PrivateKey)
	case []byte:
		privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), priKey.([]byte))
		return privateKey, nil
	case *big.Int:
		privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), priKey.(*big.Int).Bytes())
		return privateKey, nil
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

func fromECDSAPub(publicKey *ecdsa.PublicKey) []byte {
	pk := (*btc.PublicKey)(publicKey)
	return pk.SerializeUncompressed()
}

// parsePubKey parse the public key
func parsePubKey(pubKey interface{}) (*ecdsa.PublicKey, error) {
	pk := new(ecdsa.PublicKey)
	switch pubKey.(type) {
	case *ecdsa.PublicKey:
		pk = pubKey.(*ecdsa.PublicKey)
	case ecdsa.PublicKey:
		*pk = pubKey.(ecdsa.PublicKey)
	case []byte:
		publicKey, err := btc.ParsePubKey(pubKey.([]byte), btc.S256())
		if err != nil {
			return nil, err
		}
		return publicKey.ToECDSA(), nil
	case *big.Int:
		publicKey, err := btc.ParsePubKey(pubKey.(*big.Int).Bytes(), btc.S256())
		if err != nil {
			return nil, err
		}
		return publicKey.ToECDSA(), nil
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

func sign(hash []byte, privKey interface{}) ([]byte, *ecdsa.PublicKey, error) {
	pk, err := parsePriKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	signature, err := pk.Sign(hash)
	if err != nil {
		return nil, nil, err
	}
	return signature.Serialize(), &pk.PublicKey, nil
}

// Sign is used to create signature of content by private key
func (e BEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	return crypto.Sign(e.name, hash, priKey, sign)
}

func verify(hash []byte, pubKey interface{}, sig []byte) (bool, error) {
	signature, err := btc.ParseSignature(sig, btc.S256())
	if err != nil {
		return false, err
	}
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	return signature.Verify(hash, (*btc.PublicKey)(pk)), nil
}

func parseMulSig(signature []byte) [][]byte {
	currentSigLeft := -1
	var currentSig []byte
	var sigs [][]byte
	for _, v := range signature {
		if currentSigLeft == -1 {
			currentSigLeft = 0
			currentSig = []byte{}
			currentSig = append(currentSig, v)
		} else if currentSigLeft == 0 {
			currentSigLeft = int(v)
			currentSig = append(currentSig, v)
		} else {
			currentSig = append(currentSig, v)
			currentSigLeft--
			if currentSigLeft == 0 {
				sigs = append(sigs, currentSig)
				currentSigLeft = -1
			}
		}
	}
	return sigs
}

// Verify is used to verify the signature
func (e BEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	return crypto.Verify(e.name, hash, sig, verify, parseMulSig)
}

// Unmarshal unmarshal private & public key
func (e BEngine) Unmarshal(privKeyBytes, pubKeyBytes []byte) (privKey *crypto.PrivateKey, pubKey *crypto.PublicKey, err error) {

	if len(pubKeyBytes) > 0 {
		if pubKey, err = e.unmarshalPubKey(pubKeyBytes); err != nil {
			return
		}
	}
	if len(privKeyBytes) > 0 {
		if privKey, err = e.unmarshalPrivKey(privKeyBytes); err != nil {
			return
		}
	}
	return
}

func (e BEngine) unmarshalPrivKey(input []byte) (*crypto.PrivateKey, error) {
	p := crypto.PrivateKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == e.name {
		if p.SigType == crypto.Signature2PublicKey {
			pk, err := hex.DecodeString(aMap["privKey"].(string))
			if err != nil {
				return nil, err
			}
			privKey, err := parsePriKey(pk)
			if err != nil {
				return nil, err
			}
			p.PriKey = privKey
		} else if p.SigType == crypto.MultipleSignatures {
			pks := aMap["privKey"].([]interface{})
			var privKeys []interface{}
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				privKey, err := parsePriKey(pk)
				if err != nil {
					return nil, err
				}
				privKeys = append(privKeys, privKey)
			}
			p.PriKey = privKeys
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}

	return &p, nil
}

func (e BEngine) unmarshalPubKey(input []byte) (*crypto.PublicKey, error) {
	p := crypto.PublicKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == e.name {
		if p.SigType == crypto.Signature2PublicKey {
			pk, err := hex.DecodeString(aMap["pubKey"].(string))
			if err != nil {
				return nil, err
			}
			pubKey, err := parsePubKey(pk)
			if err != nil {
				return nil, err
			}
			p.PubKey = pubKey
		} else if p.SigType == crypto.MultipleSignatures {
			pks := aMap["pubKey"].([]interface{})
			var pubKeys []interface{}
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				pubKey, err := parsePubKey(pk)
				if err != nil {
					return nil, err
				}
				pubKeys = append(pubKeys, pubKey)
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

// Marshal marshal private & public key
func (e BEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	if privKey != nil {
		if privKeyBytes, err = e.marshalPrivKey(privKey); err != nil {
			return
		}
	}
	if pubKey != nil {
		if pubKeyBytes, err = e.marshalPubKey(pubKey); err != nil {
			return
		}
	}
	return
}

func (e BEngine) marshalPrivKey(a *crypto.PrivateKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk := a.PriKey.(*btc.PrivateKey)
			bpk := (btc.PrivateKey)(*pk)
			aMap["privKey"] = hex.EncodeToString(bpk.Serialize())
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PriKey.(type) {
			case []interface{}:
				pks := a.PriKey.([]interface{})
				privKey := make([]string, len(pks))
				for i, v := range pks {
					pk := v.(*btc.PrivateKey)
					bpk := (btc.PrivateKey)(*pk)
					privKey[i] = hex.EncodeToString(bpk.Serialize())
				}
				aMap["privKey"] = privKey
			}

		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}

func (e BEngine) marshalPubKey(a *crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			aMap["pubKey"] = hex.EncodeToString(fromECDSAPub(a.PubKey.(*ecdsa.PublicKey)))
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PubKey.(type) {
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks))
				for i, v := range pks {
					pubKey[i] = hex.EncodeToString(fromECDSAPub(v.(*ecdsa.PublicKey)))
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

// EncryptKey encryptKey into file
func (e BEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {
	return crypto.EncryptKey(e.name, priKey, pass, privKeyToKeyBytes)
}

func privKeyToKeyBytes(priKey interface{}) ([]byte, []byte, error) {
	btcpk, err := parsePriKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	pk := btcpk.ToECDSA()
	address := eth.PubkeyToAddress(pk.PublicKey)
	keyBytes := pk.D.Bytes()
	return keyBytes, address[:], nil
}

// DecryptKey decrypt private key from file
func (e BEngine) DecryptKey(keyJSON []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.DecryptKey(e.name, keyJSON, pass, parseKey)
}

func genKey() (interface{}, interface{}, error) {
	privKey, err := btc.NewPrivateKey(btc.S256())
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}
