// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	ta "github.com/iotexproject/iotex-core/test/testaddress"
	"github.com/iotexproject/iotex-core/testutil"
)

const (
	//host = "127.0.0.1:14014"
	host        = "api.testnet.iotex.one:80"
	account     = "io13n3382cjhaawmqfk4vmvvgllnryw4tf56qdtks"
	actionHash  = "74bf4f3e943c2285899426302669c5bc2d479f0f3799160b41ac435bfa04fa47"
	blk1000Hash = "cce0233204fba5f1f259a3aeebd1b2aa12773039ff1f11dbc142025da624c3c9"
)

var (
	testTransfer, _ = testutil.SignedTransfer(ta.Addrinfo["alfa"].String(),
		ta.Keyinfo["alfa"].PriKey, 3, big.NewInt(10), []byte{}, testutil.TestGasLimit,
		big.NewInt(testutil.TestGasPriceInt64))

	testTransferPb = testTransfer.Proto()
)

func TestServer_GetAccount(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetAccountRequest{Address: account}
	res, err := svr.GetAccount(request)
	require.NoError(err)
	accountMeta := res.AccountMeta
	require.Equal(account, accountMeta.Address)
	require.Equal("100000986477999999999724250", accountMeta.Balance)
	require.Equal(uint64(0x2af4), accountMeta.Nonce)
	require.Equal(uint64(0x2af5), accountMeta.PendingNonce)
	require.Equal(true, accountMeta.NumActions > 0x42ba)

	// failure
	_, err = svr.GetAccount(&iotexapi.GetAccountRequest{})
	require.Error(err)
}

func TestServer_GetActions(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByIndex{
			ByIndex: &iotexapi.GetActionsByIndexRequest{
				Start: 0,
				Count: 5,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(5, len(res.ActionInfo))
}

func TestServer_GetAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash:   actionHash,
				CheckPending: true,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(1, len(res.ActionInfo))
	act := res.ActionInfo[0]
	require.Equal(uint64(0), act.Action.GetCore().GetNonce())
	require.Equal("044b4e708c2f408c8a34eb2a0b8824f4642b67e04032174066438a367fa61a59f4b6cbc64453509105f550abe96206cb1a1fd0fdb055c6a2b167460a17e3d86245", hex.EncodeToString(act.Action.SenderPubKey))
}

func TestServer_GetActionsByAddress(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByAddr{
			ByAddr: &iotexapi.GetActionsByAddressRequest{
				Address: account,
				Start:   1,
				Count:   1,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal("4a074ff58459691037cc5f67889837f0d4b99ca721f217ed8ee936f789c05c98", res.ActionInfo[0].ActHash)
	require.Equal(1, len(res.ActionInfo))
}

func TestServer_GetUnconfirmedActionsByAddress(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_UnconfirmedByAddr{
			UnconfirmedByAddr: &iotexapi.GetUnconfirmedActionsByAddressRequest{
				Address: account,
				Start:   1,
				Count:   10,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(0, len(res.ActionInfo))
}

func TestServer_GetActionsByBlock(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByBlk{
			ByBlk: &iotexapi.GetActionsByBlockRequest{
				BlkHash: blk1000Hash,
				Start:   1,
				Count:   10,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(0, len(res.ActionInfo))
}

func TestServer_GetBlockMetas(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: 10,
				Count: 10,
			},
		},
	}
	res, err := svr.GetBlockMetas(request)
	require.NoError(err)
	require.Equal(10, len(res.GetBlkMetas()))
	var prevBlkPb *iotextypes.BlockMeta
	for _, blkPb := range res.BlkMetas {
		if prevBlkPb != nil {
			require.True(blkPb.Height > prevBlkPb.Height)
		}
		prevBlkPb = blkPb
	}

}

func TestServer_GetBlockMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByHash{
			ByHash: &iotexapi.GetBlockMetaByHashRequest{
				BlkHash: blk1000Hash,
			},
		},
	}
	res, err := svr.GetBlockMetas(request)
	require.NoError(err)
	require.Equal(1, len(res.BlkMetas))
	blkPb := res.BlkMetas[0]
	require.Equal(int64(1), blkPb.NumActions)
	require.Equal("0", blkPb.TransferAmount)
}

func TestServer_GetChainMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	res, err := svr.GetChainMeta(&iotexapi.GetChainMetaRequest{})
	require.NoError(err)
	chainMetaPb := res.ChainMeta
	require.Equal(true, chainMetaPb.Height > 1)
	require.Equal(true, chainMetaPb.NumActions > 1)
	require.Equal(true, chainMetaPb.Tps == 0)
	require.Equal(true, chainMetaPb.Epoch.Num > 1)
	require.Equal(true, chainMetaPb.Epoch.Height > 1)
	require.Equal(true, chainMetaPb.Epoch.GravityChainStartHeight > 1)
}

func TestServer_GetServerMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.GetServerMeta(&iotexapi.GetServerMetaRequest{})
	require.NoError(err)
	require.Equal("4977f444c32d830a55a47449bf2330202d7338cb", res.GetServerMeta().PackageCommitID)
}
func TestServer_ReadState(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	out, err := svr.ReadState(&iotexapi.ReadStateRequest{
		ProtocolID: []byte("rewarding"),
		MethodName: []byte("UnclaimedBalance"),
		Arguments:  [][]byte{[]byte(account)},
	})
	require.NoError(err)
	require.NotNil(out)
	val, ok := big.NewInt(0).SetString(string(out.Data), 10)
	require.True(ok)
	expected, ok := new(big.Int).SetString("3712000000000000000000", 10)
	require.True(ok)
	require.Equal(1, val.Cmp(expected))
}

func TestServer_SuggestGasPrice(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.SuggestGasPrice(&iotexapi.SuggestGasPriceRequest{})
	require.NoError(err)
	require.Equal(uint64(1), res.GasPrice)
}

func TestServer_EstimateGasForAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &iotexapi.EstimateGasForActionRequest{Action: testTransferPb}
	res, err := svr.EstimateGasForAction(request)
	require.NoError(err)
	require.Equal(uint64(10000), res.Gas)
}

func TestServer_GetEpochMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	res, err := svr.GetEpochMeta(&iotexapi.GetEpochMetaRequest{EpochNumber: 1})
	require.NoError(err)
	require.Equal(uint64(1), res.EpochData.Num)
	require.Equal(uint64(1), res.EpochData.Height)
	require.Equal(uint64(0x731874), res.EpochData.GravityChainStartHeight)
	require.Equal(360, int(res.TotalBlocks))
	require.Equal(24, len(res.BlockProducersInfo))
}
