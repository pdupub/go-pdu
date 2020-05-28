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

package pdu

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"

	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/crypto"
)

// PEngine is the engine of pdu
type PEngine struct {
	name string
}

// New is used to create PEngine
func New() *PEngine {
	return &PEngine{name: crypto.PDU}
}

// Name return name of pdu (PDU)
func (e PEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e PEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.GenKey(e.name, genKey, params...)
}

// parseKey parse the private key, return private key and public key
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
	return hex.EncodeToString(pk.D.Bytes()), hex.EncodeToString(fromECDSAPub(&pk.PublicKey)), nil

}

func parsePubKeyToString(pubKey interface{}) (string, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(fromECDSAPub(pk)), nil
}

// parsePriKey parse the private key, return private key
func parsePriKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	switch priKey.(type) {
	case *ecdsa.PrivateKey:
		pk = priKey.(*ecdsa.PrivateKey)
	case ecdsa.PrivateKey:
		*pk = priKey.(ecdsa.PrivateKey)
	case []byte:
		pk.D = new(big.Int).SetBytes(priKey.([]byte))
	case *big.Int:
		pk.D = new(big.Int).Set(priKey.(*big.Int))
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}

	pk.PublicKey.Curve = elliptic.P256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
}

func toECDSAPub(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(elliptic.P256(), pub)
	if x == nil {
		return nil, crypto.ErrInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

func fromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
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
		return toECDSAPub(pubKey.([]byte))
	case *big.Int:
		return toECDSAPub(pubKey.(*big.Int).Bytes())
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
	r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
	if err != nil {
		return nil, nil, err
	}

	rb := common.Bytes2Hash(r.Bytes())
	sb := common.Bytes2Hash(s.Bytes())
	signature := append(rb[:], sb[:]...)

	return signature, &pk.PublicKey, nil
}

// Sign is used to create signature of content by private key
func (e PEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	return crypto.Sign(e.name, hash, priKey, sign)
}

func verify(hash []byte, pubKey interface{}, signature []byte) (bool, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	return ecdsa.Verify(pk, hash, r, s), nil
}

func parseMulSig(signature []byte) [][]byte {
	size := len(signature) / 64
	sigs := make([][]byte, size)
	for i := 0; i < size; i++ {
		sigs[i] = signature[i*64 : (i+1)*64]
	}
	return sigs
}

// Verify is used to verify the signature
func (e PEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	return crypto.Verify(e.name, hash, sig, verify, parseMulSig)
}

// Unmarshal unmarshal private & public key from json
func (e PEngine) Unmarshal(privKeyBytes, pubKeyBytes []byte) (privKey *crypto.PrivateKey, pubKey *crypto.PublicKey, err error) {
	return crypto.Unmarshal(e.name, privKeyBytes, pubKeyBytes, parseKey, parsePubKey)
}

// Marshal marshal private & public key to json
func (e PEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	return crypto.Marshal(e.name, privKey, pubKey, parseKeyToString, parsePubKeyToString)
}

// Display display private & public key
func (e PEngine) DisplayKey(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (map[string]interface{}, map[string]interface{}, error) {
	return crypto.DisplayKey(e.name, privKey, pubKey, parseKeyToString, parsePubKeyToString)
}

// EncryptKey encryptKey into file
func (e PEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {
	return crypto.EncryptKey(e.name, priKey, pass, privKeyToKeyBytes)
}

func privKeyToKeyBytes(priKey interface{}) ([]byte, []byte, error) {
	pk, err := parsePriKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	address := eth.PubkeyToAddress(pk.PublicKey)
	keyBytes := pk.D.Bytes()
	return keyBytes, address[:], nil
}

// DecryptKey decrypt private key from file
func (e PEngine) DecryptKey(keyJSON []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.DecryptKey(e.name, keyJSON, pass, parseKey)
}

func genKey() (interface{}, interface{}, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}
