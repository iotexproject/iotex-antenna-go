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

// Caller is used to perform a send action call.
type Caller interface {
	API() iotexapi.APIServiceClient
	Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error)
}

// SendActionCaller is used to set nonce/gas etc. on top of Caller
type SendActionCaller interface {
	Caller
	SetNonce(uint64) SendActionCaller
	SetGasLimit(uint64) SendActionCaller
	SetGasPrice(*big.Int) SendActionCaller
	SetPayload([]byte) SendActionCaller
}

// ClaimRewardCaller is used to perform a claim reward call.
type ClaimRewardCaller interface {
	Caller

	SetGasPrice(*big.Int) ClaimRewardCaller
	SetGasLimit(uint64) ClaimRewardCaller
	SetData([]byte) ClaimRewardCaller
	SetNonce(uint64) ClaimRewardCaller
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
	Transfer(to address.Address, value *big.Int) SendActionCaller
	ClaimReward(value *big.Int) ClaimRewardCaller
	DeployContract(data []byte) DeployContractCaller
	// staking related
	Staking() StakingCaller
	Candidate() CandidateCaller
	Account() account.Account
	ChainID() uint32
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
	Caller

	SetGasPrice(*big.Int) ExecuteContractCaller
	SetGasLimit(uint64) ExecuteContractCaller
	SetAmount(*big.Int) ExecuteContractCaller
	SetNonce(uint64) ExecuteContractCaller
}

// DeployContractCaller is used to perform a deploy contract call.
type DeployContractCaller interface {
	Caller

	SetArgs(abi abi.ABI, args ...interface{}) DeployContractCaller
	SetGasPrice(*big.Int) DeployContractCaller
	SetGasLimit(uint64) DeployContractCaller
	SetNonce(uint64) DeployContractCaller
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
	Create(candidateName string, amount *big.Int, duration uint32, autoStake bool) SendActionCaller
	Unstake(bucketIndex uint64) SendActionCaller
	Withdraw(bucketIndex uint64) SendActionCaller
	AddDeposit(index uint64, amount *big.Int) SendActionCaller
	ChangeCandidate(candName string, bucketIndex uint64) SendActionCaller
	StakingTransfer(voterAddress address.Address, bucketIndex uint64) SendActionCaller
	Restake(index uint64, duration uint32, autoStake bool) SendActionCaller
}

// CandidateCaller is used to perform a candidate call.
type CandidateCaller interface {
	Register(name string, ownerAddr, operatorAddr, rewardAddr address.Address, amount *big.Int, duration uint32, autoStake bool, payload []byte) SendActionCaller
	Update(name string, operatorAddr, rewardAddr address.Address) SendActionCaller
}
