// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"testing"
	"encoding/hex"
	"fmt"

	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/stretchr/testify/assert"

)

var testAcct = Account{
	Address:    "io187wzp08vnhjjpkydnr97qlh8kh0dpkkytfam8j",
	PrivateKey: "0806c458b262edd333a191e92f561aff338211ee3e18ab315a074a2d82aa343f",
	PublicKey:  "044e18306ae9ef4ec9d07bf6e705442d4d1a75e6cdf750330ca2d880f2cc54607c9c33deb9eae9c06e06e04fe9ce3d43962cc67d5aa34fbeb71270d4bad3d648d9",
}

const text = "IoTeX is the auto-scalable and privacy-centric blockchain."

func TestHash160b(t *testing.T) {
	h := hash.Hash160b([]byte(text))
	assert.Equal(t, "93988dc3d2d879f703c7d3f54dcc1b473b27d015", hex.EncodeToString(h[:]))
}

func TestHash256b(t *testing.T) {
	h := hash.Hash256b([]byte(text))
	assert.Equal(t, "aada23f93a5ed1829ebf1c0693988dc3d2d879f703c7d3f54dcc1b473b27d015", hex.EncodeToString(h[:]))
}

func TestAccount_Sign(t *testing.T) {
	b, err := testAcct.Sign([]byte(text))
	assert.NoError(t, err)
	assert.Equal(t,
		"482da72c8faa48ee1ac2cf9a5f9ecd42ee3258be5ddd8d6b496c7171dc7bfe8e75e5d16e7129c88d99a21a912e5c082fa1baab6ba87d2688ebd7d27bb1ab090701",
		fmt.Sprintf("%x", b),
	)
}
