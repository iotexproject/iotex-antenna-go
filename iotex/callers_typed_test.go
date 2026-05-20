// Copyright (c) 2026 IoTeX
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotexapi/mock_iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

func newExecCaller() *executeContractCaller {
	return &executeContractCaller{sendActionCaller: &sendActionCaller{}}
}

// TestExecuteCaller_TypedSettersPropagate confirms the chained setters on
// ExecuteContractCaller write through to the embedded sendActionCaller.
func TestExecuteCaller_TypedSettersPropagate(t *testing.T) {
	c := newExecCaller()
	blob := &BlobData{
		BlobFeeCap: big.NewInt(7),
		BlobHashes: []common.Hash{common.HexToHash("0x01")},
	}
	auths := []types.SetCodeAuthorization{{
		ChainID: *uint256.NewInt(2),
		Address: common.HexToAddress("0xab"),
		Nonce:   3,
	}}
	ret := c.SetTxType(blobTxType).
		SetGasTipCap(big.NewInt(1)).
		SetGasFeeCap(big.NewInt(2)).
		SetAccessList(types.AccessList{{Address: common.HexToAddress("0xcd")}}).
		SetBlobTxData(blob).
		SetSetCodeAuthList(auths)
	require.NotNil(t, ret)

	require.Equal(t, uint32(blobTxType), c.txType)
	require.Equal(t, big.NewInt(1), c.gasTipCap)
	require.Equal(t, big.NewInt(2), c.gasFeeCap)
	require.Len(t, c.accessList, 1)
	require.NotNil(t, c.blobTxData)
	require.Equal(t, "7", c.blobTxData.GetBlobFeeCap())
	require.Len(t, c.setCodeAuthList, 1)
	require.Equal(t, uint64(3), c.setCodeAuthList[0].Nonce)
}

func execCore() *iotextypes.ActionCore {
	return &iotextypes.ActionCore{
		Action: &iotextypes.ActionCore_Execution{Execution: &iotextypes.Execution{Amount: "0"}},
	}
}

func TestApplyTypedFields_Blob(t *testing.T) {
	c := &sendActionCaller{
		core:      execCore(),
		txType:    blobTxType,
		gasTipCap: big.NewInt(1),
		gasFeeCap: big.NewInt(2),
		blobTxData: (&BlobData{
			BlobFeeCap: big.NewInt(500),
			BlobHashes: []common.Hash{common.HexToHash("0x01")},
		}).toProto(),
	}
	require.NoError(t, c.applyTypedFields())
	require.Equal(t, uint32(blobTxType), c.core.TxType)
	require.Equal(t, "1", c.core.GasTipCap)
	require.Equal(t, "2", c.core.GasFeeCap)
	require.NotNil(t, c.core.BlobTxData)
	require.Equal(t, "500", c.core.BlobTxData.GetBlobFeeCap())
}

func TestApplyTypedFields_BlobMissingData(t *testing.T) {
	c := &sendActionCaller{
		core:      execCore(),
		txType:    blobTxType,
		gasTipCap: big.NewInt(1),
		gasFeeCap: big.NewInt(2),
	}
	err := c.applyTypedFields()
	require.Error(t, err)
	require.Contains(t, err.Error(), "blob tx requires blobTxData")
}

func TestApplyTypedFields_SetCode(t *testing.T) {
	c := &sendActionCaller{
		core:      execCore(),
		txType:    setCodeTxType,
		gasTipCap: big.NewInt(1),
		gasFeeCap: big.NewInt(2),
		setCodeAuthList: []types.SetCodeAuthorization{{
			ChainID: *uint256.NewInt(4690),
			Address: common.HexToAddress("0xab"),
			Nonce:   1,
		}},
	}
	require.NoError(t, c.applyTypedFields())
	require.Equal(t, uint32(setCodeTxType), c.core.TxType)
	require.Len(t, c.core.SetCodeAuthList, 1)
	require.Equal(t, uint32(4690), c.core.SetCodeAuthList[0].GetChainID())
}

