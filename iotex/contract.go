package iotex

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"google.golang.org/grpc"
)

type Data struct {
	method string
	abi    *abi.ABI
	Raw    []byte
}

func (d Data) Unmarshal(v interface{}) error { return d.abi.Unpack(v, d.method, d.Raw) }

type ReadContractCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (Data, error)
}

type ExecuteContractCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) ExecuteContractCaller
	SetGasLimit(uint64) ExecuteContractCaller
}

type Contract interface {
	ReadOnlyContract

	Execute(method string, args ...interface{}) ExecuteContractCaller
}

type ReadOnlyContract interface {
	Read(method string, args ...interface{}) ReadContractCaller
}

type contract struct {
	address address.Address
	abi     *abi.ABI
}
