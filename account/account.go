// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"crypto/ecdsa"

	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
)

type Account struct {
	Address    string
	PrivateKey string
	PublicKey  string
}

func (act Account) Sign(data []byte) ([]byte, error) {
	priv, err := keypair.HexStringToPrivateKey(act.PrivateKey)
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(data)
	return priv.Sign(h[:])
}

func privateToAccount(private *ecdsa.PrivateKey) (acc Account, err error) {
	pri, err := keypair.BytesToPrivateKey(private.D.Bytes())
	if err != nil {
		return
	}
	addr, err := address.FromBytes(pri.PublicKey().Hash())
	if err != nil {
		return
	}
	return Account{
		Address:    addr.String(),
		PrivateKey: pri.HexString(),
		PublicKey:  pri.PublicKey().HexString(),
	}, nil
}
