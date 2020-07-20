// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"fmt"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
)

type (
	// Account is a user account
	Account interface {
		// Address returns the IoTeX address
		Address() address.Address
		// PrivateKey returns the embedded private key interface
		PrivateKey() crypto.PrivateKey
		// PublicKey returns the embedded public key interface
		PublicKey() crypto.PublicKey
		// Sign signs the message using the private key
		Sign([]byte) ([]byte, error)
		// Verify verifies the message using the public key
		Verify([]byte, []byte) bool
		// Zero zeroes the private key data
		Zero()
		// SignMessage signs the message using preamble
		SignMessage(data []byte) ([]byte, error)
	}

	account struct {
		private crypto.PrivateKey
		address address.Address
	}
)

// NewAccount generates a new account
func NewAccount() (Account, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr, err := address.FromBytes(pk.PublicKey().Hash())
	if err != nil {
		return nil, err
	}
	return &account{
		pk,
		addr,
	}, nil
}

// HexStringToAccount generates an account from private key string
func HexStringToAccount(privateKey string) (Account, error) {
	sk, err := crypto.HexStringToPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	addr, err := address.FromBytes(sk.PublicKey().Hash())
	if err != nil {
		return nil, err
	}
	return &account{
		sk,
		addr,
	}, nil
}

// PrivateKeyToAccount generates an account from an existing private key interface
func PrivateKeyToAccount(key crypto.PrivateKey) (Account, error) {
	addr, err := address.FromBytes(key.PublicKey().Hash())
	if err != nil {
		return nil, err
	}
	return &account{
		key,
		addr,
	}, nil
}

// Address returns the IoTeX address
func (act *account) Address() address.Address {
	return act.address
}

// PrivateKey return the embedded private key
func (act *account) PrivateKey() crypto.PrivateKey {
	return act.private
}

// PublicKey returns the embedded public key interface
func (act *account) PublicKey() crypto.PublicKey {
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

// SignMessage signs the message using preamble
func (act *account) SignMessage(data []byte) ([]byte, error) {
	h := HashMessage(data)
	return act.private.Sign(h[:])
}

// HashMessage hash the message using preamble
func HashMessage(data []byte) hash.Hash256 {
	preamble := fmt.Sprintf("\x16IoTeX Signed Message:\n%d", len(data))
	iotexMessage := []byte(preamble)
	iotexMessage = append(iotexMessage, data...)
	return hash.Hash256b(iotexMessage)
}

// RecoverAddress recover address by message hash and signature
func RecoverAddress(messageHash, signature []byte) (address.Address, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("wrong size for signature: got %d, want 65", len(signature))
	}

	pub, err := ethCrypto.Ecrecover(messageHash, signature)
	if err != nil {
		return nil, err
	}

	payload := hash.Hash160b(pub[1:])
	return address.FromBytes(payload[:])
}
