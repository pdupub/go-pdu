// Copyright 2021 The PDU Authors
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

package identity

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/google/uuid"
)

// EncryptedKeyJSONV3 is from geth
type EncryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}

func auth(passwd, salt []byte, addr common.Address) string {
	return fmt.Sprintf("%s%s%s%s", passwd, addr.Hex()[2:4], addr.Hex()[40:42], salt)
}

// GeneratePrivateKey generate random private key
func GeneratePrivateKey() ([]byte, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return crypto.FromECDSA(privateKey), nil
}

func InspectKeystore(keyjson, passwd []byte) ([]byte, error) {

	key, err := keystore.DecryptKey(keyjson, string(passwd))
	if err != nil {
		return nil, err
	}
	return crypto.FromECDSA(key.PrivateKey), nil
}

// GenerateKeystore generate keystore file by json
func GenerateKeystore(privKeyBytes, passwd []byte) ([]byte, error) {

	privateKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, err
	}

	currentAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	UUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	cryptoStruct, err := keystore.EncryptDataV3(privateKey.D.Bytes(), passwd, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	encryptedKeyJSONV3 := EncryptedKeyJSONV3{
		hex.EncodeToString(currentAddress[:]),
		cryptoStruct,
		UUID.String(),
		3,
	}
	return json.Marshal(encryptedKeyJSONV3)
}

// CreateKeystoreArrayAndSaveLocal used to generate some keystore and save to local path for start up test
func CreateKeystoreArrayAndSaveLocal(keyDir string, passwd, salt []byte, cnt int) ([]*DID, error) {

	var dids []*DID
	if _, err := os.Stat(keyDir); err == nil {
		// Directory already exists at keyFilePath
		return dids, errors.New("directory already exists at keyFilePath")
	} else if !os.IsNotExist(err) {
		// Error checking if keyfile exists
		return dids, err
	}

	for i := 0; i < cnt; i++ {
		var privateKey *ecdsa.PrivateKey

		// If not loaded, generate random.
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return dids, err
		}

		currentAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

		UUID, err := uuid.NewUUID()
		if err != nil {
			return dids, err
		}

		dids = append(dids, &DID{key: &keystore.Key{
			Id:         UUID[:],
			Address:    currentAddress,
			PrivateKey: privateKey,
		}})
		cryptoStruct, err := keystore.EncryptDataV3(privateKey.D.Bytes(), []byte(auth(passwd, salt, currentAddress)), keystore.StandardScryptN, keystore.StandardScryptP)
		if err != nil {
			return dids, err
		}
		encryptedKeyJSONV3 := EncryptedKeyJSONV3{
			hex.EncodeToString(currentAddress[:]),
			cryptoStruct,
			UUID.String(),
			3,
		}
		keyJSON, err := json.Marshal(encryptedKeyJSONV3)
		if err != nil {
			return dids, err
		}

		filename := fmt.Sprintf("%s.json", currentAddress.Hex())
		keyFilePath := path.Join(keyDir, filename)
		if _, err := os.Stat(keyFilePath); err == nil {
			// Keyfile already exists at keyFilePath
			return dids, err
		} else if !os.IsNotExist(err) {
			// Error checking if keyfile exists
			return dids, err
		}
		// Store the file to disk.
		if err := os.MkdirAll(filepath.Dir(keyFilePath), 0700); err != nil {
			// Could not create directory keyFilePath
			return dids, err
		}
		if err := ioutil.WriteFile(keyFilePath, keyJSON, 0600); err != nil {
			// Failed to write keyfile to keyFilePath
			return dids, err
		}

	}

	return dids, nil
}

func UnlockKeystoreArray(keyDir string, passwd, salt []byte) ([]*DID, error) {
	var ids []*DID
	rd, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			res := strings.Split(fi.Name(), ".")
			if len(res) == 2 && res[1] == "json" {
				did, err := unlockKeystore(path.Join(keyDir, fi.Name()), auth(passwd, salt, common.HexToAddress(res[0])))
				if err != nil {
					return nil, err
				}
				ids = append(ids, did)
			}

		}
	}
	return ids, nil
}

func unlockKeystore(filename string, auth string) (*DID, error) {
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	return &DID{key: key}, nil
}
