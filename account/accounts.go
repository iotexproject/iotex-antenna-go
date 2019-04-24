// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import "github.com/pkg/errors"

// Accounts type
type Accounts struct {
	accounts map[string]Account
}

// NewAccounts return Accounts instance
func NewAccounts() *Accounts {
	accounts := make(map[string]Account)
	return &Accounts{
		accounts: accounts,
	}
}

// Create new account
func (acts *Accounts) Create() (Account, error) {
	acc, err := NewAccount()
	if err != nil {
		return nil, err
	}
	acts.accounts[acc.Address()] = acc
	return acc, nil
}

// GetAccount by address
func (acts *Accounts) GetAccount(addr string) (Account, error) {
	acc, ok := acts.accounts[addr]
	if !ok {
		return nil, errors.Errorf("Account %s does not exist", addr)
	}
	return acc, nil
}

// PrivateKeyToAccount new Account by privateKey
func (acts *Accounts) PrivateKeyToAccount(privateKey string) (Account, error) {
	acc, err := NewAccountFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	acts.accounts[acc.Address()] = acc
	return acc, nil
}
