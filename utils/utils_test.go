// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromRau(t *testing.T) {
	require := require.New(t)
	convert := FromRau("1000", "Jin")
	require.Equal("1", convert)

	convert = FromRau("1000000000000000000", "Rau")
	require.Equal("1", convert)
}
func TestToRau(t *testing.T) {
	require := require.New(t)
	convert := ToRau("1", "Iotx")
	require.Equal("1000000000000000000", convert)

	convert = ToRau("1", "GRau")
	require.Equal("1000000000", convert)
}
