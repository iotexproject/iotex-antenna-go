// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package antenna

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	host = "api.testnet.iotex.one:80"
)

func TestServer_GetAccount(t *testing.T) {
	require := require.New(t)
	_, err := NewAntenna(host)
	require.NoError(err)
}
