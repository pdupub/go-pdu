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
	"fmt"

	eth "github.com/ethereum/go-ethereum/crypto"

	"math/big"

	"github.com/pdupub/go-pdu/crypto"
)

// EEngine is the engine of ETH
type EEngine struct {
	name string
}

// New create new EEngine
func New() *EEngine {
	return &EEngine{name: crypto.ETH}
}

// Name return the name of this engine (ETH)
func (e EEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e EEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
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
	return hex.EncodeToString(eth.FromECDSA(pk)), hex.EncodeToString(eth.FromECDSAPub(&pk.PublicKey)), nil
}

func parsePubKeyToString(pubKey interface{}) (string, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(eth.FromECDSAPub(pk)), nil
}

// parsePriKey parse the private key
func parsePriKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
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

	pk.PublicKey.Curve = eth.S256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
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
		return eth.UnmarshalPubkey(pubKey.([]byte))
	case *big.Int:
		return eth.UnmarshalPubkey(pubKey.(*big.Int).Bytes())
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
	signature, err := eth.Sign(signHash(hash), pk)
	if err != nil {
		return nil, nil, err
	}
	return signature, &pk.PublicKey, nil
}

// Sign is used to create signature of content by private key
func (e EEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	return crypto.Sign(e.name, hash, priKey, sign)
}

func verify(hash []byte, pubKey interface{}, signature []byte) (bool, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	recoveredPubKey, err := eth.SigToPub(signHash(hash), signature)
	if err != nil || recoveredPubKey == nil {
		return false, err
	}

	if recoveredPubKey.Y.Cmp(pk.Y) != 0 || recoveredPubKey.X.Cmp(pk.X) != 0 {
		return false, nil
	}
	return true, nil
}

func parseMulSig(signature []byte) [][]byte {
	size := len(signature) / 65
	sigs := make([][]byte, size)
	for i := 0; i < size; i++ {
		sigs[i] = signature[i*65 : (i+1)*65]
	}
	return sigs
}

// Verify is used to verify the signature
func (e EEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	return crypto.Verify(e.name, hash, sig, verify, parseMulSig)
}

// Unmarshal unmarshal private & public key from json
func (e EEngine) Unmarshal(privKeyBytes, pubKeyBytes []byte) (privKey *crypto.PrivateKey, pubKey *crypto.PublicKey, err error) {
	return crypto.Unmarshal(e.name, privKeyBytes, pubKeyBytes, parseKey, parsePubKey)
}

// Marshal marshal private & public key to json
func (e EEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	return crypto.Marshal(e.name, privKey, pubKey, parseKeyToString, parsePubKeyToString)
}

// EncryptKey encryptKey into file
func (e EEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {
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
func (e EEngine) DecryptKey(keyJSON []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.DecryptKey(e.name, keyJSON, pass, parseKey)
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return eth.Keccak256([]byte(msg))
}

func genKey() (interface{}, interface{}, error) {
	privKey, err := eth.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}
