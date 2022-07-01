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
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type deployContractCaller struct {
	*sendActionCaller
	abi  *abi.ABI
	args []interface{}
}

func (c *deployContractCaller) SetArgs(abi abi.ABI, args ...interface{}) DeployContractCaller {
	c.abi = &abi
	c.args = args
	return c
}

func (c *deployContractCaller) SetGasLimit(g uint64) DeployContractCaller {
	c.sendActionCaller.setGasLimit(g)
	return c
}

func (c *deployContractCaller) SetGasPrice(g *big.Int) DeployContractCaller {
	c.sendActionCaller.setGasPrice(g)
	return c
}

func (c *deployContractCaller) SetNonce(n uint64) DeployContractCaller {
	c.sendActionCaller.setNonce(n)
	return c
}

func (c *deployContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if len(c.payload) == 0 {
		return hash.ZeroHash256, errcodes.New("contract data can not empty", errcodes.InvalidParam)
	}
	if len(c.args) > 0 {
		var err error
		c.args, err = encodeArgument(c.abi.Constructor, c.args)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
		}
		packed, err := c.abi.Pack("", c.args...)
		if err != nil {
			return hash.ZeroHash256, errcodes.New("failed to pack args", errcodes.InvalidParam)
		}
		c.payload = append(c.payload, packed...)
	}

	exec := &iotextypes.Execution{
		Data:   c.payload,
		Amount: "0",
	}
	c.core = &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Action:  &iotextypes.ActionCore_Execution{Execution: exec},
	}
	return c.sendActionCaller.Call(ctx, opts...)
}

type contractArgs struct {
	contract address.Address
	abi      *abi.ABI
	method   string
	args     []interface{}
}

type executeContractCaller struct {
	*sendActionCaller
	contractArgs
	amount *big.Int
}

func (c *executeContractCaller) SetAmount(a *big.Int) ExecuteContractCaller {
	c.amount = a
	return c
}

func (c *executeContractCaller) SetGasLimit(g uint64) ExecuteContractCaller {
	c.sendActionCaller.setGasLimit(g)
	return c
}

func (c *executeContractCaller) SetGasPrice(g *big.Int) ExecuteContractCaller {
	c.sendActionCaller.setGasPrice(g)
	return c
}

func (c *executeContractCaller) SetNonce(n uint64) ExecuteContractCaller {
	c.sendActionCaller.setNonce(n)
	return c
}

func (c *executeContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.method == "" {
		return hash.ZeroHash256, errcodes.New("contract address and method can not empty", errcodes.InvalidParam)
	}

	method, exist := c.abi.Methods[c.method]
	if !exist {
		return hash.ZeroHash256, errcodes.New("method is not found", errcodes.InvalidParam)
	}
	var err error
	c.args, err = encodeArgument(method, c.args)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
	}

	c.payload, err = c.abi.Pack(c.method, c.args...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
	}

	exec := &iotextypes.Execution{
		Contract: c.contract.String(),
		Data:     c.payload,
		Amount:   "0",
	}
	if c.amount != nil {
		exec.Amount = c.amount.String()
	}
	c.core = &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Action:  &iotextypes.ActionCore_Execution{Execution: exec},
	}
	return c.sendActionCaller.Call(ctx, opts...)
}

type readContractCaller struct {
	api iotexapi.APIServiceClient
	contractArgs
}

func (c *readContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (Data, error) {
	if c.method == "" {
		return Data{}, errcodes.New("contract address and method can not empty", errcodes.InvalidParam)
	}

	method, exist := c.abi.Methods[c.method]
	if !exist {
		return Data{}, errcodes.New("method is not found", errcodes.InvalidParam)
	}
	var err error
	c.args, err = encodeArgument(method, c.args)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.InvalidParam)
	}

	actData, err := c.abi.Pack(c.method, c.args...)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.InvalidParam)
	}

	request := &iotexapi.ReadContractRequest{
		Execution: &iotextypes.Execution{
			Contract: c.contract.String(),
			Data:     actData,
		},
		CallerAddress: address.ZeroAddress,
	}
	response, err := c.api.ReadContract(ctx, request, opts...)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.RPCError)
	}

	decoded, err := hex.DecodeString(response.GetData())
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.BadResponse)
	}

	return Data{
		method: c.method,
		abi:    c.abi,
		Raw:    decoded,
	}, nil
}

func encodeArgument(method abi.Method, args []interface{}) ([]interface{}, error) {
	if len(method.Inputs) != len(args) {
		return nil, errcodes.New("the number of arguments is not correct", errcodes.InvalidParam)
	}
	newArgs := make([]interface{}, len(args))
	for index, input := range method.Inputs {
		switch input.Type.String() {
		case "address":
			var err error
			newArgs[index], err = addressTypeAssert(args[index])
			if err != nil {
				return nil, errcodes.NewError(err, errcodes.InvalidParam)
			}
		case "address[]":
			s := reflect.ValueOf(args[index])
			if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
				return nil, errcodes.New("fail because the type is non-slice, non-array", errcodes.InvalidParam)
			}
			newArr := make([]common.Address, s.Len())
			for j := 0; j < s.Len(); j++ {
				var err error
				newArr[j], err = addressTypeAssert(s.Index(j).Interface())
				if err != nil {
					return nil, errcodes.NewError(err, errcodes.InvalidParam)
				}
			}
			newArgs[index] = newArr
		default:
			newArgs[index] = args[index]
		}
	}
	return newArgs, nil
}
