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
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type transferCaller struct {
	account   account.Account
	api       iotexapi.APIServiceClient
	amount    *big.Int
	recipient address.Address
	payload   []byte
	gasLimit  *uint64
	gasPrice  *big.Int
	nonce     *uint64
}

func (c *transferCaller) SetPayload(pl []byte) TransferCaller {
	if pl == nil {
		return c
	}
	c.payload = make([]byte, len(pl))
	copy(c.payload, pl)
	return c
}

func (c *transferCaller) SetGasLimit(g uint64) TransferCaller {
	c.gasLimit = &g
	return c
}

func (c *transferCaller) SetGasPrice(g *big.Int) TransferCaller {
	c.gasPrice = g
	return c
}

func (c *transferCaller) SetNonce(n uint64) TransferCaller {
	c.nonce = &n
	return c
}

func (c *transferCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *transferCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("transfer amount cannot be nil", errcodes.InvalidParam)
	}

	tx := &iotextypes.Transfer{
		Amount:    c.amount.String(),
		Recipient: c.recipient.String(),
		Payload:   c.payload,
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
