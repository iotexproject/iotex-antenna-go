// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-antenna-go/v2/examples/util"
)

// Xrc20Service is the Xrc20Service interface
type Xrc20Service interface {
	// Deploy is the Deploy interface
	Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (string, error)
	// Transfer is the Transfer interface
	Transfer(ctx context.Context, to string, amount *big.Int) (string, error)
	// BalanceOf is the BalanceOf interface
	BalanceOf(ctx context.Context, addr string) (*big.Int, error)
}

type xrc20Service struct {
	util.IotexService

	contract address.Address
	abi      abi.ABI
	bin      string
	gasPrice *big.Int
	gasLimit uint64
}

// NewXrc20Service returns Xrc20Service
func NewXrc20Service(accountPrivate, abiString, binString, contract string, gasPrice *big.Int, gasLimit uint64, endpoint string, secure bool) (Xrc20Service, error) {
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	var addr address.Address
	if contract != "" {
		addr, err = address.FromString(contract)
		if err != nil {
			return nil, err
		}
	}
	return &xrc20Service{
		util.NewIotexService(accountPrivate, endpoint, secure),
		addr, abi, binString, gasPrice, gasLimit,
	}, nil
}

// Deploy is the Deploy interface
func (s *xrc20Service) Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (hash string, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	data, err := hex.DecodeString(s.bin)
	if err != nil {
		return
	}
	h, err := s.AuthClient().DeployContract(data).SetGasPrice(s.gasPrice).SetGasLimit(s.gasLimit).SetArgs(s.abi, args...).Call(ctx)
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	if waitContractAddress {
		time.Sleep(time.Second * 10)
		receiptResponse, err := s.AuthClient().GetReceipt(h).Call(ctx)
		if err != nil {
			return "", err
		}
		status := receiptResponse.GetReceiptInfo().GetReceipt().GetStatus()
		if status != uint64(iotextypes.ReceiptStatus_Success) {
			return "", errors.New("deploy error,status:" + fmt.Sprintf("%d", status))
		}
		addr := receiptResponse.GetReceiptInfo().GetReceipt().GetContractAddress()
		s.contract, err = address.FromString(addr)
		if err != nil {
			return "", err
		}
	}
	return
}

// Transfer is the Transfer interface
func (s *xrc20Service) Transfer(ctx context.Context, to string, amount *big.Int) (hash string, err error) {
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

// BalanceOf is the BalanceOf interface
func (s *xrc20Service) BalanceOf(ctx context.Context, addr string) (balance *big.Int, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	ret, err := s.ReadOnlyClient().ReadOnlyContract(s.contract, s.abi).Read("balanceOf", addr).Call(ctx)
	if err != nil {
		return
	}
	balance = new(big.Int).SetBytes(ret.Raw)
	return
}
