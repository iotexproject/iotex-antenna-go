// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package did

import (
	"encoding/hex"

	"github.com/iotexproject/go-pkgs/crypto"
)

const (
	// DIDContext is the general context
	DIDContext = "https://www.w3.org/ns/did/v1"
	// DIDPrefix is the prefix string
	DIDPrefix = "did:io:"
	// DIDAuthType is the authentication type
	DIDAuthType = "EcdsaSecp256k1VerificationKey2019"
	// DIDOwner is the suffix string
	DIDOwner = "#owner"
)

type (
	authentication struct {
		ID           string `json:"id,omitempty"`
		Type         string `json:"type,omitempty"`
		Controller   string `json:"controller,omitempty"`
		PublicKeyHex string `json:"publicKeyHex,omitempty"`
	}
	// Doc is the DID document struct
	Doc struct {
		Context        string           `json:"@context,omitempty"`
		ID             string           `json:"id,omitempty"`
		Authentication []authentication `json:"authentication,omitempty"`
	}
)

// CreateDID creates a new DID using public key
func CreateDID(pk crypto.PublicKey) Doc {
	id := DIDPrefix + "0x" + hex.EncodeToString(pk.Hash())
	return Doc{
		Context: DIDContext,
		ID:      id,
		Authentication: []authentication{
			authentication{
				ID:           id + DIDOwner,
				Type:         DIDAuthType,
				Controller:   id,
				PublicKeyHex: pk.HexString(),
			},
		},
	}
}
