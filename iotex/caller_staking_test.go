package iotex

import (
	"context"
	"errors"
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
		owner     string
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
		err        error
	}{
		{
			"1111",
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
			nil,
		},
		// failure case
		{
			"2222",
			"00000000000000000000000000000000000000000",
			"00000000000000000000000000000000000000000",
			"00000000000000000000000000000000000000000",
			big.NewInt(0),
			0,
			false,
			big.NewInt(int64(0)),
			0,
			[]byte(""),
			hash.Hash256{},
			errors.New("address error"),
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
		err        error
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
			nil,
		},
		// failure case
		{
			"00000000000000000000000000000000000000000",
			big.NewInt(int64(unit.Iotx)),
			10000,
			false,
			big.NewInt(int64(10 * unit.Qev)),
			1000000,
			[]byte("TestStaking"),
			hash.Hash256{},
			errors.New("address error"),
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
		err        error
	}{
		{
			10000,
			false,
			big.NewInt(int64(10 * unit.Qev)),
			1000000,
			[]byte("TestStaking"),
			testActionHash1,
			nil,
		},
		// failure case
		{
			0,
			false,
			big.NewInt(0),
			0,
			[]byte{},
			hash.Hash256{},
			errors.New("invalid gas price"),
		},
	}
)

func TestCandidateCaller_Register(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(candidateRegisterTests)

	stakingAPICaller := NewMockSendActionCaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	candidateCaller := NewMockCandidateCaller(ctrl)
	candidateCaller.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	client := NewMockAuthedClient(ctrl)
	client.EXPECT().Candidate().Return(candidateCaller).Times(testTimes)

	for _, test := range candidateRegisterTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, test.err).Times(1)
		ownerAddr, _ := address.FromString(test.owner)
		operatorAddr, _ := address.FromString(test.operator)
		rewardAddr, _ := address.FromString(test.reward)
		ret, err := client.Candidate().
			Register(test.name, ownerAddr, operatorAddr, rewardAddr, test.amount, test.duration, test.autoStake, test.payload).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.Equal(test.err, err)
		require.Equal(test.actionHash, ret)
	}
}

func TestStakingCaller_Create(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(stakeTests)

	stakingAPICaller := NewMockSendActionCaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingCaller := NewMockStakingCaller(ctrl)
	stakingCaller.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	clt := NewMockAuthedClient(ctrl)
	clt.EXPECT().Staking().Return(stakingCaller).Times(testTimes)

	for _, test := range stakeTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, test.err).Times(1)
		ret, err := clt.Staking().
			Create(test.candidateName, test.amount, test.duration, test.autoStake).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.Equal(test.err, err)
		require.Equal(test.actionHash, ret)
	}
}

func TestStakingCaller_Unstake(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testTimes := len(unstakeTests)

	stakingAPICaller := NewMockSendActionCaller(ctrl)
	stakingAPICaller.EXPECT().SetGasPrice(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetGasLimit(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingAPICaller.EXPECT().SetPayload(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	stakingCaller := NewMockStakingCaller(ctrl)
	stakingCaller.EXPECT().Unstake(gomock.Any()).Return(stakingAPICaller).Times(testTimes)
	clt := NewMockAuthedClient(ctrl)
	clt.EXPECT().Staking().Return(stakingCaller).Times(testTimes)

	for _, test := range unstakeTests {
		stakingAPICaller.EXPECT().Call(gomock.Any()).Return(test.actionHash, test.err).Times(1)
		ret, err := clt.Staking().
			Unstake(test.bucketIndex).
			SetGasPrice(test.gasPrice).
			SetGasLimit(test.gasLimit).
			SetPayload(test.payload).
			Call(context.Background())
		require.Equal(test.err, err)
		require.Equal(test.actionHash, ret)
	}
}
