// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package did

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

func TestCreateDID(t *testing.T) {
	r := require.New(t)

	a, err := account.NewAccount()
	r.NoError(err)
	id := DIDPrefix + "0x" + hex.EncodeToString(a.PublicKey().Hash())

	d := CreateDID(a)
	r.Equal(id, d.ID)
	r.Equal(1, len(d.Authentication))
	r.Equal(id+DIDOwner, d.Authentication[0].ID)
	r.Equal(DIDAuthType, d.Authentication[0].Type)
	r.Equal(id, d.Authentication[0].Controller)
	r.Equal(hex.EncodeToString(a.PublicKey().Bytes()), d.Authentication[0].PublicKeyHex)
}
