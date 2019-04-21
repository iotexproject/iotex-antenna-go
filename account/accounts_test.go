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
	acts := NewAccounts()
	act, err := acts.Create()
	assert.NoError(t, err)
	assert.NotEmpty(t, act.PrivateKey)
	assert.NotEmpty(t, act.PublicKey)
	assert.NotEmpty(t, act.Address)
}

func TestAccounts_PrivateKeyToAccount(t *testing.T) {
	acts := NewAccounts()
	act, _ := acts.PrivateKeyToAccount(testAcct.PrivateKey)
	assert.Equal(t, testAcct.Address, act.Address)
	assert.Equal(t, testAcct.PrivateKey, act.PrivateKey)
	assert.Equal(t, testAcct.PublicKey, act.PublicKey)
}
