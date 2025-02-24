// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
)

const (
	_testnet           = "api.testnet.iotex.one:443"
	_mainnet           = "api.iotex.one:443"
	_accountPrivateKey = "73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426"
	_to                = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
	_testNetChainID    = 4690
)

func TestTransfer(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)

	to, err := address.FromString(_to)
	require.NoError(err)
	v := big.NewInt(160000000000000000)
	for _, test := range []struct {
		chainID uint32
		err     string
	}{
		{0, "0 is not a valid chain ID (use 1 for mainnet, 2 for testnet)"},
		{1, "ChainID does not match, expecting 2, got 1"},
		{2, ""},
	} {
		c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), test.chainID, acc)
		caller := c.Transfer(to, v)
		hash, err := caller.SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).Call(context.Background())
		if len(test.err) > 0 {
			require.Contains(err.Error(), test.err)
		} else {
			require.NoError(err)
			require.NotEmpty(hash)
		}
	}
}

func TestStake(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount("afeefca74d9a325cf1d6b6911d61a65c32afa8e02bd5e78e2e4ac2910bab45f5")
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc)

	one := big.NewInt(int64(unit.Iotx))
	_, err = c.Staking().Create("robotbp00001", one.Lsh(one, 7), 0, false).
		SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(20000).Call(context.Background())
	require.Contains(err.Error(), "insufficient funds for gas * price + value")
}

func TestClaimReward(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc)

	require.NoError(err)
	v := big.NewInt(1000000000000000000)
	hash, err := c.ClaimReward(v).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).Call(context.Background())
	require.NoError(err)
	require.NotEmpty(hash)
}

func TestDeployContract(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc)

	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`))
	require.NoError(err)
	data, err := hex.DecodeString("608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a723058208d4f6c9737f34d9b28ef070baa8127c0876757fbf6f3945a7ea8d4387ca156590029")
	require.NoError(err)

	hash, err := c.DeployContract(data).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).SetArgs(abi, big.NewInt(10)).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
}

func TestExecuteContract(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc)
	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`))
	require.NoError(err)
	contract, err := address.FromString("io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0")
	require.NoError(err)

	hash, err := c.Contract(contract, abi).Execute("set", big.NewInt(8)).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
}

func TestExecuteContractWithAddressArgument(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 2, acc)
	abi, err := abi.JSON(strings.NewReader(`[
	{
		"constant": false,
		"inputs": [
			{
				"name": "recipients",
				"type": "address[]"
			},
			{
				"name": "amounts",
				"type": "uint256[]"
			},
			{
				"name": "payload",
				"type": "string"
			}
		],
		"name": "multiSend",
		"outputs": [],
		"payable": true,
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "recipient",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "refund",
				"type": "uint256"
			}
		],
		"name": "Refund",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "payload",
				"type": "string"
			}
		],
		"name": "Payload",
		"type": "event"
	}
]`))
	require.NoError(err)
	contract, err := address.FromString("io1up8gd9nxhc0k0fjff7nrl6jn626vkdzj7y3g09")
	require.NoError(err)

	recipient1, err := address.FromString("io18jaldgzc8wlyfnzamgas62yu3kg5nw527czg37")
	require.NoError(err)
	recipient2, err := address.FromString("io1ntprz4p5zw38fvtfrcczjtcv3rkr3nqs6sm3pj")
	require.NoError(err)

	recipients := [2]address.Address{recipient1, recipient2}
	// recipients := [2]string{"io18jaldgzc8wlyfnzamgas62yu3kg5nw527czg37", "io1ntprz4p5zw38fvtfrcczjtcv3rkr3nqs6sm3pj"}
	amounts := [2]*big.Int{big.NewInt(1), big.NewInt(2)}
	actionHash, err := c.Contract(contract, abi).Execute("multiSend", recipients, amounts, "payload").SetGasPrice(big.NewInt(1000000000000)).SetGasLimit(1000000).Call(context.Background())
	require.NoError(err)
	require.NotNil(actionHash)
}

