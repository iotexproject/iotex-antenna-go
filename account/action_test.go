// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/keypair"
)

func TestAction_Envelope(t *testing.T) {
	gasLimit := uint64(888)
	gasPrice := big.NewInt(999)
	tx, err := action.NewTransfer(123, big.NewInt(456),
		testAcct.Address, []byte("hello world!"), gasLimit, gasPrice)
	assert.NoError(t, err)

	builder := &action.EnvelopeBuilder{}
	elp := builder.
		SetAction(tx).
		SetGasLimit(gasLimit).
		SetGasLimit(gasPrice.Uint64()).
		Build()

	privKey, err := keypair.HexStringToPrivateKey(testAcct.PrivateKey)
	assert.NoError(t, err)

	sealed, err := action.Sign(elp, privKey)
	actionPb := sealed.Proto()
	assert.Equal(
		t,
		hex.EncodeToString(actionPb.Signature),
		"bc64e8689e29c8afd31c10fad9b4973d94842db77b4dd32f026d3887849d057c1b7e92e2b5a6b388aeeeb852f16c3584bc07d5124d09ef169774d9c8bff30bdd00",
	)
}

func TestAction_SerializationDeserialization(t *testing.T) {
	gasLimit := uint64(888)
	gasPrice := big.NewInt(999)
	action, err := action.NewTransfer(
		123,
		big.NewInt(456),
		testAcct.Address,
		[]byte("hello world!"),
		gasLimit,
		gasPrice,
	)
	assert.NoError(t, err)
	actionPb := action.Proto()
	marshaled, err := proto.Marshal(actionPb)
	assert.NoError(t, err)
	assert.Equal(
		t,
		hex.EncodeToString(marshaled),
		"0a033435361229696f313837777a703038766e686a6a706b79646e723937716c68386b683064706b6b797466616d386a1a0c68656c6c6f20776f726c6421",
	)
}
