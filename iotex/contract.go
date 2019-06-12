package iotex

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"google.golang.org/grpc"
)

type Data struct {
	method string
	abi    abi.ABI
	Raw    []byte
}

func (d Data) Unmarshal(v interface{}) error { return d.abi.Unpack(v, d.method, d.Raw) }

type ReadContractCaller interface {
	SetGasPrice(*big.Int) ReadContractCaller
	SetGasLimit(uint64) ReadContractCaller
	Call(ctx context.Context, opts ...grpc.CallOption) (Data, error)
}

type ExecuteContractCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) ExecuteContractCaller
	SetGasLimit(uint64) ExecuteContractCaller
}

type Contract interface {
	Read(method string, args ...interface{}) ReadContractCaller
	Execute(method string, args ...interface{}) ExecuteContractCaller
}

type contract struct {
	address string
	abi     abi.ABI
}