func TestReadContract(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))

	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`))
	require.NoError(err)
	contract, err := address.FromString("io17sn486alutrnzlrdz9vv44g7qyc38hygf7s6h0")
	require.NoError(err)

	_, err = c.ReadOnlyContract(contract, abi).Read("get").Call(context.Background())
	require.NoError(err)
}

func TestGetReceipt(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))

	actionHash, err := hash.HexStringToHash256("163ece70353acfe8fa7929e756d96b1b3cfec1246bc5a8f397ca77f20a0d5c5f")
	require.NoError(err)
	_, err = c.GetReceipt(actionHash).Call(context.Background())
	require.NoError(err)
}

func TestGetRlpTx(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	c := iotexapi.NewAPIServiceClient(conn)
	ctx := context.Background()
	req := iotexapi.GetRawBlocksRequest{
		Count:               1,
		WithTransactionLogs: true,
	}
	for _, v := range []struct {
		height uint64
		hash   string
	}{
		{8638208, "bc34eaac7fc9a4e3edc78aff514fa5631f722702ed35b4f696229d5da3f7e914"},
		{8638638, "a43e8ee6c444da9b7524663afdde84dd3c4976f0127a3fbc3129e18480213386"},
		{8638658, "76c2b74ea767529a2ecc45721e968b9a75c930be91742ac84dd200375af5ab76"},
	} {
		req.StartHeight = v.height
		res, err := c.GetRawBlocks(ctx, &req)
		require.NoError(err)
		require.Equal(1, len(res.Blocks))
		txLog := res.Blocks[0].TransactionLogs
		require.Equal(1, len(txLog.Logs))
		log := txLog.Logs[0]
		require.Equal(v.hash, hex.EncodeToString(log.ActionHash))
		require.Equal(2, len(log.Transactions))
		for i, tx := range log.Transactions {
			if i == 0 {
				// first log is for native transfer
				require.Equal(
					iotextypes.TransactionLogType_name[int32(iotextypes.TransactionLogType_NATIVE_TRANSFER)],
					tx.Type.String())
			} else {
				// second log is for gas fee in the amount of 0.01 IOTX
				require.Equal(
					iotextypes.TransactionLogType_name[int32(iotextypes.TransactionLogType_GAS_FEE)],
					tx.Type.String())
				require.Equal("10000000000000000", tx.Amount)
			}
		}

		// verify hash
		act := res.Blocks[0].GetBlock().GetBody().GetActions()[0]
		require.Equal(act.Encoding, iotextypes.Encoding_ETHEREUM_RLP)
		h, err := ActionHash(act, _testNetChainID)
		require.NoError(err)
		require.Equal(h[:], log.ActionHash)
	}
}

func TestGetLogs(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))

	// https://testnet.iotexscan.io/action/22cd0c2d1f7d65298cec7599e2d0e3c650dd8b4ed2b1c816d909026c60d785b2
	name, _ := hex.DecodeString("000000000000000000000000000000000000007472616e736665725374616b65")
	index, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000039")
	receiver, _ := hex.DecodeString("000000000000000000000000cb68a8318de4d4061e0956de69927c327bcfb352")
	sender, _ := hex.DecodeString("00000000000000000000000053fbc28faf9a52dfe5f591948a23189e900381b5")
	filterTopics := [][]byte{name, index, receiver, sender}
	blkHash, _ := hex.DecodeString("b13199e4cc712b3fee4feda52e39cec664ef5cbbc775ee1a66535305ff3a1af7")

	req := &iotexapi.GetLogsRequest{
		Filter: &iotexapi.LogsFilter{
			Address: []string{"io1qnpz47hx5q6r3w876axtrn6yz95d70cjl35r53"},
			Topics: []*iotexapi.Topics{
				&iotexapi.Topics{},
				&iotexapi.Topics{Topic: filterTopics},
			},
		},
		Lookup: &iotexapi.GetLogsRequest_ByBlock{
			ByBlock: &iotexapi.GetLogsByBlock{
				BlockHash: blkHash,
			},
		},
	}
	logs, err := c.GetLogs(req).Call(context.Background())
	require.NoError(err)
	require.Equal(1, len(logs.Logs))
	log := logs.Logs[0]
	for i := range filterTopics {
		require.Equal(filterTopics[i], log.Topics[i])
	}

	// pulling with second topic (bucket ID) = 0x31 in 500 blocks
	index, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000031")
	req.Filter.Topics = []*iotexapi.Topics{
		&iotexapi.Topics{},
		&iotexapi.Topics{Topic: [][]byte{index}},
	}
	req.Lookup = &iotexapi.GetLogsRequest_ByRange{
		ByRange: &iotexapi.GetLogsByRange{
			FromBlock: 4795567,
			ToBlock:   4796567,
		},
	}
	logs, err = c.GetLogs(req).Call(context.Background())
	require.NoError(err)
	require.Equal(5, len(logs.Logs))

	// verify index == 0x31
	for _, log := range logs.Logs {
		require.True(len(log.Topics) >= 2)
		require.Equal(index, log.Topics[1])
	}
}

// mainnet tests
func mainnetGrpcConn() (*grpc.ClientConn, error) {
	return NewDefaultGRPCConn(_mainnet)
}

func TestMainnetGetActions(t *testing.T) {
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	ctx := context.Background()
	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	res, err := c.API().GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash: "84755e8be7ee7ba7891a57cf03b3c3dc65d2699b8935552704b900c70f38d4af",
			},
		},
	})
	require.NoError(err)
	require.Equal(1, len(res.ActionInfo))
	act := res.ActionInfo[0]
	require.Equal("84755e8be7ee7ba7891a57cf03b3c3dc65d2699b8935552704b900c70f38d4af", act.ActHash)
	require.Equal("7e42c615c7dde2297c98220520eb68b0af9d1e6ff856b793974a311aed8221b5", act.BlkHash)
	require.EqualValues(8007, act.BlkHeight)
	require.Equal("io17ch0jth3dxqa7w9vu05yu86mqh0n6502d92lmp", act.Sender)
	require.Equal("11493000000000000", act.GasFee)

	res, err = c.API().GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByIndex{
			ByIndex: &iotexapi.GetActionsByIndexRequest{
				Start: 0,
				Count: 5,
			},
		},
	})
	require.NoError(err)
	require.Equal(5, len(res.ActionInfo))
	for i, v := range []struct {
		height                        uint64
		actHash, blkHash, sender, gas string
	}{
		{1, "dd2e83336f1ff219b1e54558f0627e1f556ed2caeedb44b758b0e107aa246531", "230ba8095d5a505e355652f9dcc2b13605419a8fa3d3fd5ddc6d24fd6a902641", "io1vtm2zgn830pn6auc2cvnchgwdaefa9gr4z0s86", "0"},
		{1, "b7024bc52f315fafb9cc7677e730aec79767b28fbaa6bdd1f37c1861dd699aba", "230ba8095d5a505e355652f9dcc2b13605419a8fa3d3fd5ddc6d24fd6a902641", "io1vtm2zgn830pn6auc2cvnchgwdaefa9gr4z0s86", "0"},
		{2, "833de10c30ffbedd898ae8669123362bdd1f4012ac0b979d784f201364d4dda0", "e6bdc2fd1d36f47ec8f9c5554503855eb29453a7cd4138b3ceb1af670fee6d75", "io1jafqlvntcxgyp6e0uxctt3tljzc3vyv5hg4ukh", "0"},
		{3, "dbbd4c606b7743be3e064aa84e78e70dc786d1599c9bf37c03681dc8038b761e", "cb6055d916355600471f1ca2cf15098d56c7f59cff3dae0123228c8a2da8c1f2", "io1v0q5g2f8z6l3v25krl677chdx7g5pwt9kvqfpc", "0"},
		{4, "ddd8149259ae3f169cf0f84ae8bb420089860a9ce2cb7a79c28edd475e002ae8", "67b9de1ffda9995c033e37f2d3dabdfab798aa240908f02dc8096037f12613da", "io1nyjs526mnqcsx4twa7nptkg08eclsw5c2dywp4", "0"},
	} {
		act := res.ActionInfo[i]
		require.Equal(v.actHash, act.ActHash)
		require.Equal(v.blkHash, act.BlkHash)
		require.Equal(v.height, act.BlkHeight)
		require.Equal(v.sender, act.Sender)
		require.Equal(v.gas, act.GasFee)
	}

	res, err = c.API().GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByBlk{
			ByBlk: &iotexapi.GetActionsByBlockRequest{
				BlkHash: "e222dfbd059e064369b65f6553455329a1640810c88776154c8db08071b8c860",
				Start:   0,
				Count:   4,
			},
		},
	})
	require.NoError(err)
	require.Equal(3, len(res.ActionInfo))
	for i, v := range []string{
		"70714000000000000",
		"304233000000000000",
		"0",
	} {
		act = res.ActionInfo[i]
		require.Equal("e222dfbd059e064369b65f6553455329a1640810c88776154c8db08071b8c860", act.BlkHash)
		require.EqualValues(14527394, act.BlkHeight)
		require.Equal(v, act.GasFee)
	}
}

func TestMainnetGetBlock(t *testing.T) {
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	req := iotexapi.GetRawBlocksRequest{
		Count:        1,
		WithReceipts: true,
	}
	type blk struct {
		height                             uint64
		prev, txRoot, delta, rxRoot, bloom string
		actSize, rxSize                    int
	}
	for _, v := range []blk{
		{8007, "857a6e2adb32020fc700c5bb89d4eebbb173c3bd763dfdb2be3a6cf4067a05bc", "ead1f54f396bf6016bf27f45e5252cf26c1a2e5c2a69b8f7e066754d45b5f1ba", "1f8b25e0333e9794e6a61d783b66822f4efeff010a15be2f677ed02e74deea28", "fee660dcce2cf0f4f6fc09b2730c5df1a0d5ebe68d6c410b58748bedc17220c5", "", 2, 2},
		{13713788, "e3c2eb8aaf531b877fa461aeb000209bbe95ac17626237952db9f25c7e89305f", "17caa398ed77b50329f90f8c9d72f423f435b29866d4295568370759cf596e69", "91ba29f363d970eb4022dade4229c7d8beadeb30e7292d56a9376fac11bb6836", "f4b59d54db803c1590428e90513a17427525f8ebfed0e5b9a10dd2cd5ad317c1", "00000000080004000081200000000100000000000000000000000000000000100000001000000200002000000000000001800080000800800000008000000000000000060000000000000800000000000082000088100800000040000000000000000a0040000000000000400000000000000000000000801000000000000000100000010000000008020000000000000000000000010000000000000000000080000000000000040000100a000800000000000000000000000000011000200000000000000000002808000010000000000000000000000800000000200000100000000050001000000000000200000000000000000001010000100010000000", 6, 6},
	} {
		req.StartHeight = v.height
		res, err := c.API().GetRawBlocks(context.Background(), &req)
		require.NoError(err)
		require.Equal(1, len(res.Blocks))
		blk := res.Blocks[0]
		require.Equal(v.rxSize, len(blk.Receipts))
		core := blk.Block.Header.Core
		require.Equal(v.height, core.Height)
		require.Equal(v.prev, hex.EncodeToString(core.PrevBlockHash))
		require.Equal(v.txRoot, hex.EncodeToString(core.TxRoot))
		require.Equal(v.delta, hex.EncodeToString(core.DeltaStateDigest))
		require.Equal(v.rxRoot, hex.EncodeToString(core.ReceiptRoot))
		require.Equal(v.bloom, hex.EncodeToString(core.LogsBloom))
		require.Equal(v.actSize, len(blk.Block.Body.Actions))
	}
}

func TestMainnetGetMeta(t *testing.T) {
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	ctx := context.Background()
	c := iotexapi.NewAPIServiceClient(conn)
	res, err := c.GetBlockMetas(ctx, &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: 13655501,
				Count: 100,
			},
		},
	})
	require.NoError(err)
	require.Equal(100, len(res.BlkMetas))
	meta := res.BlkMetas[0]
	require.Equal(meta.Hash, "fd7b5a375300940fe0c0da482ee79a2e04440c4dc0490f2fc57e191b9c322d71")
	require.EqualValues(meta.Height, 13655501)
	require.EqualValues(meta.NumActions, 8)
	require.Equal(meta.ProducerAddress, "io1zy9q4myp3jezxdsv82z3f865ruamhqm7a6mggh")
	require.Equal(meta.TransferAmount, "0")
	require.Equal(meta.TxRoot, "2baa264c7bbaf4b40ab8852e9ca8c900ee7c0d50425738dd337d3831a8b332d9")
	require.Equal(meta.ReceiptRoot, "55d4e59f7a1f2c070300ff5eae84fcbeacd80970226dc8cb754bbbb916311878")
	require.Equal(meta.DeltaStateDigest, "cba67830d9a1f5d88da1bcafa7c9e79a780f0a3cf8af89b93cce50527a81c441")
	require.Equal(meta.LogsBloom, "00100501000000200000000024000000400004000c00008080000000000000000000000000000000000000010000000002800000000020000000000100400200000000000000410000200080804000000002000000000000000040400040004000028200000040c200000140001000000000000000002021000200000200020010000020000000000802040000000000110004000200020020005010000000040004200000000000000000002008000001000000004020000000000800302000000000000000000080200000000000000000000020020100000000000800010000480000c0000000000000400000000000020008080000820000000010000880")
	require.EqualValues(meta.GasLimit, 21208628)
	require.EqualValues(meta.GasUsed, 953375)
}

func TestMainnetGetReceipt(t *testing.T) {
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	for _, v := range []struct {
		actHash, blkHash, contract string
		status, height, gas        uint64
		logSize                    int
		logData                    string
	}{
		{"84755e8be7ee7ba7891a57cf03b3c3dc65d2699b8935552704b900c70f38d4af", "7e42c615c7dde2297c98220520eb68b0af9d1e6ff856b793974a311aed8221b5", "", 0, 8007, 11493, 0, ""},
		{"cd93427afc41e0e219b1aa6370b04ec89989e249438c7f8451a83f48503a3660", "7e42c615c7dde2297c98220520eb68b0af9d1e6ff856b793974a311aed8221b5", "io154mvzs09vkgn0hw6gg3ayzw5w39jzp47f8py9v", 1, 8007, 0, 1, "1229696f3163357a776832347063347a3837747177346d367a3663347935343470777376386e72726d36361a143136303030303030303030303030303030303030"},
	} {
		h, err := hex.DecodeString(v.actHash)
		require.NoError(err)
		res, err := c.GetReceipt(hash.BytesToHash256(h)).Call(context.Background())
		require.NoError(err)
		require.Equal(v.blkHash, res.ReceiptInfo.BlkHash)
		r := res.ReceiptInfo.Receipt
		require.Equal(v.status, r.Status)
		require.Equal(v.height, r.BlkHeight)
		require.Equal(v.actHash, hex.EncodeToString(r.ActHash))
		require.Equal(v.gas, r.GasConsumed)
		require.Equal(v.contract, r.ContractAddress)
		require.Equal(v.logSize, len(r.Logs))
		if v.logSize > 0 {
			// verify first log
			l := r.Logs[0]
			require.Equal(v.contract, l.ContractAddress)
			require.Zero(len(l.Topics))
			require.Equal(v.logData, hex.EncodeToString(l.Data))
			require.Equal(v.height, l.BlkHeight)
			require.Equal(v.actHash, hex.EncodeToString(l.ActHash))
			require.Zero(l.BlkHash) // somehow defined but not used
		}
	}
}

func TestMainnetTransactionLogs(t *testing.T) {
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	req := iotexapi.GetTransactionLogByBlockHeightRequest{
		BlockHeight: 11399317,
	}
	res, err := c.API().GetTransactionLogByBlockHeight(context.Background(), &req)
	require.NoError(err)
	require.Equal(1, len(res.TransactionLogs.Logs))
	log := res.TransactionLogs.Logs[0]
	require.Equal("fa45221706f521e3c18638610c52d45166ae5d91b8c23d89a5f4a49cba5e5b1c", hex.EncodeToString(log.ActionHash))
	require.EqualValues(2, log.NumTransactions)
	type txLog struct {
		send, recv, amount string
		txType             int
	}
	for i, v := range []txLog{
		{"io1yrnwxklgkyu464zlpgx9laa6ax4203znvkc3z8", "io0000000000000000000000rewardingprotocol", "29999000000000000", 6},
		{"io1yrnwxklgkyu464zlpgx9laa6ax4203znvkc3z8", "io1qp00xg6a0x864a7k68l3wm7l9zggk5wv3wwtem", "0", 0},
	} {
		tx := log.Transactions[i]
		require.Equal(v.send, tx.Sender)
		require.Equal(v.recv, tx.Recipient)
		require.Equal(v.amount, tx.Amount)
		require.EqualValues(v.txType, tx.Type)
	}

	req.BlockHeight = 15146543
	res, err = c.API().GetTransactionLogByBlockHeight(context.Background(), &req)
	require.NoError(err)
	require.Equal(4, len(res.TransactionLogs.Logs))
	log = res.TransactionLogs.Logs[0]
	require.Equal("d814942989e07a87c9a6258b369d891b6abf0dc0b10a879c6965d4d168a7c43b", hex.EncodeToString(log.ActionHash))
	require.EqualValues(1, log.NumTransactions)
	tx := log.Transactions[0]
	require.Equal("io18z6st3gkxfda3t7dhksk63cgd7dr7jnkwwy22f", tx.Sender)
	require.Equal("io0000000000000000000000rewardingprotocol", tx.Recipient)
	require.Equal("155692000000000000", tx.Amount)
	require.EqualValues(6, tx.Type)

	log = res.TransactionLogs.Logs[1]
	require.Equal("85e65a6062bc3c35c55562bb49bc39f80fdd55a570ff9adb717f566f7d9625c4", hex.EncodeToString(log.ActionHash))
	require.EqualValues(2, log.NumTransactions)
	for i, v := range []txLog{
		{"io1lhukp867ume3qn2g7cxn4e47pj0ugfxeqj7nm8", "io1fur9ctq05wms488l3hjn2s59lpqq34wuakqv88", "61415308145556742", 7},
		{"io1lhukp867ume3qn2g7cxn4e47pj0ugfxeqj7nm8", "io0000000000000000000000rewardingprotocol", "10000000000000000", 6},
	} {
		tx := log.Transactions[i]
		require.Equal(v.send, tx.Sender)
		require.Equal(v.recv, tx.Recipient)
		require.Equal(v.amount, tx.Amount)
		require.EqualValues(v.txType, tx.Type)
	}

	log = res.TransactionLogs.Logs[2]
	require.Equal("d67f964c1eef51315d7484179ea0fe3529e5fa5ec7e985c8854814ffac44ed0c", hex.EncodeToString(log.ActionHash))
	require.EqualValues(2, log.NumTransactions)
	for i, v := range []txLog{
		{"io1ayw5k53sh4pfetmgj887yzuhhge42fxzruga6l", "io18q8zqlwwtkr6wwjd2smrqlfjza29mn5c4wnr46", "14650000000000000000000", 7},
		{"io1ayw5k53sh4pfetmgj887yzuhhge42fxzruga6l", "io0000000000000000000000rewardingprotocol", "10000000000000000", 6},
	} {
		tx := log.Transactions[i]
		require.Equal(v.send, tx.Sender)
		require.Equal(v.recv, tx.Recipient)
		require.Equal(v.amount, tx.Amount)
		require.EqualValues(v.txType, tx.Type)
	}

	log = res.TransactionLogs.Logs[3]
	require.Equal("44b857ed0514df6474c2af7e14ec09d1414585019b536f183a2e192cea71db87", hex.EncodeToString(log.ActionHash))
	require.EqualValues(1, log.NumTransactions)
	tx = log.Transactions[0]
	require.Equal("io1dyj80tp303jshhvumlfkfemclqnmnr6hng95ep", tx.Sender)
	require.Equal("io0000000000000000000000rewardingprotocol", tx.Recipient)
	require.Equal("85714000000000000", tx.Amount)
	require.EqualValues(6, tx.Type)
}

func TestMainnetReadCandAndReclaim(t *testing.T) {
	t.Skipf("skip mainnet test because the data may be updated")
	require := require.New(t)
	conn, err := mainnetGrpcConn()
	require.NoError(err)
	defer conn.Close()

	c := iotexapi.NewAPIServiceClient(conn)
	ctx := context.Background()
	owner := "io16frg0hz9ztxhqkn5e42f2aznnreyy79dmqlsn2"
	readStakingDataRequest := &iotexapi.ReadStakingDataRequest{
		Request: &iotexapi.ReadStakingDataRequest_CandidateByAddress_{
			CandidateByAddress: &iotexapi.ReadStakingDataRequest_CandidateByAddress{
				OwnerAddr: owner,
			},
		},
	}
	requestData, err := proto.Marshal(readStakingDataRequest)
	require.NoError(err)
	method := &iotexapi.ReadStakingDataMethod{
		Method: iotexapi.ReadStakingDataMethod_CANDIDATE_BY_ADDRESS,
	}
	methodData, err := proto.Marshal(method)
	require.NoError(err)
	rep, err := c.ReadState(ctx, &iotexapi.ReadStateRequest{
		ProtocolID: []byte("staking"),
		MethodName: methodData,
		Arguments:  [][]byte{requestData},
	})
	require.NoError(err)
	cand := iotextypes.CandidateV2{}
	require.NoError(proto.Unmarshal(rep.Data, &cand))
	require.Equal("coredev", cand.Name)
	require.Equal(owner, cand.OwnerAddress)
	require.Equal("io1hgxksz37qtqq9n5n6lkkc9qhaajdklhvkh5969", cand.OperatorAddress)
	require.Equal("io12mgttmfa2ffn9uqvn0yn37f4nz43d248l2ga85", cand.RewardAddress)
	require.EqualValues(11, cand.SelfStakeBucketIdx)

	res, err := c.GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash: "2b2419b09aeaa6e650359dea8893083b2befa7a541adce4f73504e7c380f23ed",
			},
		},
	})
	require.NoError(err)
	act := res.ActionInfo[0].Action.Core.GetStakeTransferOwnership()
	require.NotNil(act)
	reclaim := string(act.Payload)
	require.Contains(reclaim, "\"{\\\"bucket\\\":1339")
	require.Contains(reclaim, "\\\"nonce\\\":1")
	require.Contains(reclaim, "\\\"recipient\\\":\\\"io10ek5h563632nu5p2ur3ysj79ctxca9zh75eayv\\\"")
	require.Contains(reclaim, "caaade52f022eb0e50693b7b4997dbacdb0b5b1547b095b0870e2441e2f7c72c72699ba8024261f32f2736736d39a782863f6b13dcae632bb8681d820b9d232101")
}
