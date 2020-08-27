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
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
)

const (
	_testnet           = "api.testnet.iotex.one:443"
	_mainnet           = "api.iotex.one:443"
	_accountPrivateKey = "73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426"
	_to                = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
)

func TestTransfer(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)

	to, err := address.FromString(_to)
	require.NoError(err)
	v := big.NewInt(1000000000000000000)
	hash, err := c.Transfer(to, v).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).Call(context.Background())
	require.NoError(err)
	require.NotEmpty(hash)
}

func TestClaimReward(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)

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
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)

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
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
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
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
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
	conn, err := NewDefaultGRPCConn(_mainnet)
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
			Count:     1000,
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
