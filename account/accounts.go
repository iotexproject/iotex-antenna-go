// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"github.com/iotexproject/iotex-core/pkg/keypair"
)

type Accounts struct {
	acts []Account
}

func (acts *Accounts) Create() (Account, error) {
	private, err := keypair.GenerateKey()
	if err != nil {
		return Account{}, err
	}
	return privateToAccount(private.EcdsaPrivateKey())
}
func (acts *Accounts) GetAccount(addr string) (string, bool) {
	for _, acc := range acts.acts {
		if acc.Address == addr {
			return acc.PrivateKey, true
		}
	}
	return "", false
}
func (acts *Accounts) AddAccount(privateKey string) error {
	acc, err := acts.PrivateKeyToAccount(privateKey)
	if err != nil {
		return err
	}
	_, exist := acts.GetAccount(acc.Address)
	if !exist {
		acts.acts = append(acts.acts, acc)
	}
	return nil
}
func (acts *Accounts) PrivateKeyToAccount(privateKey string) (Account, error) {
	private, err := keypair.HexStringToPrivateKey(privateKey)
	if err != nil {
		return Account{}, nil
	}
	return privateToAccount(private.EcdsaPrivateKey())
}

func (acts *Accounts) Sign(data []byte, privateKey string) ([]byte, error) {
	act, err := acts.PrivateKeyToAccount(privateKey)
	if err != nil {
		return nil, err
	}
	return act.Sign(data)
}
