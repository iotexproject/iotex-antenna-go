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
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type transferCaller struct {
	*sendActionCaller
	amount    *big.Int
	recipient address.Address
}

func (c *transferCaller) SetPayload(pl []byte) SendActionCaller {
	c.sendActionCaller.setPayload(pl)
	return c
}

func (c *transferCaller) SetGasLimit(g uint64) SendActionCaller {
	c.sendActionCaller.setGasLimit(g)
	return c
}

func (c *transferCaller) SetGasPrice(g *big.Int) SendActionCaller {
	c.sendActionCaller.setGasPrice(g)
	return c
}

func (c *transferCaller) SetNonce(n uint64) SendActionCaller {
	c.sendActionCaller.setNonce(n)
	return c
}

func (c *transferCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("transfer amount cannot be nil", errcodes.InvalidParam)
	}

	tx := iotextypes.Transfer{
		Amount:    c.amount.String(),
		Recipient: c.recipient.String(),
		Payload:   c.payload,
	}
	c.core = &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Action:  &iotextypes.ActionCore_Transfer{Transfer: &tx},
	}
	return c.sendActionCaller.Call(ctx, opts...)
}
