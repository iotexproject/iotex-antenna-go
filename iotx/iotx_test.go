// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testnet           = "api.testnet.iotex.one:80"
	mainnet           = "api.iotex.one:443"
	accountPrivateKey = "9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	to                = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
)

func TestTransfer(t *testing.T) {
	require := require.New(t)
	iotx, err := NewIotx(testnet, false)
	require.NoError(err)
	defer iotx.Close()
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)

	req := &TransferRequest{
		From:     acc.Address(),
		To:       to,
		Value:    "1000000000000000000",
		Payload:  "",
		GasLimit: "1000000",
		GasPrice: "1",
	}

	hash, err := iotx.SendTransfer(req)
	require.NoError(err)
	require.NotEmpty(hash)
}

func TestDeployContract(t *testing.T) {
	require := require.New(t)
	iotx, err := NewIotx(testnet, false)
	require.NoError(err)
	defer iotx.Close()
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)

	req := &ContractRequest{
		From:     acc.Address(),
		Abi:      `[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`,
		Data:     "608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a723058208d4f6c9737f34d9b28ef070baa8127c0876757fbf6f3945a7ea8d4387ca156590029",
		GasLimit: "1000000",
		GasPrice: "1",
	}

	hash, err := iotx.DeployContract(req, big.NewInt(10))

	require.NoError(err)
	require.NotNil(hash)
}

func TestReadContract(t *testing.T) {
	// TODO: re-enable test when new gPRC API has been pushed to testnet/mainnet
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	require := require.New(t)
	iotx, err := NewIotx(mainnet, true)
	require.NoError(err)
	defer iotx.Close()
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)

	req := &ContractRequest{
		From:     acc.Address(),
		Address:  "io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0",
		Abi:      `[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`,
		Method:   "get",
		GasLimit: "1000000",
		GasPrice: "1",
	}

	result, err := iotx.ReadContractByMethod(req)
	require.NoError(err)
	require.NotNil(result)
}

func TestExecuteContract(t *testing.T) {
	require := require.New(t)
	iotx, err := NewIotx(testnet, false)
	require.NoError(err)
	defer iotx.Close()
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)
	require.NotNil(acc)

	req := &ContractRequest{
		From:     acc.Address(),
		Address:  "io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0",
		Abi:      `[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`,
		Method:   "set",
		Amount:   "0",
		GasLimit: "1000000",
		GasPrice: "1",
	}

	result, err := iotx.ExecuteContract(req, big.NewInt(8))

	require.NoError(err)
	require.NotNil(result)
}

func TestReadContractByHash(t *testing.T) {
	// TODO: re-enable test when new gPRC API has been pushed to testnet/mainnet
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	require := require.New(t)
	iotx, err := NewIotx(mainnet, true)
	require.NoError(err)
	defer iotx.Close()

	result, err := iotx.ReadContractByHash("edf65e7ccbfb05e4fbd394db1acc276029c309994879e3a8c07023a753ea8886")
	require.NoError(err)
	require.NotNil(result)
}
