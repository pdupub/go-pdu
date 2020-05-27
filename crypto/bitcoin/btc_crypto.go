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

func parseKeyToString(privKey interface{}) (string, string, error) {
	pk, err := parsePriKey(privKey)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(pk.Serialize()), hex.EncodeToString(fromECDSAPub(&pk.PublicKey)), nil
}

func parsePubKeyToString(pubKey interface{}) (string, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(fromECDSAPub(pk)), nil
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
	return crypto.Unmarshal(e.name, privKeyBytes, pubKeyBytes, parseKey, parsePubKey)
}

// Marshal marshal private & public key
func (e BEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	return crypto.Marshal(e.name, privKey, pubKey, parseKeyToString, parsePubKeyToString)
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
