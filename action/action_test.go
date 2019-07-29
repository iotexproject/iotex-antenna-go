// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/stretchr/testify/assert"

	"github.com/iotexproject/iotex-antenna-go/account"
)

var (
	PrivateKey = "0806c458b262edd333a191e92f561aff338211ee3e18ab315a074a2d82aa343f"
)

func TestActionTransfer(t *testing.T) {
	assert := assert.New(t)

	testAcct, err := account.HexStringToAccount(PrivateKey)
	assert.NoError(err)
	ac, err := NewTransfer(123, uint64(888), big.NewInt(999), big.NewInt(456), testAcct.Address(), []byte("hello world!"))
	assert.NoError(err)
	sac, err := ac.Sign(testAcct)
	assert.NoError(err)
	assert.Equal(
		"555cc8af4181bf85c044c3201462eeeb95374f78aa48c67b87510ee63d5e502372e53082f03e9a11c1e351de539cedf85d8dff87de9d003cb9f92243541541a000",
		hex.EncodeToString(sac.Signature),
	)
	marshaled, err := proto.Marshal(ac)
	assert.NoError(err)
	h := hash.Hash256b(marshaled)
	assert.Equal("0f17cd7f43bdbeff73dfe8f5cb0c0045f2990884e5050841de887cf22ca35b50", hex.EncodeToString(h[:]))
	assert.True(testAcct.Verify(marshaled, sac.Signature))
	marshaled, err = proto.Marshal(sac)
	assert.NoError(err)
	h = hash.Hash256b(marshaled)
	assert.Equal("6c84ac119058e859a015221f87a4e187c393d0c6ee283959342eac95fad08c33", hex.EncodeToString(h[:]))
	assert.False(testAcct.Verify(marshaled, sac.Signature))
}