// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

// ProtocolVersion is the iotex protocol version to use. Currently 1.
const ProtocolVersion = 1

type sendActionCaller struct {
	account  account.Account
	api      iotexapi.APIServiceClient
	nonce    uint64
	gasLimit uint64
	gasPrice *big.Int
	chainID  uint32
	payload  []byte
	core     *iotextypes.ActionCore

	// Eth-typed transaction fields. When txType != 0, the action is built as
	// an Ethereum typed tx and submitted via Encoding_TX_CONTAINER instead of
	// the proto-native path.
	txType          uint32
	gasTipCap       *big.Int
	gasFeeCap       *big.Int
	accessList      types.AccessList
	blobTxData      *iotextypes.BlobTxData
	setCodeAuthList []SetCodeAuthorization
}

// API returns api
func (c *sendActionCaller) API() iotexapi.APIServiceClient {
	return c.api
}

func (c *sendActionCaller) setNonce(n uint64) {
	c.nonce = n
}

func (c *sendActionCaller) setGasLimit(g uint64) {
	c.gasLimit = g
}

func (c *sendActionCaller) setGasPrice(g *big.Int) {
	c.gasPrice = g
}

func (c *sendActionCaller) setPayload(pl []byte) {
	if pl == nil {
		return
	}
	c.payload = make([]byte, len(pl))
	copy(c.payload, pl)
}

func (c *sendActionCaller) setTxType(t uint32)      { c.txType = t }
func (c *sendActionCaller) setGasTipCap(v *big.Int) { c.gasTipCap = v }
func (c *sendActionCaller) setGasFeeCap(v *big.Int) { c.gasFeeCap = v }

// setAccessList stores a defensive deep copy so that a caller mutating the
// passed-in slice between the setter and Call() cannot corrupt SDK state.
func (c *sendActionCaller) setAccessList(l types.AccessList) {
	if l == nil {
		c.accessList = nil
		return
	}
	cp := make(types.AccessList, len(l))
	for i, tup := range l {
		cp[i].Address = tup.Address
		if n := len(tup.StorageKeys); n > 0 {
			cp[i].StorageKeys = make([]common.Hash, n)
			copy(cp[i].StorageKeys, tup.StorageKeys)
		}
	}
	c.accessList = cp
}

func (c *sendActionCaller) setBlobTxData(b *iotextypes.BlobTxData) {
	c.blobTxData = b
}
func (c *sendActionCaller) setSetCodeAuthList(a []SetCodeAuthorization) {
	c.setCodeAuthList = a
}

func (c *sendActionCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	if c.chainID == 0 {
		return hash.ZeroHash256, errcodes.New("0 is not a valid chain ID (use 1 for mainnet, 2 for testnet)", errcodes.InvalidParam)
	}
	c.core.ChainID = c.chainID

	if c.nonce == 0 {
		res, err := c.api.GetAccount(ctx, &iotexapi.GetAccountRequest{Address: c.account.Address().String()}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.nonce = res.GetAccountMeta().GetPendingNonce()
	}
	c.core.Nonce = c.nonce

	if c.txType != 0 {
		if c.gasLimit == 0 {
			return hash.ZeroHash256, errcodes.New("gas limit must be set explicitly for typed transactions", errcodes.InvalidParam)
		}
		if err := c.applyTypedFields(); err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InvalidParam)
		}
		c.core.GasLimit = c.gasLimit
		sealed, err := signContainer(c.account, c.core)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
		}
		return c.send(ctx, sealed, opts...)
	}

	if c.gasLimit == 0 {
		sealed, err := sign(c.account, c.core)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
		}
		request := &iotexapi.EstimateGasForActionRequest{Action: sealed}
		response, err := c.api.EstimateGasForAction(ctx, request, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasLimit = response.GetGas()
	}
	c.core.GasLimit = c.gasLimit

	if c.gasPrice == nil {
		response, err := c.api.SuggestGasPrice(ctx, &iotexapi.SuggestGasPriceRequest{}, opts...)
		if err != nil {
			return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
		}
		c.gasPrice = big.NewInt(0).SetUint64(response.GetGasPrice())
	}
	c.core.GasPrice = c.gasPrice.String()

	sealed, err := sign(c.account, c.core)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.InternalError)
	}
	return c.send(ctx, sealed, opts...)
}

