// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

// This example shows how to programmatically deploy a contract to IoTeX blockchain and interact with it
// To run:
// go build; ./chaininfo

package main

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/iotexproject/iotex-address/address"

	"github.com/iotexproject/iotex-antenna-go/v2/examples/service"
)

// Xrc20Example is the Xrc20Example interface
type Xrc20Example interface {
	// Transfer is the Transfer interface
	Transfer(ctx context.Context, to string, amount *big.Int) (string, error)
}

type iotexService struct {
	service.IotexService

	contract address.Address
	abi      abi.ABI
	gasPrice *big.Int
	gasLimit uint64
}

// NewIotexService returns xrc20Example
func NewIotexService(accountPrivate, abiString, contract string, gasPrice *big.Int,
	gasLimit uint64, endpoint string, secure bool) (Xrc20Example, error) {
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}

	addr, err := address.FromString(contract)
	if err != nil {
		return nil, err
	}
	return &iotexService{
		service.NewIotexService(accountPrivate, endpoint, secure),
		addr, abi, gasPrice, gasLimit,
	}, nil
}

// Transfer is the Transfer interface
func (s *iotexService) Transfer(ctx context.Context, to string, amount *big.Int) (hash string, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	addr, err := address.FromString(to)
	if err != nil {
		return
	}
	ethAddr := common.HexToAddress(hex.EncodeToString(addr.Bytes()))
	h, err := s.AuthClient().Contract(s.contract, s.abi).Execute("transfer", ethAddr, amount).SetGasPrice(s.gasPrice).SetGasLimit(s.gasLimit).Call(ctx)
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}
