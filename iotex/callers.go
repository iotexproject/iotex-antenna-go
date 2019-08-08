package iotex

import (
	"context"
	"encoding/hex"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
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
	nonce    *uint64
}

func (c *sendActionCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.nonce == nil {
		res, err := c.api.GetAccount(ctx, &iotexapi.GetAccountRequest{Address: c.account.Address().String()}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		nonce := res.GetAccountMeta().GetPendingNonce()
		c.nonce = &nonce
	}
	core := &iotextypes.ActionCore{
		Version: ProtocolVersion,
		Nonce:   *c.nonce,
	}

	switch a := c.action.(type) {
	case *iotextypes.Execution:
		core.Action = &iotextypes.ActionCore_Execution{Execution: a}
	case *iotextypes.Transfer:
		core.Action = &iotextypes.ActionCore_Transfer{Transfer: a}
	case *iotextypes.ClaimFromRewardingFund:
		core.Action = &iotextypes.ActionCore_ClaimFromRewardingFund{ClaimFromRewardingFund: a}
	default:
		return hash.ZeroHash256, errcodes.New("not support action call", errcodes.InternalError)
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
	nonce     *uint64
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

func (c *transferCaller) SetNonce(n uint64) TransferCaller {
	c.nonce = &n
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
		nonce:    c.nonce,
		action:   tx,
	}
	return sc.Call(ctx, opts...)
}

type claimRewardCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	amount   *big.Int
	data     []byte
	gasLimit *uint64
	gasPrice *big.Int
	nonce    *uint64
}

func (c *claimRewardCaller) SetData(data []byte) ClaimRewardCaller {
	if data == nil {
		return c
	}
	c.data = make([]byte, len(data))
	copy(c.data, data)
	return c
}

func (c *claimRewardCaller) SetGasLimit(g uint64) ClaimRewardCaller {
	c.gasLimit = &g
	return c
}

func (c *claimRewardCaller) SetGasPrice(g *big.Int) ClaimRewardCaller {
	c.gasPrice = g
	return c
}

func (c *claimRewardCaller) SetNonce(n uint64) ClaimRewardCaller {
	c.nonce = &n
	return c
}

func (c *claimRewardCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *claimRewardCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.amount == nil {
		return hash.ZeroHash256, errcodes.New("claim amount cannot be nil", errcodes.InvalidParam)
	}

	tx := &iotextypes.ClaimFromRewardingFund{
		Amount: c.amount.String(),
		Data:   c.data,
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		nonce:    c.nonce,
		action:   tx,
	}
	return sc.Call(ctx, opts...)
}

type deployContractCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	gasLimit *uint64
	gasPrice *big.Int
	nonce    *uint64
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

func (c *deployContractCaller) SetNonce(n uint64) DeployContractCaller {
	c.nonce = &n
	return c
}

func (c *deployContractCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *deployContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if len(c.data) == 0 {
		return hash.ZeroHash256, errcodes.New("contract data can not empty", errcodes.InvalidParam)
	}
	if len(c.args) > 0 {
		var err error
		c.args, err = encodeArgument(c.abi.Constructor, c.args)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
		}
		packed, err := c.abi.Pack("", c.args...)
		if err != nil {
			return hash.ZeroHash256, errcodes.New("failed to pack args", errcodes.InvalidParam)
		}
		c.data = append(c.data, packed...)
	}

	exec := &iotextypes.Execution{
		Data:   c.data,
		Amount: "0",
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		nonce:    c.nonce,
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
	nonce    *uint64
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

func (c *executeContractCaller) SetNonce(n uint64) ExecuteContractCaller {
	c.nonce = &n
	return c
}

func (c *executeContractCaller) API() iotexapi.APIServiceClient { return c.api }

func (c *executeContractCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.method == "" {
		return hash.ZeroHash256, errcodes.New("contract address and method can not empty", errcodes.InvalidParam)
	}

	method, exist := c.abi.Methods[c.method]
	if !exist {
		return hash.ZeroHash256, errcodes.New("method is not found", errcodes.InvalidParam)
	}
	var err error
	c.args, err = encodeArgument(method, c.args)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
	}

	actData, err := c.abi.Pack(c.method, c.args...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
	}

	exec := &iotextypes.Execution{
		Contract: c.contract.String(),
		Data:     actData,
		Amount:   "0",
	}
	if c.amount != nil {
		exec.Amount = c.amount.String()
	}
	sc := &sendActionCaller{
		account:  c.account,
		api:      c.api,
		gasLimit: c.gasLimit,
		gasPrice: c.gasPrice,
		nonce:    c.nonce,
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

	method, exist := c.rc.abi.Methods[c.method]
	if !exist {
		return Data{}, errcodes.New("method is not found", errcodes.InvalidParam)
	}
	var err error
	c.args, err = encodeArgument(method, c.args)
	if err != nil {
		return Data{}, errcodes.NewError(err, errcodes.InvalidParam)
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

func encodeArgument(method abi.Method, args []interface{}) ([]interface{}, error) {
	if len(method.Inputs) != len(args) {
		return nil, errcodes.New("the number of arguments is not correct", errcodes.InvalidParam)
	}
	newArgs := make([]interface{}, len(args))
	for index, input := range method.Inputs {
		if input.Type.String() == "address" {
			s, ok := args[index].(string)
			if !ok {
				return nil, errcodes.New("fail to convert from interface to string", errcodes.InvalidParam)
			}
			ethAddress, err := address.FromString(s)
			if err != nil {
				return nil, errcodes.New("fail to convert string to address(iotex)", errcodes.InvalidParam)
			}
			newArgs[index] = common.HexToAddress(hex.EncodeToString(ethAddress.Bytes()))

		} else if input.Type.String() == "address[]" {
			s := reflect.ValueOf(args[index])
			if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
				panic("InterfaceSlice() given a non-slice and non-array type")
			}
			ret := make([]interface{}, s.Len())
			for i := 0; i < s.Len(); i++ {
				ret[i] = s.Index(i).Interface()
			}
			newArr := make([]common.Address, s.Len())
			for j, elem := range ret {
				str, ok := elem.(string)
				if !ok {
					return nil, errcodes.New("fail to convert from interface to string in array", errcodes.InvalidParam)
				}
				ethAddress, err := address.FromString(str)
				if err != nil {
					return nil, errcodes.New("fail to convert from string to address(iotex)", errcodes.InvalidParam)
				}
				newArr[j] = common.HexToAddress(hex.EncodeToString(ethAddress.Bytes()))
			}
			newArgs[index] = newArr
		} else {
			newArgs[index] = args[index]
		}
	}
	return newArgs, nil
}
