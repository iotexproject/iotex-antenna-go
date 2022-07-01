// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"context"
	"math/big"

	"google.golang.org/grpc"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-antenna-go/v2/errcodes"
)

type (
	stakingCaller struct {
		*sendActionCaller
		action interface{}
	}

	// reclaim to differentiate unstake and withdraw
	reclaim struct {
		action     *iotextypes.StakeReclaim
		isWithdraw bool
	}
)

//Create Staking
func (c *stakingCaller) Create(candidateName string, amount *big.Int, duration uint32, autoStake bool) SendActionCaller {
	tx := iotextypes.StakeCreate{
		CandidateName:  candidateName,
		StakedDuration: duration,
		AutoStake:      autoStake,
		StakedAmount:   amount.String(),
	}
	c.action = &tx
	return c
}

//Unstake Staking
func (c *stakingCaller) Unstake(bucketIndex uint64) SendActionCaller {
	tx := iotextypes.StakeReclaim{
		BucketIndex: bucketIndex,
	}
	c.action = &reclaim{&tx, false}
	return c
}

//Withdraw Staking
func (c *stakingCaller) Withdraw(bucketIndex uint64) SendActionCaller {
	tx := iotextypes.StakeReclaim{
		BucketIndex: bucketIndex,
	}
	c.action = &reclaim{&tx, true}
	return c
}

//AddDeposit Staking
func (c *stakingCaller) AddDeposit(index uint64, amount *big.Int) SendActionCaller {
	tx := iotextypes.StakeAddDeposit{
		BucketIndex: index,
		Amount:      amount.String(),
	}
	c.action = &tx
	return c
}

//ChangeCandidate Staking
func (c *stakingCaller) ChangeCandidate(candName string, bucketIndex uint64) SendActionCaller {
	tx := iotextypes.StakeChangeCandidate{
		CandidateName: candName,
		BucketIndex:   bucketIndex,
	}
	c.action = &tx
	return c
}

//StakingTransfer Staking
func (c *stakingCaller) StakingTransfer(voterAddress address.Address, bucketIndex uint64) SendActionCaller {
	tx := iotextypes.StakeTransferOwnership{
		VoterAddress: voterAddress.String(),
		BucketIndex:  bucketIndex,
	}
	c.action = &tx
	return c
}

//Restake Staking
func (c *stakingCaller) Restake(index uint64, duration uint32, autoStake bool) SendActionCaller {
	tx := iotextypes.StakeRestake{
		BucketIndex:    index,
		StakedDuration: duration,
		AutoStake:      autoStake,
	}
	c.action = &tx
	return c
}

//Register Staking
func (c *stakingCaller) Register(name string, ownerAddr, operatorAddr, rewardAddr address.Address, amount *big.Int, duration uint32, autoStake bool, payload []byte) SendActionCaller {
	basic := iotextypes.CandidateBasicInfo{
		Name:            name,
		OperatorAddress: operatorAddr.String(),
		RewardAddress:   rewardAddr.String(),
	}
	tx := iotextypes.CandidateRegister{
		Candidate:      &basic,
		StakedAmount:   amount.String(),
		StakedDuration: duration,
		AutoStake:      autoStake,
		OwnerAddress:   ownerAddr.String(),
		Payload:        payload,
	}
	c.action = &tx
	return c
}

//Update Staking
func (c *stakingCaller) Update(name string, operatorAddr, rewardAddr address.Address) SendActionCaller {
	tx := iotextypes.CandidateBasicInfo{
		Name:            name,
		OperatorAddress: operatorAddr.String(),
		RewardAddress:   rewardAddr.String(),
	}
	c.action = &tx
	return c
}

func (c *stakingCaller) SetGasLimit(g uint64) SendActionCaller {
	c.sendActionCaller.setGasLimit(g)
	return c
}

func (c *stakingCaller) SetGasPrice(g *big.Int) SendActionCaller {
	c.sendActionCaller.setGasPrice(g)
	return c
}

func (c *stakingCaller) SetNonce(n uint64) SendActionCaller {
	c.sendActionCaller.setNonce(n)
	return c
}

func (c *stakingCaller) SetPayload(pl []byte) SendActionCaller {
	c.sendActionCaller.setPayload(pl)
	return c
}

//Call call sendActionCaller
func (c *stakingCaller) Call(ctx context.Context, opts ...grpc.CallOption) (hash.Hash256, error) {
	c.core = &iotextypes.ActionCore{
		Version: ProtocolVersion,
	}

	hasPayload := len(c.payload) > 0
	switch a := c.action.(type) {
	case *iotextypes.StakeCreate:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_StakeCreate{StakeCreate: a}
	case *reclaim:
		if hasPayload {
			a.action.Payload = c.payload
		}
		if a.isWithdraw {
			c.core.Action = &iotextypes.ActionCore_StakeWithdraw{StakeWithdraw: a.action}
		} else {
			c.core.Action = &iotextypes.ActionCore_StakeUnstake{StakeUnstake: a.action}
		}
	case *iotextypes.StakeAddDeposit:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_StakeAddDeposit{StakeAddDeposit: a}
	case *iotextypes.StakeRestake:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_StakeRestake{StakeRestake: a}
	case *iotextypes.StakeChangeCandidate:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_StakeChangeCandidate{StakeChangeCandidate: a}
	case *iotextypes.StakeTransferOwnership:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_StakeTransferOwnership{StakeTransferOwnership: a}
	case *iotextypes.CandidateRegister:
		if hasPayload {
			a.Payload = c.payload
		}
		c.core.Action = &iotextypes.ActionCore_CandidateRegister{CandidateRegister: a}
	case *iotextypes.CandidateBasicInfo:
		c.core.Action = &iotextypes.ActionCore_CandidateUpdate{CandidateUpdate: a}
	default:
		return hash.ZeroHash256, errcodes.New("not support action call", errcodes.InternalError)
	}
	return c.sendActionCaller.Call(ctx, opts...)
}
