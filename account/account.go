// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
)

type (
	// Account is a user account
	Account interface {
		// Address returns the IoTeX address
		Address() string
		// PrivateKey returns the embedded private key interface
		PrivateKey() keypair.PrivateKey
		// PublicKey returns the embedded public key interface
		PublicKey() keypair.PublicKey
		// Sign signs the message using the private key
		Sign([]byte) ([]byte, error)
		// Verify verifies the message using the public key
		Verify([]byte, []byte) bool
		// Zero zeroes the private key data
		Zero()
	}

	account struct {
		private keypair.PrivateKey
		address string
	}
)

// NewAccount generates a new account
func NewAccount() (Account, error) {
	pk, err := keypair.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr, err := address.FromBytes(pk.PublicKey().Hash())
	if err != nil {
		return nil, err
	}
	return &account{
		pk,
		addr.String(),
	}, nil
}

// NewAccountFromPrivateKey generates an account from private key string
func NewAccountFromPrivateKey(privateKey string) (Account, error) {
	pk, err := keypair.HexStringToPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	addr, err := address.FromBytes(pk.PublicKey().Hash())
	if err != nil {
		return nil, err
	}
	return &account{
		pk,
		addr.String(),
	}, nil
}

// Address returns the IoTeX address
func (act *account) Address() string {
	return act.address
}

// PrivateKey return the embedded private key
func (act *account) PrivateKey() keypair.PrivateKey {
	return act.private
}

// PublicKey returns the embedded public key interface
func (act *account) PublicKey() keypair.PublicKey {
	return act.private.PublicKey()
}

// Sign signs the message using the private key
func (act *account) Sign(data []byte) ([]byte, error) {
	h := hash.Hash256b(data)
	return act.private.Sign(h[:])
}

// Verify verifies the message using the public key
func (act *account) Verify(data []byte, sig []byte) bool {
	h := hash.Hash256b(data)
	return act.PublicKey().Verify(h[:], sig)
}

// Zero zeroes the private key data
func (act *account) Zero() {
	act.private.Zero()
}
