// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccounts_Create(t *testing.T) {
	assert := assert.New(t)

	acts := NewAccounts()
	act, err := acts.Create()
	assert.NoError(err)
	assert.NotEmpty(act.PrivateKey())
	assert.NotEmpty(act.PublicKey())
	assert.NotEmpty(act.Address)

	b, err := acts.GetAccount(act.Address())
	assert.NoError(err)
	assert.Equal(act, b)
}

func TestAccounts_PrivateKeyToAccount(t *testing.T) {
	assert := assert.New(t)

	acts := NewAccounts()
	act, err := acts.PrivateKeyToAccount(PrivateKey)
	assert.NoError(err)

	b, err := acts.GetAccount(act.Address())
	assert.NoError(err)
	assert.Equal(act, b)
}
