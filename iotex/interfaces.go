package iotex

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

type SendActionCaller interface {
	API() iotexapi.APIServiceClient
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

	SetArgs(abi abi.ABI, args ...interface{}) DeployContractCaller
	SetGasPrice(*big.Int) DeployContractCaller
	SetGasLimit(uint64) DeployContractCaller
}

type GetReceiptCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error)
}

type AuthedClient interface {
	ReadOnlyClient

	Contract(contract address.Address, abi abi.ABI) Contract
	Transfer(to address.Address, value *big.Int) TransferCaller
	DeployContract(data []byte) DeployContractCaller
}

type ReadOnlyClient interface {
	ReadOnlyContract(contract address.Address, abi abi.ABI) ReadOnlyContract
	GetReceipt(actionHash hash.Hash256) GetReceiptCaller
}

type ReadContractCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (Data, error)
}

type ExecuteContractCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) ExecuteContractCaller
	SetGasLimit(uint64) ExecuteContractCaller
	SetAmount(*big.Int) ExecuteContractCaller
}

type Contract interface {
	ReadOnlyContract

	Execute(method string, args ...interface{}) ExecuteContractCaller
}

type ReadOnlyContract interface {
	Read(method string, args ...interface{}) ReadContractCaller
}
