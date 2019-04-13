// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	ta "github.com/iotexproject/iotex-core/test/testaddress"
	"github.com/iotexproject/iotex-core/testutil"
)

const (
	//export accountPrivateKey="9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	//export accountAddress="io14gnqxf9dpkn05g337rl7eyt2nxasphf5m6n0rd"
	//export accountBalance="99994712164399999990848350"
	//export accountNonce="337529"
	//export accountPendingNonce="337530"
	//export accountNumActions="506609"
	//export actionHash="b9a938e1f249d3c7ab8152e377132989535e25ea9ee376323179d1943dc15b4e"
	//export actionActionInfoLen="1"
	//export actionActionNonce="1"
	//export getActionsByAddressActionHash="10efdacf68be2fa0afdc2a46786b3caf0a59fe2386485a9f075acf2a41c93d78"
	//export blk60801Hash="a331841e6b29becdeeb65cb0948d896076b514e5f8d69d560b4424282a7882d7"
	//export blk60801HashNumActions="1960"
	//export blk60801HashTransferAmount="2389000000000000000000"
	//export getServerMetaPackageCommitID="e24aadf90f98d800b9f117354ddd5b3dbe58dde9"
	//export accountAddressUnclaimedBalance="0"
	//export getReceiptByActionBlkHeight="517"
	//
	//export epochDataHeight="1"
	//export epochGravityChainStartHeight="7502300"
	//export readContractActionHash="63c74277bcbfcef195f57713131b05cb54c47461f7e64b8f32fb58f9b8445265"
	host = "api.iotex.one:80"
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
	accountAddress := os.Getenv("accountAddress")

	accountBalance := os.Getenv("accountBalance")
	accountNonce := os.Getenv("accountNonce")
	accountNonceInt, err := strconv.ParseUint(accountNonce, 10, 64)
	require.NoError(err)
	accountPendingNonce := os.Getenv("accountPendingNonce")
	accountPendingNonceInt, err := strconv.ParseUint(accountPendingNonce, 10, 64)
	require.NoError(err)
	accountNumActions := os.Getenv("accountNumActions")
	accountNumActionsInt, err := strconv.ParseUint(accountNumActions, 10, 64)
	require.NoError(err)
	request := &iotexapi.GetAccountRequest{Address: accountAddress}
	res, err := svr.GetAccount(request)
	require.NoError(err)
	accountMeta := res.AccountMeta
	require.Equal(accountAddress, accountMeta.Address)
	require.NotEqual(accountBalance, accountMeta.Balance)
	require.Equal(true, accountNonceInt < accountMeta.Nonce)
	require.Equal(true, accountPendingNonceInt < accountMeta.PendingNonce)
	require.Equal(true, accountNumActionsInt < accountMeta.NumActions)

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

func TestServer_SendAction(t *testing.T) {
	t.Skip("Skip,make sure env is right")
	require := require.New(t)
	rpc, err := NewRPCMethod(host)
	require.NoError(err)
	accountPrivateKey := os.Getenv("accountPrivateKey")
	accountPendingNonce := os.Getenv("accountPendingNonce")
	accountPendingNonceInt, err := strconv.ParseUint(accountPendingNonce, 10, 64)
	priKey, err := keypair.HexStringToPrivateKey(accountPrivateKey)
	require.NoError(err)

	testTransfer, err := testutil.SignedTransfer("io15jcpv957y5rn3zkyvd22cerfxcw4wc86hghyhn",
		priKey, accountPendingNonceInt, big.NewInt(1000000000000000000), []byte{}, 2000000,
		big.NewInt(1000000000000))
	require.NoError(err)
	testTransferPb := testTransfer.Proto()
	request := &iotexapi.SendActionRequest{Action: testTransferPb}
	res, err := rpc.SendAction(request)
	require.NoError(err)
	fmt.Println("res:", res)
}

func TestServer_GetAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	actionHash := os.Getenv("actionHash")
	actionActionInfoLen := os.Getenv("actionActionInfoLen")
	actionActionInfoLenInt, err := strconv.ParseInt(actionActionInfoLen, 10, 64)
	require.NoError(err)
	actionActionNonce := os.Getenv("actionActionNonce")
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
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	accountAddress := os.Getenv("accountAddress")
	getActionsByAddressActionHash := os.Getenv("getActionsByAddressActionHash")
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByAddr{
			ByAddr: &iotexapi.GetActionsByAddressRequest{
				Address: accountAddress,
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
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	accountAddress := os.Getenv("accountAddress")
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_UnconfirmedByAddr{
			UnconfirmedByAddr: &iotexapi.GetUnconfirmedActionsByAddressRequest{
				Address: accountAddress,
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
	blk60801Hash := os.Getenv("blk60801Hash")
	request := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByBlk{
			ByBlk: &iotexapi.GetActionsByBlockRequest{
				BlkHash: blk60801Hash,
				Start:   1,
				Count:   10,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	require.Equal(10, len(res.ActionInfo))
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
	blk60801Hash := os.Getenv("blk60801Hash")
	blk60801HashNumActions := os.Getenv("blk60801HashNumActions")
	blk60801HashNumActionsInt, err := strconv.ParseInt(blk60801HashNumActions, 10, 64)
	blk60801HashTransferAmount := os.Getenv("blk60801HashTransferAmount")

	request := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByHash{
			ByHash: &iotexapi.GetBlockMetaByHashRequest{
				BlkHash: blk60801Hash,
			},
		},
	}
	res, err := svr.GetBlockMetas(request)
	require.NoError(err)
	require.Equal(1, len(res.BlkMetas))
	blkPb := res.BlkMetas[0]
	require.Equal(blk60801HashNumActionsInt, blkPb.NumActions)
	require.Equal(blk60801HashTransferAmount, blkPb.TransferAmount)
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
	require.Equal(true, chainMetaPb.Tps > 0)
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
	getServerMetaPackageCommitID := os.Getenv("getServerMetaPackageCommitID")
	require.Equal(getServerMetaPackageCommitID, res.GetServerMeta().PackageCommitID)
}

func TestServer_ReadState(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	accountAddress := os.Getenv("accountAddress")
	accountAddressUnclaimedBalance := os.Getenv("accountAddressUnclaimedBalance")
	out, err := svr.ReadState(&iotexapi.ReadStateRequest{
		ProtocolID: []byte("rewarding"),
		MethodName: []byte("UnclaimedBalance"),
		Arguments:  [][]byte{[]byte(accountAddress)},
	})
	require.NoError(err)
	require.NotNil(out)
	val, ok := big.NewInt(0).SetString(string(out.Data), 10)
	require.True(ok)
	expected, ok := new(big.Int).SetString(accountAddressUnclaimedBalance, 10)
	require.True(ok)
	require.Equal(0, val.Cmp(expected))
}

func TestServer_GetReceiptByAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	actionHash := os.Getenv("actionHash")
	getReceiptByActionBlkHeight := os.Getenv("getReceiptByActionBlkHeight")
	getReceiptByActionBlkHeightInt, err := strconv.ParseUint(getReceiptByActionBlkHeight, 10, 64)
	request := &iotexapi.GetReceiptByActionRequest{ActionHash: actionHash}
	res, err := svr.GetReceiptByAction(request)
	require.NoError(err)
	require.NotNil(res)
	receiptPb := res.ReceiptInfo.Receipt
	require.Equal(uint64(1), receiptPb.Status)
	require.Equal(getReceiptByActionBlkHeightInt, receiptPb.BlkHeight)
	require.NotEqual(hash.ZeroHash256, res.ReceiptInfo.BlkHash)
}

func TestServer_ReadContract(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	readContractActionHash := os.Getenv("readContractActionHash")
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
	request2 := &iotexapi.ReadContractRequest{Action: res.ActionInfo[0].Action}

	res2, err := svr.ReadContract(request2)
	require.NoError(err)
	require.Equal("", res2.Data)
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
	epochDataHeight := os.Getenv("epochDataHeight")
	epochDataHeightInt, err := strconv.ParseUint(epochDataHeight, 10, 64)
	epochGravityChainStartHeight := os.Getenv("epochGravityChainStartHeight")
	epochGravityChainStartHeightInt, err := strconv.ParseUint(epochGravityChainStartHeight, 10, 64)

	res, err := svr.GetEpochMeta(&iotexapi.GetEpochMetaRequest{EpochNumber: 1})
	require.NoError(err)
	require.Equal(uint64(1), res.EpochData.Num)
	require.Equal(epochDataHeightInt, res.EpochData.Height)
	require.Equal(epochGravityChainStartHeightInt, res.EpochData.GravityChainStartHeight)
	require.Equal(360, int(res.TotalBlocks))
	require.Equal(24, len(res.BlockProducersInfo))
}
