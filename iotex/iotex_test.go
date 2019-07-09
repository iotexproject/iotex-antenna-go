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
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

const (
	_testnet           = "api.testnet.iotex.one:443"
	_mainnet           = "api.iotex.one:443"
	_accountPrivateKey = "9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
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
	hash, err := c.Transfer(to, v).SetGasPrice(big.NewInt(1)).SetGasLimit(1000000).Call(context.Background())
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

	hash, err := c.DeployContract(data).SetGasPrice(big.NewInt(1)).SetGasLimit(1000000).SetArgs(abi, big.NewInt(10)).Call(context.Background())
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

	hash, err := c.Contract(contract, abi).Execute("set", big.NewInt(8)).SetGasPrice(big.NewInt(1)).SetGasLimit(1000000).Call(context.Background())
	require.NoError(err)
	require.NotNil(hash)
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

	conn, err := NewDefaultGRPCConn(_mainnet)
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))

	h, err := hash.HexStringToHash256("f322902f272ee9f8d716ee843e7e1edd426ba4d37ef715b8430ee34f5b40bc75")
	require.NoError(err)
	res, err := c.GetReceipt(h).Call(context.Background())
	require.NoError(err)
	require.Equal("8715ffa12a776c6bfc9c04ba904ca5ff21c0a25e369ae1262a246b1e02b3e9cf", res.ReceiptInfo.BlkHash)
	require.Equal(uint64(1), res.ReceiptInfo.Receipt.Status)
	require.Equal(uint64(378583), res.ReceiptInfo.Receipt.BlkHeight)
	require.Equal("f322902f272ee9f8d716ee843e7e1edd426ba4d37ef715b8430ee34f5b40bc75", hex.EncodeToString(res.ReceiptInfo.Receipt.ActHash))
}

func TestGetExecutionResult(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_mainnet)
	require.NoError(err)
	defer conn.Close()

	c := NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))

	h, err := hash.HexStringToHash256("edf65e7ccbfb05e4fbd394db1acc276029c309994879e3a8c07023a753ea8886")
	require.NoError(err)
	res, err := c.GetExecutionResult(h).Call(context.Background())
	require.NoError(err)
	require.NotNil(res)
}
