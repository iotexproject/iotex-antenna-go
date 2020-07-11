package iotex

import (
	"context"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
)

var (
	testActionHash1, _ = hash.HexStringToHash256("16fbdac39a19f433b32e457eba9c64c48d9fafcdffaddb11a0b7d264c0cf1418")
	// CandidateRegister tests.
	candidateRegisterTests = []struct {
		// input
		name      string
		operator  string
		reward    string
		amount    *big.Int
		duration  uint32
		autoStake bool
		gasPrice  *big.Int
		gasLimit  uint64
		payload   []byte
		// expect
		actionHash hash.Hash256
	}{
		{
			"io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he",
			"io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he",
			"io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he",
			big.NewInt(int64(unit.Iotx)),
			10000,
			false,
			big.NewInt(int64(10 * unit.Qev)),
			1000000,
			[]byte("TestCandidateRegister"),
			testActionHash1,
		},
	}
	// Stake tests.
	stakeTests = []struct {
		candidateName string
		amount        *big.Int
		duration      uint32
		autoStake     bool
		gasPrice      *big.Int
		gasLimit      uint64
		payload       []byte
		// expect
		actionHash hash.Hash256
	}{
		{
			"io10a298zmzvrt4guq79a9f4x7qedj59y7ery84he",
			big.NewInt(int64(unit.Iotx)),
			10000,
			false,
			big.NewInt(int64(10 * unit.Qev)),
			1000000,
			[]byte("TestStaking"),
			testActionHash1,
		},
	}
	// Unstake tests.
	unstakeTests = []struct {
		bucketIndex uint64
		autoStake   bool
		gasPrice    *big.Int
		gasLimit    uint64
		payload     []byte
		// expect
		actionHash hash.Hash256
	}{
		{
			10000,
			false,
			big.NewInt(int64(10 * unit.Qev)),
			1000000,
			[]byte("TestStaking"),
			testActionHash1,
		},
	}
)

func TestCandidateCaller_Register(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(candidateRegisterTests)

	stakingAPICaller := NewMockStakingAPICaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	candidateCaller := NewMockCandidateCaller(ctrl)
	candidateCaller.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	client := NewMockAuthedClient(ctrl)
	client.EXPECT().Candidate().Return(candidateCaller).Times(testTimes)

	for _, test := range candidateRegisterTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, nil).Times(1)

		nameAddr, err := address.FromString(test.name)
		require.NoError(err)
		operatorAddr, err := address.FromString(test.operator)
		require.NoError(err)
		rewardAddr, err := address.FromString(test.reward)
		require.NoError(err)
		ret, err := client.Candidate().
			Register(nameAddr, operatorAddr, rewardAddr, test.amount, test.duration, test.autoStake).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.NoError(err)
		require.Equal(test.actionHash, ret)
	}
}

func TestStakingCaller_Create(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(stakeTests)

	stakingAPICaller := NewMockStakingAPICaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingCaller := NewMockStakingCaller(ctrl)
	stakingCaller.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	clt := NewMockAuthedClient(ctrl)
	clt.EXPECT().Staking().Return(stakingCaller).Times(testTimes)

	for _, test := range stakeTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, nil).Times(1)

		ret, err := clt.Staking().
			Create(test.candidateName, test.amount, test.duration, test.autoStake).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.NoError(err)
		require.Equal(test.actionHash, ret)
	}
}

func TestStakingCaller_Unstake(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(unstakeTests)

	stakingAPICaller := NewMockStakingAPICaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingCaller := NewMockStakingCaller(ctrl)
	stakingCaller.EXPECT().Unstake(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	clt := NewMockAuthedClient(ctrl)
	clt.EXPECT().Staking().Return(stakingCaller).Times(testTimes)

	for _, test := range unstakeTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, nil).Times(1)

		ret, err := clt.Staking().
			Unstake(test.bucketIndex).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.NoError(err)
		require.Equal(test.actionHash, ret)
	}
}
