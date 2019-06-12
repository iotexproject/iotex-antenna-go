package iotex

import (
	"context"
	"math/big"

	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

type SendActionCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (actHash string, err error)
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

	Contract(contractAddr, abi string) (Contract, error)
	Transfer(to string, value big.Int) TransferCaller
	DeployContract(abi string, data []byte) DeployContractCaller
}

type ReadOnlyClient interface {
	GetReceipt(actionHash string) GetReceiptCaller
}

func NewAuthedClient(iotexapi.APIServiceClient, account.Account) AuthedClient { return nil }

func NewReadOnlyClient(iotexapi.APIServiceClient) ReadOnlyClient { return nil }
