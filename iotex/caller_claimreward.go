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
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type claimRewardCaller struct {
	*sendActionCaller
	amount *big.Int
}

func (c *claimRewardCaller) SetData(data []byte) ClaimRewardCaller {
	c.sendActionCaller.setPayload(data)
	return c
}

func (c *claimRewardCaller) SetGasLimit(g uint64) ClaimRewardCaller {
	c.sendActionCaller.setGasLimit(g)
	return c
}

func (c *claimRewardCaller) SetGasPrice(g *big.Int) ClaimRewardCaller {
	c.sendActionCaller.setGasPrice(g)
	return c
}

func (c *claimRewardCaller) SetNonce(n uint64) ClaimRewardCaller {
	c.sendActionCaller.setNonce(n)
	return c
}

func (c *claimRewardCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("claim amount cannot be nil", errcodes.InvalidParam)
	}

	tx := iotextypes.ClaimFromRewardingFund{
		Amount: c.amount.String(),
		Data:   c.payload,
	}
	c.core = &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Action:  &iotextypes.ActionCore_ClaimFromRewardingFund{ClaimFromRewardingFund: &tx},
	}
	return c.sendActionCaller.Call(ctx, opts...)
}
