package iotex

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

type SendActionCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error)
}

type TransferCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) TransferCaller
	SetGasLimit(uint64) TransferCaller
	SetPayload([]byte) TransferCaller
}

type DeployContractCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) DeployContractCaller
	SetGasLimit(uint64) DeployContractCaller
}

type GetReceiptCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error)
}

type AuthedClient interface {
	ReadOnlyClient

	Contract(contract address.Address, abi abi.ABI) (Contract, error)
	Transfer(to address.Address, value *big.Int) TransferCaller
	DeployContract(data []byte) DeployContractCaller
}

type ReadOnlyClient interface {
	ReadOnlyContract(contract address.Address, abi abi.ABI) (ReadOnlyContract, error)
	GetReceipt(actionHash hash.Hash256) GetReceiptCaller
}

func NewAuthedClient(iotexapi.APIServiceClient, account.Account) AuthedClient { return nil }

func NewReadOnlyClient(iotexapi.APIServiceClient) ReadOnlyClient { return nil }
