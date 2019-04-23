// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	host              = "api.testnet.iotex.one:80"
	accountPrivateKey = "9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	accountAddress    = "io14gnqxf9dpkn05g337rl7eyt2nxasphf5m6n0rd"
	to                = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
)

func TestTransfer(t *testing.T) {
	require := require.New(t)
	iotx, err := New(host)
	require.NoError(err)
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)
	require.EqualValues(acc.Address, accountAddress)

	req := &TransferRequest{
		From:     accountAddress,
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
	iotx, err := New(host)
	require.NoError(err)
	acc, err := iotx.Accounts.PrivateKeyToAccount(accountPrivateKey)
	require.NoError(err)
	require.EqualValues(acc.Address, accountAddress)

	req := &ContractRequest{
		From:     accountAddress,
		Data:     "608060405234801561001057600080fd5b5060bf8061001f6000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582071e7563748698bc5d23ab47a5d532b099d97c16ae8ff555cd22d25f6951582850029",
		GasLimit: "1000000",
		GasPrice: "1",
	}

	hash, err := iotx.DeployContract(req)
	require.NoError(err)
	require.NotNil(hash)
}
