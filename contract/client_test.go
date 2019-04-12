// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// use accountPrivateKey="9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	// use accountAddress="io14gnqxf9dpkn05g337rl7eyt2nxasphf5m6n0rd"
	host = "api.iotex.one:80"
)

func TestServer_Deploy(t *testing.T) {
	require := require.New(t)
	sct, err := NewSmartContract("testdata/array-return.json", host)
	require.NoError(err)
	err = sct.DeployContracts()
	require.NoError(err)
}
