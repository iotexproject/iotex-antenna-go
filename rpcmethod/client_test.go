// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/action"
)

var (
	Address            = "io15jcpv957y5rn3zkyvd22cerfxcw4wc86hghyhn"
	PrivateKey         = "0806c458b262edd333a191e92f561aff338211ee3e18ab315a074a2d82aa343f"
	mainnetAddress     = "io1066kus4vlyvk0ljql39fzwqw0k22h7j8wmef3n"
	mainnetBlockHash   = "89bbf8b1d3cbfb6020a1074a11c5430ef77eb220c00143dbd6f76d1cab94a1c2"
	mainnetReceiptHash = "246b9b47f390a6faee9d725d9637b00b7ec56fa7cdffe3d39aeaad277edbb8f4"
)

const (
	testnet = "api.testnet.iotex.one:80"
	mainnet = "api.iotex.one:443"
)

func TestServer_GetAccount(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(mainnet, true)
	require.NoError(err)

	account := "io1066kus4vlyvk0ljql39fzwqw0k22h7j8wmef3n"
	request := &iotexapi.GetAccountRequest{Address: account}
	res, err := svr.GetAccount(request)
	require.NoError(err)
	accountMeta := res.AccountMeta
	require.Equal(account, accountMeta.Address)
	require.True(7 <= accountMeta.Nonce)
	require.True(9 == accountMeta.NumActions)
}

func TestServer_GetActions(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
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

func TestServer_SendAction(t *testing.T) {
	require := require.New(t)
	rpc, err := NewRPCMethod(testnet, false)
	require.NoError(err)
	accountPrivateKey := os.Getenv("accountPrivateKey")
	accountPendingNonce := os.Getenv("accountPendingNonce")
	if accountPrivateKey == "" || accountPendingNonce == "" {
		t.Skip("skipping test; some params not set")
	}
	accountPendingNonceInt, err := strconv.ParseUint(accountPendingNonce, 10, 64)
	require.NoError(err)

	act, err := account.NewAccountFromPrivateKey(accountPrivateKey)
	require.NoError(err)
	transfer, err := action.NewTransfer(
		accountPendingNonceInt,
		2000000,
		big.NewInt(1000000000000),
		big.NewInt(1000000000000000000),
		Address,
		nil)
	require.NoError(err)
	sealed, err := transfer.Sign(act)
	require.NoError(err)
	request := &iotexapi.SendActionRequest{Action: sealed.Action}
	_, err = rpc.SendAction(request)
	require.NoError(err)
}

func TestServer_GetAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	actionHash := "93de5923763c4ea79a01be023b49000838b1a4c22bdceed99dc23eeea8c9c757"
	actionActionInfoLen := "1"
	actionActionNonce := "27"

	actionActionInfoLenInt, err := strconv.ParseInt(actionActionInfoLen, 10, 64)
	require.NoError(err)
	actionActionNonceInt, err := strconv.ParseUint(actionActionNonce, 10, 64)
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
	require.Equal(int(actionActionInfoLenInt), len(res.ActionInfo))
	act := res.ActionInfo[0]
	require.Equal(actionActionNonceInt, act.Action.GetCore().GetNonce())
}

