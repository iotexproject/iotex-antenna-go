package iotex

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

// SendActionCaller is used to perform a send action call.
type SendActionCaller interface {
	API() iotexapi.APIServiceClient
	Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error)
}

// TransferCaller is used to perform a transfer call.
type TransferCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) TransferCaller
	SetGasLimit(uint64) TransferCaller
	SetPayload([]byte) TransferCaller
	SetNonce(uint64) TransferCaller
}

// ClaimRewardCaller is used to perform a claim reward call.
type ClaimRewardCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) ClaimRewardCaller
	SetGasLimit(uint64) ClaimRewardCaller
	SetData([]byte) ClaimRewardCaller
	SetNonce(uint64) ClaimRewardCaller
}

// DeployContractCaller is used to perform a deploy contract call.
type DeployContractCaller interface {
	SendActionCaller

	SetArgs(abi abi.ABI, args ...interface{}) DeployContractCaller
	SetGasPrice(*big.Int) DeployContractCaller
	SetGasLimit(uint64) DeployContractCaller
	SetNonce(uint64) DeployContractCaller
}

// GetReceiptCaller is used to perform a get receipt call.
type GetReceiptCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetReceiptByActionResponse, error)
}

// GetLogsCaller is used to get logs
type GetLogsCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (*iotexapi.GetLogsResponse, error)
}

// AuthedClient is an iotex client which associate with an account credentials, so it can perform write actions.
type AuthedClient interface {
	ReadOnlyClient

	Contract(contract address.Address, abi abi.ABI) Contract
	Transfer(to address.Address, value *big.Int) TransferCaller
	ClaimReward(value *big.Int) ClaimRewardCaller
	DeployContract(data []byte) DeployContractCaller
	// staking related
	Staking() StakingCaller
	Candidate() CandidateCaller
	Account() account.Account
}

// ReadOnlyClient is an iotex client which can perform read actions.
type ReadOnlyClient interface {
	ReadOnlyContract(contract address.Address, abi abi.ABI) ReadOnlyContract
	GetReceipt(actionHash hash.Hash256) GetReceiptCaller
	GetLogs(request *iotexapi.GetLogsRequest) GetLogsCaller
	API() iotexapi.APIServiceClient
}

// ReadContractCaller is used to perform a read contract call.
type ReadContractCaller interface {
	Call(ctx context.Context, opts ...grpc.CallOption) (Data, error)
}

// ExecuteContractCaller is used to perform an execute contract call.
type ExecuteContractCaller interface {
	SendActionCaller

	SetGasPrice(*big.Int) ExecuteContractCaller
	SetGasLimit(uint64) ExecuteContractCaller
	SetAmount(*big.Int) ExecuteContractCaller
	SetNonce(uint64) ExecuteContractCaller
}

// Contract allows to read or execute on this contract's methods.
type Contract interface {
	ReadOnlyContract

	Execute(method string, args ...interface{}) ExecuteContractCaller
}

// ReadOnlyContract allows to read on this contract's methods.
type ReadOnlyContract interface {
	Read(method string, args ...interface{}) ReadContractCaller
}

// StakingCaller is used to perform a staking call.
type StakingCaller interface {
	Create(candidateName string, amount *big.Int, duration uint32, autoStake bool) StakingAPICaller
	Unstake(bucketIndex uint64) StakingAPICaller
	Withdraw(bucketIndex uint64) StakingAPICaller
	AddDeposit(index uint64, amount *big.Int) StakingAPICaller
	ChangeCandidate(candName string, bucketIndex uint64) StakingAPICaller
	StakingTransfer(voterAddress address.Address, bucketIndex uint64) StakingAPICaller
	Restake(index uint64, duration uint32, autoStake bool) StakingAPICaller
}

// CandidateCaller is used to perform a candidate call.
type CandidateCaller interface {
	Register(name string, ownerAddr, operatorAddr, rewardAddr address.Address, amount *big.Int, duration uint32, autoStake bool, payload []byte) StakingAPICaller
	Update(name string, operatorAddr, rewardAddr address.Address) StakingAPICaller
}

// StakingAPICaller is used to perform extra info call.
type StakingAPICaller interface {
	SendActionCaller
	SetGasPrice(*big.Int) StakingAPICaller
	SetGasLimit(uint64) StakingAPICaller
	SetNonce(uint64) StakingAPICaller
	SetPayload([]byte) StakingAPICaller
}
