// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	require := require.New(t)

	tsf, err := NewTransfer(1, 100, big.NewInt(-1), big.NewInt(10), "test", nil)
	require.Error(err)
	require.Nil(tsf)
	tsf, err = NewTransfer(1, 100, big.NewInt(1), big.NewInt(-10), "test", nil)
	require.Error(err)
	require.Nil(tsf)
	tsf, err = NewTransfer(1, 100, big.NewInt(1), big.NewInt(10), "", nil)
	require.Error(err)
	require.Nil(tsf)
	tsf, err = NewTransfer(1, 100, big.NewInt(1), big.NewInt(10), "test", nil)
	require.NoError(err)
	require.NotNil(tsf)
}
