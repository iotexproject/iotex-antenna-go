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

	"github.com/iotexproject/iotex-antenna-go/v2/examples/service"
)

type contractExample interface {
	// Deploy is the Deploy interface
	Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (string, error)
	// BalanceOf is the BalanceOf interface
	BalanceOf(ctx context.Context, addre string) (balance *big.Int, err error)
}

type iotexService struct {
	service.IotexService

	contract address.Address
	abi      abi.ABI
	bin      string
	gasPrice *big.Int
	gasLimit uint64
}

// NewIotexService returns contractExample service
func NewIotexService(accountPrivate, abiString, binString, contract string, gasPrice *big.Int, gasLimit uint64, endpoint string, secure bool) (contractExample, error) {
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

	return &iotexService{
		service.NewIotexService(accountPrivate, endpoint, secure),
		addr, abi, binString, gasPrice, gasLimit,
	}, nil
}

// Deploy is the Deploy interface
func (s *iotexService) Deploy(ctx context.Context, waitContractAddress bool, args ...interface{}) (hash string, err error) {
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

// BalanceOf is the BalanceOf interface
func (s *iotexService) BalanceOf(ctx context.Context, addre string) (balance *big.Int, err error) {
	err = s.Connect()
	if err != nil {
		return
	}
	addr, err := address.FromString(addre)
	if err != nil {
		return
	}
	ethAddr := common.HexToAddress(hex.EncodeToString(addr.Bytes()))
	ret, err := s.ReadOnlyClient().ReadOnlyContract(s.contract, s.abi).Read("balanceOf", ethAddr).Call(ctx)
	if err != nil {
		return
	}
	balance = big.NewInt(0)
	err = ret.Unmarshal(&balance)
	return
}
