// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

// ProtocolVersion is the iotex protocol version to use. Currently 1.
const ProtocolVersion = 1

type sendActionCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	nonce    uint64
	gasLimit uint64
	gasPrice *big.Int
	chainID  uint32
	payload  []byte
	core     *iotextypes.ActionCore
}

//API returns api
func (c *sendActionCaller) API() iotexapi.APIServiceClient {
	return c.api
}

func (c *sendActionCaller) setNonce(n uint64) {
	c.nonce = n
}

func (c *sendActionCaller) setGasLimit(g uint64) {
	c.gasLimit = g
}

func (c *sendActionCaller) setGasPrice(g *big.Int) {
	c.gasPrice = g
}

func (c *sendActionCaller) setPayload(pl []byte) {
	if pl == nil {
		return
	}
	c.payload = make([]byte, len(pl))
	copy(c.payload, pl)
}

func (c *sendActionCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.chainID == 0 {
		return hash.ZeroHash256, errcodes.New("0 is not a valid chain ID (use 1 for mainnet, 2 for testnet)", errcodes.InvalidParam)
	}
	c.core.ChainID = c.chainID

	if c.nonce == 0 {
		res, err := c.api.GetAccount(ctx, &iotexapi.GetAccountRequest{Address: c.account.Address().String()}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.nonce = res.GetAccountMeta().GetPendingNonce()
	}
	c.core.Nonce = c.nonce

	if c.gasLimit == 0 {
		sealed, err := sign(c.account, c.core)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
		}
		request := &iotexapi.EstimateGasForActionRequest{Action: sealed}
		response, err := c.api.EstimateGasForAction(ctx, request, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasLimit = response.GetGas()
	}
	c.core.GasLimit = c.gasLimit

	if c.gasPrice == nil {
		response, err := c.api.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasPrice = big.NewInt(0).SetUint64(response.GetGasPrice())
	}
	c.core.GasPrice = c.gasPrice.String()

	sealed, err := sign(c.account, c.core)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
	}

	response, err := c.api.SendAction(ctx, &iotexapi.SendActionRequest{Action: sealed}, opts...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
	}
	h, err := hash.HexStringToHash256(response.GetActionHash())
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.BadResponse)
	}
	c.nonce = 0 // reset before next time use
	return h, nil
}

type getReceiptCaller struct {
	api        iotexapi.APIServiceClient
	actionHash hash.Hash256
}

func (c *getReceiptCaller) Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error) {
	h := hex.EncodeToString(c.actionHash[:])
	return c.api.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h}, opts...)
}

type getLogsCaller struct {
	api     iotexapi.APIServiceClient
	Request *iotexapi.GetLogsRequest
}

func (c *getLogsCaller) Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetLogsResponse, error) {
	return c.api.GetLogs(ctx, c.Request, opts...)
}

func addressTypeAssert(preVal interface{}) (common.Address, error) {
	switch v := preVal.(type) {
	case string:
		ioAddress, err := address.FromString(v)
		if err != nil {
			return common.Address{}, errcodes.New("fail to convert string to ioAddress", errcodes.InvalidParam)
		}
		return common.HexToAddress(hex.EncodeToString(ioAddress.Bytes())), nil
	case address.Address:
		return common.HexToAddress(hex.EncodeToString(v.Bytes())), nil
	case common.Address:
		return v, nil
	default:
		return common.Address{}, errcodes.New("fail to convert from interface to string/ioAddress/ethAddress", errcodes.InvalidParam)
	}
}
