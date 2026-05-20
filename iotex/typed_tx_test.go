// Copyright (c) 2026 IoTeX
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
)

// typedTxKey is read from IOTEX_INTEGRATION_KEY (a funded testnet private key
// hex string). The tests in this file are skipped when it's empty so that
// `go test ./...` does not require network access or testnet funds.
func typedTxKey(t *testing.T) string {
	t.Helper()
	k := os.Getenv("IOTEX_INTEGRATION_KEY")
	if k == "" {
		t.Skip("set IOTEX_INTEGRATION_KEY to a funded testnet private key to run typed-tx integration tests")
	}
	return k
}

func newTypedTxClient(t *testing.T) (AuthedClient, account.Account) {
	t.Helper()
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	acc, err := account.HexStringToAccount(typedTxKey(t))
	require.NoError(t, err)
	return NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc), acc
}

func TestTypedTransfer_AccessList(t *testing.T) {
	c, _ := newTypedTxClient(t)
	to, err := address.FromString(_to)
	require.NoError(t, err)

	hash, err := c.Transfer(to, big.NewInt(1)).
		SetTxType(accessListTxType).
		SetGasPrice(big.NewInt(int64(unit.Qev))).
		SetGasLimit(50000).
		SetAccessList(types.AccessList{}).
		Call(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	t.Logf("AccessList transfer hash: %x", hash[:])
}

func TestTypedTransfer_DynamicFee(t *testing.T) {
	c, _ := newTypedTxClient(t)
	to, err := address.FromString(_to)
	require.NoError(t, err)

	tip := big.NewInt(int64(unit.Qev))
	feeCap := new(big.Int).Mul(big.NewInt(2), tip)

	hash, err := c.Transfer(to, big.NewInt(1)).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(tip).
		SetGasFeeCap(feeCap).
		SetGasLimit(50000).
		Call(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, hash)
	t.Logf("DynamicFee transfer hash: %x", hash[:])
}

func TestTypedExecute_DynamicFee(t *testing.T) {
	c, _ := newTypedTxClient(t)
	parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`))
	require.NoError(t, err)
	contract, err := address.FromString("io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0")
	require.NoError(t, err)

	tip := big.NewInt(int64(unit.Qev))
	feeCap := new(big.Int).Mul(big.NewInt(2), tip)

	hash, err := c.Contract(contract, parsedABI).Execute("set", big.NewInt(int64(time.Now().Unix()))).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(tip).
		SetGasFeeCap(feeCap).
		SetGasLimit(1000000).
		Call(context.Background())
	require.NoError(t, err)
	require.NotNil(t, hash)
	t.Logf("DynamicFee execute hash: %x", hash[:])
}

// TestTypedTransfer_RequiresGasLimit confirms that the SDK rejects a typed tx
// when GasLimit is unset. Validation runs before any RPC; nonce is set to
// skip the GetAccount call so the nil API client never gets touched.
func TestTypedTransfer_RequiresGasLimit(t *testing.T) {
	acc, err := account.HexStringToAccount("73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426")
	require.NoError(t, err)
	to, err := address.FromString(_to)
	require.NoError(t, err)

	c := NewAuthedClient(nil, 2, acc)
	_, err = c.Transfer(to, big.NewInt(1)).
		SetNonce(1).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(big.NewInt(1)).
		SetGasFeeCap(big.NewInt(2)).
		Call(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "gas limit must be set explicitly")
}

// TestTypedExecute_DynamicFeeMissingFeeCap confirms validation when a typed
// tx is missing its required fee cap.
func TestTypedExecute_DynamicFeeMissingFeeCap(t *testing.T) {
	acc, err := account.HexStringToAccount("73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426")
	require.NoError(t, err)
	contractAddr, err := address.FromString("io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0")
	require.NoError(t, err)
	parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	require.NoError(t, err)

	c := NewAuthedClient(nil, 2, acc)
	_, err = c.Contract(contractAddr, parsedABI).Execute("set", big.NewInt(1)).
		SetNonce(1).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(big.NewInt(1)).
		SetGasLimit(1000000).
		Call(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "gas tip cap and fee cap must be set")
}
