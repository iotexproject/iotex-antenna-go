package iotex

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
)

func TestStake(t *testing.T) {
	require := require.New(t)
	conn, err := NewDefaultGRPCConn(_testnet)
	require.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(_accountPrivateKey)
	require.NoError(err)
	c := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	name, _ := address.FromString("io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he")
	operator, _ := address.FromString("io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he")
	reward, _ := address.FromString("io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he")

	// CandidateRegister
	ret, err := c.Candidate().Register(name, operator, reward, big.NewInt(1), 10, false).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).SetPayload([]byte("payload")).Call(context.Background())
	require.NoError(err)
	require.NotEqual(hash.ZeroHash256, ret)

	// need to fix when testnet ready
	//time.Sleep(time.Second * 20)
	//receipt, err := c.GetReceipt(ret).Call(context.Background())
	//require.NoError(err)
	//require.Equal(iotextypes.ReceiptStatus_Success, receipt.ReceiptInfo.Receipt.Status)

	// StakeCreate
	ret, err = c.Staking().Create("io19d0p3ah4g8ww9d7kcxfq87yxe7fnr8rpth5shj", big.NewInt(100), uint32(10000), true).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).SetPayload([]byte("payload")).Call(context.Background())
	require.NoError(err)
	require.NotEqual(hash.ZeroHash256, ret)

	// need to fix when testnet ready
	//time.Sleep(time.Second * 10)
	//receipt, err = c.GetReceipt(hash).Call(context.Background())
	//require.NoError(err)
	//require.Equal(iotextypes.ReceiptStatus_Success, receipt.ReceiptInfo.Receipt.Status)

	// StakeUnstake
	ret, err = c.Staking().Unstake(1).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).SetPayload([]byte("payload")).Call(context.Background())
	require.NoError(err)
	require.NotEqual(hash.ZeroHash256, ret)
	// need to fix when testnet ready
	time.Sleep(time.Second * 10)
	_, err = c.GetReceipt(ret).Call(context.Background())
	require.Error(err)
	//require.NotEqual(iotextypes.ReceiptStatus_Success, receipt.ReceiptInfo.Receipt.Status)
}
