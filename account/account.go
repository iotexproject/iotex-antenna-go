// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"crypto/ecdsa"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
)

// Account type
type Account struct {
	private    *keypair.PrivateKey
	Address    string
	PrivateKey string
	PublicKey  string
}

// Sign by acount private key
func (act *Account) Sign(data []byte) ([]byte, error) {
	priv, err := keypair.HexStringToPrivateKey(act.PrivateKey)
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(data)
	return priv.Sign(h[:])
}

// Private return keypair private key
func (act *Account) Private() *keypair.PrivateKey {
	return act.private
}

// FromPrivateKey create Account from private key string
func FromPrivateKey(privateKey string) (*Account, error) {
	private, err := keypair.HexStringToPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return privateToAccount(private.EcdsaPrivateKey())
}

func privateToAccount(privateKey *ecdsa.PrivateKey) (acc *Account, err error) {
	pri, err := keypair.BytesToPrivateKey(privateKey.D.Bytes())
	if err != nil {
		return
	}
	addr, err := address.FromBytes(pri.PublicKey().Hash())
	if err != nil {
		return
	}
	return &Account{
		private:    &pri,
		Address:    addr.String(),
		PrivateKey: pri.HexString(),
		PublicKey:  pri.PublicKey().HexString(),
	}, nil
}