func TestServer_GetActionsByAddress(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	getActionsByAddressActionHash := "633cf62ab47611476423d7416bb74395be9c9b602062074ac36744ddd31fd122"
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByAddr{
			ByAddr: &iotexapi.GetActionsByAddressRequest{
				Address: mainnetAddress,
				Start:   1,
				Count:   1,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(getActionsByAddressActionHash, res.ActionInfo[0].ActHash)
	require.Equal(1, len(res.ActionInfo))
}

func TestServer_GetUnconfirmedActionsByAddress(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_UnconfirmedByAddr{
			UnconfirmedByAddr: &iotexapi.GetUnconfirmedActionsByAddressRequest{
				Address: mainnetAddress,
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
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByBlk{
			ByBlk: &iotexapi.GetActionsByBlockRequest{
				BlkHash: mainnetBlockHash,
				Start:   1,
				Count:   10,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(1, len(res.ActionInfo))
}

func TestServer_GetBlockMetas(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
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
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	request := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByHash{
			ByHash: &iotexapi.GetBlockMetaByHashRequest{
				BlkHash: mainnetBlockHash,
			},
		},
	}
	res, err := svr.GetBlockMetas(request)
	require.NoError(err)
	require.Equal(1, len(res.BlkMetas))
	blkPb := res.BlkMetas[0]
	require.Equal(int64(2), blkPb.NumActions)
	require.Equal("0", blkPb.TransferAmount)
}

func TestServer_GetChainMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	res, err := svr.GetChainMeta(&iotexapi.GetChainMetaRequest{})
	require.NoError(err)
	chainMetaPb := res.ChainMeta
	require.Equal(true, chainMetaPb.Height > 208646)
	require.Equal(true, chainMetaPb.NumActions > 211443)
	require.Equal(true, chainMetaPb.Tps >= 0)
	require.Equal(true, chainMetaPb.Epoch.Num >= 580)
	require.Equal(true, chainMetaPb.Epoch.Height >= 208441)
	require.Equal(true, chainMetaPb.Epoch.GravityChainStartHeight >= 7769900)
}

func TestServer_GetServerMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	res, err := svr.GetServerMeta(&iotexapi.GetServerMetaRequest{})
	require.NoError(err)
	require.Equal("0810e5166d12c7ae06110cf6429f332c59585056", res.GetServerMeta().PackageCommitID)
}

func TestServer_ReadState(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	out, err := svr.ReadState(&iotexapi.ReadStateRequest{
		ProtocolID: []byte("rewarding"),
		MethodName: []byte("UnclaimedBalance"),
		Arguments:  [][]byte{[]byte(mainnetAddress)},
	})
	require.NoError(err)
	require.NotNil(out)
	val, ok := big.NewInt(0).SetString(string(out.Data), 10)
	require.True(ok)
	expected, ok := new(big.Int).SetString("39860707937452088904761", 10)
	require.True(ok)
	require.Equal(1, val.Cmp(expected))
}

func TestServer_GetReceiptByAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	request := &iotexapi.GetReceiptByActionRequest{ActionHash: mainnetReceiptHash}
	res, err := svr.GetReceiptByAction(request)
	require.NoError(err)
	require.NotNil(res)
	receiptPb := res.ReceiptInfo.Receipt
	require.Equal(uint64(1), receiptPb.Status)
	require.Equal(uint64(56664), receiptPb.BlkHeight)
	require.NotEqual(hash.ZeroHash256, res.ReceiptInfo.BlkHash)
}

func TestServer_ReadContract(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(testnet, false)
	require.NoError(err)
	readContractActionHash := os.Getenv("readContractActionHash")
	if readContractActionHash == "" {
		t.Skip("skipping test; some params not set")
	}

	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash:   readContractActionHash,
				CheckPending: true,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	action := res.ActionInfo[0].Action
	require.NotNil(action)
	caller, _ := address.FromBytes(action.SenderPubKey)
	_, err = svr.ReadContract(&iotexapi.ReadContractRequest{
		Execution:     action.GetCore().GetExecution(),
		CallerAddress: caller.String(),
	})
	require.Error(err)
}

func TestServer_SuggestGasPrice(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)
	res, err := svr.SuggestGasPrice(&iotexapi.SuggestGasPriceRequest{})
	require.NoError(err)
	require.Equal(uint64(1000000000000), res.GasPrice)
}

func TestServer_EstimateGasForAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	act, err := account.NewAccountFromPrivateKey(PrivateKey)
	require.NoError(err)
	transfer, err := action.NewTransfer(
		3,
		20000,
		big.NewInt(10),
		big.NewInt(0),
		Address,
		nil)
	require.NoError(err)
	sealed, err := transfer.Sign(act)
	require.NoError(err)
	request := &iotexapi.EstimateGasForActionRequest{Action: sealed.Action}
	res, err := svr.EstimateGasForAction(request)
	require.NoError(err)
	require.Equal(uint64(10000), res.Gas)
}

func TestServer_GetEpochMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCWithTLSEnabled(mainnet)
	require.NoError(err)

	res, err := svr.GetEpochMeta(&iotexapi.GetEpochMetaRequest{EpochNumber: 1})
	require.NoError(err)
	require.Equal(uint64(1), res.EpochData.Num)
	require.Equal(uint64(1), res.EpochData.Height)
	require.Equal(uint64(0x743088), res.EpochData.GravityChainStartHeight)
	require.Equal(360, int(res.TotalBlocks))
	require.Equal(36, len(res.BlockProducersInfo))
}
