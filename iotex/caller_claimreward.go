// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"math/big"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type claimRewardCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	amount   *big.Int
	data     []byte
	gasLimit *uint64
	gasPrice *big.Int
	nonce    *uint64
}

func (c *claimRewardCaller) SetData(data []byte) ClaimRewardCaller {
	if data == nil {
		return c
	}
	c.data = make([]byte, len(data))
	copy(c.data, data)
	return c
}

func (c *claimRewardCaller) SetGasLimit(g uint64) ClaimRewardCaller {
	c.gasLimit = &g
	return c
}

func (c *claimRewardCaller) SetGasPrice(g *big.Int) ClaimRewardCaller {
	c.gasPrice = g
	return c
}

func (c *claimRewardCaller) SetNonce(n uint64) ClaimRewardCaller {
	c.nonce = &n
	return c
}

func (c *claimRewardCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *claimRewardCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("claim amount cannot be nil", errcodes.InvalidParam)
	}

	tx := &iotextypes.ClaimFromRewardingFund{
		Amount: c.amount.String(),
		Data:   c.data,
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		nonce:    c.nonce,
		action:   tx,
	}
	return sc.Call(ctx, opts...)
}
