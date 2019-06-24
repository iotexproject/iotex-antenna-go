package iotex

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/errcodes"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"
)

// ProtocolVersion is the iotex protocol version to use. Currently 1.
const ProtocolVersion = 1

type sendActionCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	gasLimit *uint64
	gasPrice *big.Int
	action   interface{}
}

func (c *sendActionCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	res, err := c.api.GetAccount(ctx, &iotexapi.GetAccountRequest{Address: c.account.Address().String()}, opts...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
	}
	core := &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Nonce:   res.GetAccountMeta().GetPendingNonce(),
	}

	switch a := c.action.(type) {
	case *iotextypes.Execution:
		core.Action = &iotextypes.ActionCore_Execution{Execution: a}
	case *iotextypes.Transfer:
		core.Action = &iotextypes.ActionCore_Transfer{Transfer: a}
	default:
		return hash.ZeroHash256, errcodes.New("not support action core", errcodes.InternalError)
	}

	if c.gasLimit == nil {
		sealed, err := sign(c.account, core)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
		}
		request := &iotexapi.EstimateGasForActionRequest{Action: sealed}
		response, err := c.api.EstimateGasForAction(ctx, request, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		limit := response.GetGas()
		c.gasLimit = &limit
	}
	core.GasLimit = *c.gasLimit

	if c.gasPrice == nil {
		response, err := c.api.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasPrice = big.NewInt(0).SetUint64(response.GetGasPrice())
	}
	core.GasPrice = c.gasPrice.String()

	sealed, err := sign(c.account, core)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
	}

	response, err := c.api.SendAction(ctx, &iotexapi.SendActionRequest{Action: sealed}, opts...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
	}
	h, err := hash.HexStringToHash256(response.GetActionHash())
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.BadResponse)
	}
	return h, nil
}

type transferCaller struct {
	account   account.Account
	api       iotexapi.APIServiceClient
	amount    *big.Int
	recipient address.Address
	payload   []byte
	gasLimit  *uint64
	gasPrice  *big.Int
}

func (c *transferCaller) SetPayload(pl []byte) TransferCaller {
	if pl == nil {
		return c
	}
	c.payload = make([]byte, len(pl))
	copy(c.payload, pl)
	return c
}

func (c *transferCaller) SetGasLimit(g uint64) TransferCaller {
	c.gasLimit = &g
	return c
}

func (c *transferCaller) SetGasPrice(g *big.Int) TransferCaller {
	c.gasPrice = g
	return c
}

func (c *transferCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *transferCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("transfer amount cannot be nil", errcodes.InvalidParam)
	}

	tx := &iotextypes.Transfer{
		Amount:    c.amount.String(),
		Recipient: c.recipient.String(),
		Payload:   c.payload,
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		action:   tx,
	}
	return sc.Call(ctx, opts...)
}

type deployContractCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	gasLimit *uint64
	gasPrice *big.Int
	abi      *abi.ABI
	args     []interface{}
	data     []byte
}

func (c *deployContractCaller) SetArgs(abi abi.ABI, args ...interface{}) DeployContractCaller {
	c.abi = &abi
	c.args = args
	return c
}

func (c *deployContractCaller) SetGasLimit(g uint64) DeployContractCaller {
	c.gasLimit = &g
	return c
}

func (c *deployContractCaller) SetGasPrice(g *big.Int) DeployContractCaller {
	c.gasPrice = g
	return c
}

func (c *deployContractCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *deployContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if len(c.data) == 0 {
		return hash.ZeroHash256, errcodes.New("contract data can not empty", errcodes.InvalidParam)
	}

	if len(c.args) > 0 {
		packed, err := c.abi.Pack("", c.args...)
		if err != nil {
			return hash.ZeroHash256, errcodes.New("failed to pack args", errcodes.InvalidParam)
		}
		c.data = append(c.data, packed...)
	}

	exec := &iotextypes.Execution{
		Data: c.data,
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		action:   exec,
	}
	return sc.Call(ctx, opts...)
}

type executeContractCaller struct {
	abi      *abi.ABI
	contract address.Address
	account  account.Account
	api      iotexapi.APIServiceClient
	method   string
	args     []interface{}
	amount   *big.Int
	gasLimit *uint64
	gasPrice *big.Int
}

func (c *executeContractCaller) SetAmount(a *big.Int) ExecuteContractCaller {
	c.amount = a
	return c
}

func (c *executeContractCaller) SetGasLimit(g uint64) ExecuteContractCaller {
	c.gasLimit = &g
	return c
}

func (c *executeContractCaller) SetGasPrice(g *big.Int) ExecuteContractCaller {
	c.gasPrice = g
	return c
}

func (c *executeContractCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *executeContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.method == "" {
		return hash.ZeroHash256, errcodes.New("contract address and method can not empty", errcodes.InvalidParam)
	}

	actData, err := c.abi.Pack(c.method, c.args...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
	}

	exec := &iotextypes.Execution{
		Contract: c.contract.String(),
		Data:     actData,
	}
	if c.amount != nil {
		exec.Amount = c.amount.String()
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		action:   exec,
	}
	return sc.Call(ctx, opts...)
}

type readContractCaller struct {
	method string
	args   []interface{}
	sender address.Address
	rc     *readOnlyContract
}

func (c *readContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (Data, error) {
	if c.method == "" {
		return Data{}, errcodes.New("contract address and method can not empty", errcodes.InvalidParam)
	}

	actData, err := c.rc.abi.Pack(c.method, c.args...)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.InvalidParam)
	}

	request := &iotexapi.ReadContractRequest{
		Execution: &iotextypes.Execution{
			Contract: c.rc.address.String(),
			Data:     actData,
		},
		CallerAddress: c.sender.String(),
	}
	response, err := c.rc.api.ReadContract(ctx, request, opts...)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.RPCError)
	}

	decoded, err := hex.DecodeString(response.GetData())
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.BadResponse)
	}

	return Data{
		method: c.method,
		abi:    c.rc.abi,
		Raw:    decoded,
	}, nil
}

type getReceiptCaller struct {
	api        iotexapi.APIServiceClient
	actionHash hash.Hash256
}

func (c *getReceiptCaller) Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error) {
	h := hex.EncodeToString(c.actionHash[:])
	return c.api.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h}, opts...)
}
