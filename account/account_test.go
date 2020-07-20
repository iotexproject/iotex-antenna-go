// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"encoding/hex"
	"testing"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/stretchr/testify/assert"
)

var (
	Address    = "io187wzp08vnhjjpkydnr97qlh8kh0dpkkytfam8j"
	PrivateKey = "0806c458b262edd333a191e92f561aff338211ee3e18ab315a074a2d82aa343f"
	PublicKey  = "044e18306ae9ef4ec9d07bf6e705442d4d1a75e6cdf750330ca2d880f2cc54607c9c33deb9eae9c06e06e04fe9ce3d43962cc67d5aa34fbeb71270d4bad3d648d9"
)

const text = "IoTeX is the auto-scalable and privacy-centric blockchain."

func TestHash160b(t *testing.T) {
	h := hash.Hash160b([]byte(text))
	assert.Equal(t, "93988dc3d2d879f703c7d3f54dcc1b473b27d015", hex.EncodeToString(h[:]))
}

func TestHash256b(t *testing.T) {
	h := hash.Hash256b([]byte(text))
	assert.Equal(t, "aada23f93a5ed1829ebf1c0693988dc3d2d879f703c7d3f54dcc1b473b27d015", hex.EncodeToString(h[:]))
}

func TestAccount(t *testing.T) {
	assert := assert.New(t)

	act, err := HexStringToAccount(PrivateKey)
	assert.NoError(err)
	assert.Equal(Address, act.Address().String())
	assert.Equal(PublicKey, act.PrivateKey().PublicKey().HexString())

	act1, err := PrivateKeyToAccount(act.PrivateKey())
	assert.NoError(err)
	assert.Equal(act, act1)

	b, err := act.Sign([]byte(text))
	assert.NoError(err)
	assert.Equal(
		"482da72c8faa48ee1ac2cf9a5f9ecd42ee3258be5ddd8d6b496c7171dc7bfe8e75e5d16e7129c88d99a21a912e5c082fa1baab6ba87d2688ebd7d27bb1ab090701",
		hex.EncodeToString(b),
	)
	// verify the signature
	assert.True(act.Verify([]byte(text), b))

	act.Zero()
	b, err = act.Sign([]byte(text))
	assert.Equal("invalid private key", err.Error())
}

func TestHashMessage(t *testing.T) {
	assert := assert.New(t)

	act, err := HexStringToAccount(PrivateKey)
	assert.NoError(err)

	h := act.HashMessage([]byte("hello"))
	assert.Equal(
		"5077b388a631936d73d9c6c9a0bf6016843a8b594540d1d968f7ea40d1541c58",
		hex.EncodeToString(h[:]),
	)
}

func TestSignMessage(t *testing.T) {
	assert := assert.New(t)

	act, err := HexStringToAccount(PrivateKey)
	assert.NoError(err)

	b, err := act.SignMessage([]byte("hello"))
	assert.NoError(err)
	assert.Equal(
		"f09c729cc8617aeda344defba6c0eb0eb3ee71732e26f22d1a9fac5beeaa86da3a368417e31779b44e3df4440dfec89a9ecb40567b60228efb67c79672288cef01",
		hex.EncodeToString(b),
	)
}

func TestRecover(t *testing.T) {
	assert := assert.New(t)

	h, _ := hex.DecodeString("5077b388a631936d73d9c6c9a0bf6016843a8b594540d1d968f7ea40d1541c58")
	sig, _ := hex.DecodeString("f09c729cc8617aeda344defba6c0eb0eb3ee71732e26f22d1a9fac5beeaa86da3a368417e31779b44e3df4440dfec89a9ecb40567b60228efb67c79672288cef01")

	addr, err := Recover(h, sig)
	assert.NoError(err)
	assert.Equal(Address, addr.String())
}
