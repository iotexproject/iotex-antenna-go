// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package utils

import "math/big"

// FromRau is a function to convert Rau to Iotx.
func FromRau(rau *big.Int) int64 {
	rau = rau.Div(rau, big.NewInt(1e18))
	return rau.Int64()
}

// ToRau is a function to convert various units to Rau.
func ToRau(iotx int64) *big.Int {
	itx := big.NewInt(iotx)
	return itx.Mul(itx, big.NewInt(1e18))
}
