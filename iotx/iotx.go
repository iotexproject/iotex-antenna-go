// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
)

// RPCMethod provides simple interface tp invoke rpc method
type Iotx = rpcmethod.RPCMethod

// NewRPCMethod returns RPCMethod interacting with endpoint
func NewRPCMethod(endpoint string) (*Iotx, error) {
	return rpcmethod.NewRPCMethod(endpoint)
}
