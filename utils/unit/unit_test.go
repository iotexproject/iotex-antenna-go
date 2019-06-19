// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package unit

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromRau(t *testing.T) {
	require := require.New(t)
	n, _ := big.NewInt(0).SetString("1000", 10)
	convert := FromRau(n, "KRau")
	require.Equal("1", convert.Text(10))

	n, _ = big.NewInt(0).SetString("1000000000000000000", 10)
	convert = FromRau(n, "Iotx")
	require.Equal("1", convert.Text(10))
}
func TestToRau(t *testing.T) {
	require := require.New(t)
	n, _ := big.NewInt(0).SetString("1", 10)
	convert := ToRau(n, "Iotx")
	require.Equal("1000000000000000000", convert.Text(10))

	convert = ToRau(n, "GRau")
	require.Equal("1000000000", convert.Text(10))
}
