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

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-antenna-go/v2/examples/util"
)

// OpenOracleService is the OpenOracleService interface
type OpenOracleService interface {
	// Deploy is the Deploy interface
	Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (string, error)
	// Put is the Put interface
	Put(ctx context.Context, message []byte, signature []byte) (string, error)
	// Get is the Get interface
	Get(ctx context.Context, source, key string) (string, error)
}

type openOracleService struct {
	util.IotexService

	contract address.Address
	abi      abi.ABI
	bin      string
	gasPrice *big.Int
	gasLimit uint64
}

// NewOpenOracleService returns OpenOracleService
func NewOpenOracleService(accountPrivate, abiString, binString, contract string, gasPrice *big.Int, gasLimit uint64, endpoint string, secure bool) (OpenOracleService, error) {
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

	return &openOracleService{
		util.NewIotexService(accountPrivate, endpoint, secure),
		addr, abi, binString, gasPrice, gasLimit,
	}, nil
}

// Deploy is the Deploy interface
func (s *openOracleService) Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (hash string, err error) {
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

// Put is the Put interface
func (s *openOracleService) Put(ctx context.Context, message, signature []byte) (hash string, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	h, err := s.AuthClient().Contract(s.contract, s.abi).Execute("put", message, signature).SetGasPrice(s.gasPrice).SetGasLimit(s.gasLimit).Call(ctx)
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}

// Get is the Get interface
func (s *openOracleService) Get(ctx context.Context, source, key string) (ret string, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	data, err := s.ReadOnlyClient().ReadOnlyContract(s.contract, s.abi).Read("get", source, key).Call(ctx)
	if err != nil {
		return
	}
	ret = hex.EncodeToString(data.Raw)
	return
}
