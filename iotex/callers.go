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
	gasLimit *uint64
	gasPrice *big.Int
	action   interface{}
	nonce    *uint64
}

func (c *sendActionCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.nonce == nil {
		res, err := c.api.GetAccount(ctx, &iotexapi.GetAccountRequest{Address: c.account.Address().String()}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		nonce := res.GetAccountMeta().GetPendingNonce()
		c.nonce = &nonce
	}
	core := &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Nonce:   *c.nonce,
	}

	switch a := c.action.(type) {
	case *iotextypes.Execution:
		core.Action = &iotextypes.ActionCore_Execution{Execution: a}
	case *iotextypes.Transfer:
		core.Action = &iotextypes.ActionCore_Transfer{Transfer: a}
	case *iotextypes.ClaimFromRewardingFund:
		core.Action = &iotextypes.ActionCore_ClaimFromRewardingFund{ClaimFromRewardingFund: a}
	default:
		return hash.ZeroHash256, errcodes.New("not support action call", errcodes.InternalError)
	}

	if c.gasLimit == nil {
		sealed, err := sign(c.account, core)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
		}
		request := &iotexapi.EstimateGasForActionRequest{Action: sealed}
		response, err := c.api.EstimateGasForAction(ctx, request, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		limit := response.GetGas()
		c.gasLimit = &limit
	}
	core.GasLimit = *c.gasLimit

	if c.gasPrice == nil {
		response, err := c.api.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasPrice = big.NewInt(0).SetUint64(response.GetGasPrice())
	}
	core.GasPrice = c.gasPrice.String()

	sealed, err := sign(c.account, core)
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
