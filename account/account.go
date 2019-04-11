// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/pkg/hash"
)

type Account struct {
	Address string
	PublicKey string
	PrivateKey string
}

func (act Account) Sign(data []byte) ([]byte, error) {
	priv, err := keypair.DecodePrivateKey(act.PrivateKey)
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(data)
	fmt.Printf("h: %+x\n", h)
	return crypto.Sign(h[:], priv)
}

func privateToAccount(private *ecdsa.PrivateKey) (Account, error) {
	pkHash := keypair.HashPubKey(&private.PublicKey)
	addr, _ := address.FromBytes(pkHash[:])
	priKeyBytes := keypair.PrivateKeyToBytes(private)
	pubKeyBytes := keypair.PublicKeyToBytes(&private.PublicKey)
	return Account{
		Address: addr.String(),
		PublicKey: fmt.Sprintf("%x", pubKeyBytes),
		PrivateKey: fmt.Sprintf("%x", priKeyBytes),
	}, nil
}
