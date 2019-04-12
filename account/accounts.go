// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/ethereum/go-ethereum/crypto"
)

type Accounts struct {
	acts []Account
}

func (acts *Accounts) Create() (Account, error) {
	private, err := crypto.GenerateKey()
	if err != nil {
		return Account{}, err
	}
	return privateToAccount(private)
}

func (acts *Accounts) PrivateKeyToAccount(privateKey string) (Account, error) {
	private, err := keypair.DecodePrivateKey(privateKey)
	if err != nil {
		return Account{}, nil
	}

	return privateToAccount(private)
}

func (acts *Accounts) Sign(data []byte, privateKey string) ([]byte, error) {
	act, err := acts.PrivateKeyToAccount(privateKey)
	if err != nil {
		return nil, err
	}

	return act.Sign(data)
}
