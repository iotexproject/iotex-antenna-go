// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package antenna

import (
	"math/big"

	"github.com/iotexproject/iotex-antenna-go/iotx"
	"github.com/iotexproject/iotex-antenna-go/utils"
)

// Antenna
type Antenna = iotx.Iotx

// NewRPCMethod returns RPCMethod interacting with endpoint
func NewRPCMethod(endpoint string) (*Antenna, error) {
	return iotx.NewRPCMethod(endpoint)
}

// FromRau is a function to convert Rau to Iotx.
func FromRau(rau *big.Int) int64 {
	return utils.FromRau(rau)
}

// ToRau is a function to convert various units to Rau.
func ToRau(iotxs int64) *big.Int {
	return utils.ToRau(iotxs)
}