func TestApplyTypedFields_SetCodeEmptyList(t *testing.T) {
	c := &sendActionCaller{
		core:      execCore(),
		txType:    setCodeTxType,
		gasTipCap: big.NewInt(1),
		gasFeeCap: big.NewInt(2),
	}
	err := c.applyTypedFields()
	require.Error(t, err)
	require.Contains(t, err.Error(), "non-empty auth list")
}

func TestApplyTypedFields_SetCodeChainIDOverflow(t *testing.T) {
	c := &sendActionCaller{
		core:      execCore(),
		txType:    setCodeTxType,
		gasTipCap: big.NewInt(1),
		gasFeeCap: big.NewInt(2),
		setCodeAuthList: []types.SetCodeAuthorization{{
			ChainID: *new(uint256.Int).Lsh(uint256.NewInt(1), 33),
			Address: common.HexToAddress("0xab"),
		}},
	}
	err := c.applyTypedFields()
	require.Error(t, err)
	require.Contains(t, err.Error(), "overflows uint32")
}

// TestCall_TypedTx_RoutesToContainer drives a full typed-tx Call() against a
// mock API and confirms the routing produces a TX_CONTAINER action carrying a
// decodable DynamicFee tx — covering the path applyTypedFields -> signContainer
// -> send that the network-gated integration tests otherwise skip.
func TestCall_TypedTx_RoutesToContainer(t *testing.T) {
	ctrl := gomock.NewController(t)

	acc, err := account.HexStringToAccount(_rlpTestKey)
	require.NoError(t, err)
	to, err := address.FromString(_rlpTestTo)
	require.NoError(t, err)

	api := mock_iotexapi.NewMockAPIServiceClient(ctrl)
	var captured *iotextypes.Action
	api.EXPECT().
		SendAction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, req *iotexapi.SendActionRequest, _ ...grpc.CallOption) (*iotexapi.SendActionResponse, error) {
			captured = req.GetAction()
			return &iotexapi.SendActionResponse{ActionHash: strings.Repeat("0", 64)}, nil
		})

	c := NewAuthedClient(api, _rlpTestChainID, acc)
	_, err = c.Transfer(to, big.NewInt(1)).
		SetNonce(9).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(big.NewInt(1000000000)).
		SetGasFeeCap(big.NewInt(2000000000)).
		SetGasLimit(50000).
		Call(context.Background())
	require.NoError(t, err)

	require.NotNil(t, captured)
	require.Equal(t, iotextypes.Encoding_TX_CONTAINER, captured.Encoding)
	require.Equal(t, acc.PublicKey().Bytes(), captured.SenderPubKey)
	require.Len(t, captured.Signature, 65)

	decoded := new(types.Transaction)
	require.NoError(t, decoded.UnmarshalBinary(captured.GetCore().GetTxContainer().GetRaw()))
	require.Equal(t, uint8(types.DynamicFeeTxType), decoded.Type())
	require.Equal(t, uint64(9), decoded.Nonce())
	require.Equal(t, uint64(50000), decoded.Gas())
}

// TestSetAccessList_DefensiveCopy confirms mutating the caller's slice after
// the setter does not change what the SDK will send.
func TestSetAccessList_DefensiveCopy(t *testing.T) {
	c := &sendActionCaller{}
	list := types.AccessList{{
		Address:     common.HexToAddress("0x01"),
		StorageKeys: []common.Hash{common.HexToHash("0xaa")},
	}}
	c.setAccessList(list)

	// mutate the original after the setter
	list[0].Address = common.HexToAddress("0x02")
	list[0].StorageKeys[0] = common.HexToHash("0xbb")

	require.Equal(t, common.HexToAddress("0x01"), c.accessList[0].Address)
	require.Equal(t, common.HexToHash("0xaa"), c.accessList[0].StorageKeys[0])
}
