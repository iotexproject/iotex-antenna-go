package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccounts_Create(t *testing.T) {
	acts := Accounts{}
	act, err := acts.Create()
	assert.NoError(t, err)
	assert.NotEmpty(t, act.PrivateKey)
	assert.NotEmpty(t, act.PublicKey)
	assert.NotEmpty(t, act.Address)
}

func TestAccounts_PrivateKeyToAccount(t *testing.T) {
	acts := Accounts{}
	act, _ := acts.PrivateKeyToAccount(expAct.PrivateKey)
	assert.Equal(t, expAct, act)
}
