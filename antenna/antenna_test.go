// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package antenna

import (
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	ta "github.com/iotexproject/iotex-core/test/testaddress"
	"github.com/iotexproject/iotex-core/testutil"
)

const (
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
	request := &GetAccountRequest{Address: accountAddress}
	res, err := svr.GetAccount(request)
	require.NoError(err)
	accountMeta := res.AccountMeta
	require.Equal(accountAddress, accountMeta.Address)
	require.NotEqual(accountBalance, accountMeta.Balance)
	require.Equal(true, accountNonceInt < accountMeta.Nonce)
	require.Equal(true, accountPendingNonceInt < accountMeta.PendingNonce)
	require.Equal(true, accountNumActionsInt < accountMeta.NumActions)

	// failure
	_, err = svr.GetAccount(&GetAccountRequest{})
	require.Error(err)
}

func TestServer_GetActions(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_ByIndex{
			ByIndex: &GetActionsByIndexRequest{
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
	request := &SendActionRequest{Action: testTransferPb}
	_, err = rpc.SendAction(request)
	require.NoError(err)
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
	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_ByHash{
			ByHash: &GetActionByHashRequest{
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
	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_ByAddr{
			ByAddr: &GetActionsByAddressRequest{
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
	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_UnconfirmedByAddr{
			UnconfirmedByAddr: &GetUnconfirmedActionsByAddressRequest{
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
	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_ByBlk{
			ByBlk: &GetActionsByBlockRequest{
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

	request := &GetBlockMetasRequest{
		Lookup: &GetBlockMetasRequest_ByIndex{
			ByIndex: &GetBlockMetasByIndexRequest{
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

	request := &GetBlockMetasRequest{
		Lookup: &GetBlockMetasRequest_ByHash{
			ByHash: &GetBlockMetaByHashRequest{
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

	res, err := svr.GetChainMeta(&GetChainMetaRequest{})
	require.NoError(err)
	chainMetaPb := res.ChainMeta
	require.Equal(true, chainMetaPb.Height > 1)
	require.Equal(true, chainMetaPb.NumActions > 1)
	require.Equal(true, chainMetaPb.Tps >= 0)
	require.Equal(true, chainMetaPb.Epoch.Num > 1)
	require.Equal(true, chainMetaPb.Epoch.Height > 1)
	require.Equal(true, chainMetaPb.Epoch.GravityChainStartHeight > 1)
}

func TestServer_GetServerMeta(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.GetServerMeta(&GetServerMetaRequest{})
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
	out, err := svr.ReadState(&ReadStateRequest{
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
	request := &GetReceiptByActionRequest{ActionHash: actionHash}
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
	expectedReadContractActionHash := os.Getenv("expectedReadContractActionHash")
	request := &GetActionsRequest{
		Lookup: &GetActionsRequest_ByHash{
			ByHash: &GetActionByHashRequest{
				ActionHash:   readContractActionHash,
				CheckPending: true,
			},
		},
	}
	res, err := svr.GetActions(request)
	require.NoError(err)
	request2 := &ReadContractRequest{Action: res.ActionInfo[0].Action}

	res2, err := svr.ReadContract(request2)
	require.NoError(err)
	require.Equal(expectedReadContractActionHash, res2.Data)
}

func TestServer_SuggestGasPrice(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.SuggestGasPrice(&SuggestGasPriceRequest{})
	require.NoError(err)
	require.Equal(uint64(1), res.GasPrice)
}

func TestServer_EstimateGasForAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)

	request := &EstimateGasForActionRequest{Action: testTransferPb}
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

	res, err := svr.GetEpochMeta(&GetEpochMetaRequest{EpochNumber: 1})
	require.NoError(err)
	require.Equal(uint64(1), res.EpochData.Num)
	require.Equal(epochDataHeightInt, res.EpochData.Height)
	require.Equal(epochGravityChainStartHeightInt, res.EpochData.GravityChainStartHeight)
	require.Equal(360, int(res.TotalBlocks))
	require.Equal(24, len(res.BlockProducersInfo))
}
func TestFromRau(t *testing.T) {
	require := require.New(t)
	raw, ok := new(big.Int).SetString("12000000000000000000", 10)
	require.True(ok)
	convert := FromRau(raw)
	require.Equal(int64(12), convert)
}
func TestToRau(t *testing.T) {
	require := require.New(t)
	raw := int64(10)
	convert := ToRau(raw)
	expected, ok := new(big.Int).SetString("10000000000000000000", 10)
	require.True(ok)
	require.Equal(0, convert.Cmp(expected))
}