func (c *sendActionCaller) send(ctx context.Context, sealed *iotextypes.Action, opts ...grpc.CallOption) (hash.Hash256, error) {
	response, err := c.api.SendAction(ctx, &iotexapi.SendActionRequest{Action: sealed}, opts...)
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.RPCError)
	}
	h, err := hash.HexStringToHash256(response.GetActionHash())
	if err != nil {
		return hash.ZeroHash256, errcodes.NewError(err, errcodes.BadResponse)
	}
	// reset per-call state so a reused caller does not carry stale fields
	c.nonce = 0
	c.txType = 0
	c.gasTipCap = nil
	c.gasFeeCap = nil
	c.accessList = nil
	c.blobTxData = nil
	c.setCodeAuthList = nil
	return h, nil
}

// applyTypedFields copies the user-supplied typed-tx fields onto c.core and
// validates that required fields for the chosen TxType are set.
func (c *sendActionCaller) applyTypedFields() error {
	c.core.TxType = c.txType
	if len(c.accessList) > 0 {
		c.core.AccessList = accessListToProto(c.accessList)
	}
	switch c.txType {
	case accessListTxType:
		if c.gasPrice == nil {
			return errcodes.New("gas price must be set for access list tx", errcodes.InvalidParam)
		}
		c.core.GasPrice = c.gasPrice.String()
	case dynamicFeeTxType, blobTxType, setCodeTxType:
		if c.gasTipCap == nil || c.gasFeeCap == nil {
			return errcodes.New("gas tip cap and fee cap must be set for typed transactions", errcodes.InvalidParam)
		}
		c.core.GasTipCap = c.gasTipCap.String()
		c.core.GasFeeCap = c.gasFeeCap.String()
		if c.txType == blobTxType {
			if c.blobTxData == nil {
				return errcodes.New("blob tx requires blobTxData", errcodes.InvalidParam)
			}
			c.core.BlobTxData = c.blobTxData
		}
		if c.txType == setCodeTxType {
			if len(c.setCodeAuthList) == 0 {
				return errcodes.New("setcode tx requires non-empty auth list", errcodes.InvalidParam)
			}
			authList, err := authListToProto(c.setCodeAuthList)
			if err != nil {
				return errcodes.NewError(err, errcodes.InvalidParam)
			}
			c.core.SetCodeAuthList = authList
		}
	default:
		return errcodes.New("unsupported tx type", errcodes.InvalidParam)
	}
	return nil
}

// accessListToProto mirrors iotex-core's encoding: hex strings without 0x.
func accessListToProto(list types.AccessList) []*iotextypes.AccessTuple {
	if len(list) == 0 {
		return nil
	}
	out := make([]*iotextypes.AccessTuple, len(list))
	for i, v := range list {
		out[i] = &iotextypes.AccessTuple{
			Address: hex.EncodeToString(v.Address.Bytes()),
		}
		if n := len(v.StorageKeys); n > 0 {
			out[i].StorageKeys = make([]string, n)
			for j, k := range v.StorageKeys {
				out[i].StorageKeys[j] = hex.EncodeToString(k.Bytes())
			}
		}
	}
	return out
}

type getReceiptCaller struct {
	api        iotexapi.APIServiceClient
	actionHash hash.Hash256
}

func (c *getReceiptCaller) Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error) {
	h := hex.EncodeToString(c.actionHash[:])
	return c.api.GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{ActionHash: h}, opts...)
}

type getLogsCaller struct {
	api     iotexapi.APIServiceClient
	Request *iotexapi.GetLogsRequest
}

func (c *getLogsCaller) Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetLogsResponse, error) {
	return c.api.GetLogs(ctx, c.Request, opts...)
}

func addressTypeAssert(preVal interface{}) (common.Address, error) {
	switch v := preVal.(type) {
	case string:
		ioAddress, err := address.FromString(v)
		if err != nil {
			return common.Address{}, errcodes.New("fail to convert string to ioAddress", errcodes.InvalidParam)
		}
		return common.HexToAddress(hex.EncodeToString(ioAddress.Bytes())), nil
	case address.Address:
		return common.HexToAddress(hex.EncodeToString(v.Bytes())), nil
	case common.Address:
		return v, nil
	default:
		return common.Address{}, errcodes.New("fail to convert from interface to string/ioAddress/ethAddress", errcodes.InvalidParam)
	}
}
