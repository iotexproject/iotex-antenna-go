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

	"github.com/iotexproject/iotex-core/testutil"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	ta "github.com/iotexproject/iotex-core/test/testaddress"
)

const (
	//export accountPrivateKey              = "9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	//export accountAddress                 = "io14gnqxf9dpkn05g337rl7eyt2nxasphf5m6n0rd"
	//export accountBalance                 = "99994712164399999990848350"
	//export accountNonce                   = "337529"
	//export accountPendingNonce            = "337530"
	//export accountNumActions              = "506609"
	//export actionHash                     = "b9a938e1f249d3c7ab8152e377132989535e25ea9ee376323179d1943dc15b4e"
	//export actionActionInfoLen            = "1"
	//export actionActionNonce              = "1"
	//export getActionsByAddressActionHash  = "10efdacf68be2fa0afdc2a46786b3caf0a59fe2386485a9f075acf2a41c93d78"
	//export blk60801Hash                   = "a331841e6b29becdeeb65cb0948d896076b514e5f8d69d560b4424282a7882d7"
	//export blk60801HashNumActions         = "1960"
	//export blk60801HashTransferAmount     = "2389000000000000000000"
	//export getServerMetaPackageCommitID   = "e24aadf90f98d800b9f117354ddd5b3dbe58dde9"
	//export accountAddressUnclaimedBalance = "0"
	//export getReceiptByActionBlkHeight    = "517"
	//
	//export epochDataHeight              = "1"
	//export epochGravityChainStartHeight = "7502300"
	//export readContractActionHash       = "28a4b4979cb9922e18c60a1c4238cbea3775f757fb947456342b90eca7e52e08"
	//export senderPriKey = "8c379a71721322d16912a88b1602c5596ca9e99a5f70777561c3029efa71a435"
	host = "api.iotex.one:80"
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
	res, err := svr.GetAccount(accountAddress)
	require.NoError(err)
	accountMeta := res.AccountMeta
	require.Equal(accountAddress, accountMeta.Address)
	require.NotEqual(accountBalance, accountMeta.Balance)
	require.Equal(true, accountNonceInt < accountMeta.Nonce)
	require.Equal(true, accountPendingNonceInt < accountMeta.PendingNonce)
	require.Equal(true, accountNumActionsInt < accountMeta.NumActions)

	// failure
	_, err = svr.GetAccount("")
	require.Error(err)
}

func TestServer_GetActionsByIndex(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.GetActionsByIndex(0, 5)
	require.NoError(err)
	require.Equal(5, len(res.ActionInfo))
}

func TestServer_SendTransfer(t *testing.T) {
	require := require.New(t)
	rpc, err := NewRPCMethod(host)
	require.NoError(err)
	recipientAddr := os.Getenv("accountAddress")
	senderPriKey := os.Getenv("senderPriKey")
	_, err = rpc.SendTransfer(recipientAddr, senderPriKey, 10, big.NewInt(1000000), nil, 1000000, big.NewInt(10000000000))
	require.NoError(err)
}

func TestServer_GetActionByHash(t *testing.T) {
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
	res, err := svr.GetActionsByHash(actionHash, true)
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
	res, err := svr.GetActionsByAddress(accountAddress, 1, 1)
	require.NoError(err)
	require.Equal(getActionsByAddressActionHash, res.ActionInfo[0].ActHash)
	require.Equal(1, len(res.ActionInfo))
}

func TestServer_GetUnconfirmedActionsByAddress(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	accountAddress := os.Getenv("accountAddress")
	res, err := svr.GetUnconfirmedActionsByAddress(accountAddress, 1, 10)
	require.NoError(err)
	require.Equal(0, len(res.ActionInfo))
}

func TestServer_GetActionsByBlock(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	blk60801Hash := os.Getenv("blk60801Hash")
	res, err := svr.GetActionsByBlock(blk60801Hash, 1, 10)
	require.NoError(err)
	require.Equal(10, len(res.ActionInfo))
}

func TestServer_GetBlockMetas(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.GetBlockMetasByIndexAndCount(10, 10)
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

	res, err := svr.GetBlockMetasByBlockHash(blk60801Hash)
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

	res, err := svr.GetChainMeta()
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
	res, err := svr.GetServerMeta()
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
	out, err := svr.ReadState("rewarding", "UnclaimedBalance", [][]byte{[]byte(accountAddress)})
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
	res, err := svr.GetReceiptByAction(actionHash)
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
	res2, err := svr.ReadContract(readContractActionHash, true)
	require.NoError(err)
	require.Equal("", res2.Data)
}

func TestServer_SuggestGasPrice(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.SuggestGasPrice()
	require.NoError(err)
	require.Equal(uint64(1), res.GasPrice)
}

func TestServer_EstimateGasForAction(t *testing.T) {
	require := require.New(t)
	svr, err := NewRPCMethod(host)
	require.NoError(err)
	res, err := svr.EstimateGasForAction(ta.Addrinfo["alfa"].String(), ta.Keyinfo["alfa"].PriKey.HexString(), 3, big.NewInt(10000000), 1000000, big.NewInt(testutil.TestGasPriceInt64), nil)
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

	res, err := svr.GetEpochMeta(1)
	require.NoError(err)
	require.Equal(uint64(1), res.EpochData.Num)
	require.Equal(epochDataHeightInt, res.EpochData.Height)
	require.Equal(epochGravityChainStartHeightInt, res.EpochData.GravityChainStartHeight)
	require.Equal(360, int(res.TotalBlocks))
	require.Equal(24, len(res.BlockProducersInfo))
}
