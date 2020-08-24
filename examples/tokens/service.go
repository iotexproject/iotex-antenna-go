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

type xrc20Example interface {
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
	gasLimit uint64, endpoint string, secure bool) (xrc20Example, error) {
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
