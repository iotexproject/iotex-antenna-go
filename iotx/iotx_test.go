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
	iotx, err := NewIotx(host)
	require.NoError(err)
	err = iotx.Accounts.AddAccount(accountPrivateKey)
	require.NoError(err)

	req := &TransferRequest{From: accountAddress, To: to, Value: "1000000000000000000", Payload: "", GasLimit: "1000000", GasPrice: "1000000000000"}
	err = iotx.SendTransfer(req)
	require.NoError(err)
}
