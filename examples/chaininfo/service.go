// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-antenna-go/v2/examples/util"
)

const (
	protocolID          = "staking"
	readBucketsLimit    = 30000
	readCandidatesLimit = 20000
)

// GetInfoService is the GetInfoService interface
type GetInfoService interface {
	// GetChainMeta is the GetChainMeta interface
	GetChainMeta(ctx context.Context, in *iotexapi.GetChainMetaRequest) (*iotexapi.GetChainMetaResponse, error)
	// GetBlockMetas is the GetBlockMetas interface
	GetBlockMetas(ctx context.Context, in *iotexapi.GetBlockMetasRequest) (*iotexapi.GetBlockMetasResponse, error)
	// GetActions is the GetActions interface
	GetActions(ctx context.Context, in *iotexapi.GetActionsRequest) (*iotexapi.GetActionsResponse, error)
	// GetStakingBuckets is the GetStakingBuckets interface
	GetStakingBuckets(ctx context.Context, height uint64) (*iotextypes.VoteBucketList, error)
	// GetStakingCandidates is the GetStakingCandidates interface
	GetStakingCandidates(ctx context.Context, height uint64) (*iotextypes.CandidateListV2, error)
}

type getInfoService struct {
	util.IotexService
}

// NewGetInfoService returns GetInfoService
func NewGetInfoService(accountPrivate, endpoint string, secure bool) GetInfoService {
	return &getInfoService{
		util.NewIotexService(accountPrivate, endpoint, secure),
	}
}

// GetChainMeta is the GetChainMeta interface
func (s *getInfoService) GetChainMeta(ctx context.Context, in *iotexapi.GetChainMetaRequest) (*iotexapi.GetChainMetaResponse, error) {
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	return s.ReadOnlyClient().API().GetChainMeta(ctx, in)
}

// GetBlockMetas is the GetBlockMetas interface
func (s *getInfoService) GetBlockMetas(ctx context.Context, in *iotexapi.GetBlockMetasRequest) (*iotexapi.GetBlockMetasResponse, error) {
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	return s.ReadOnlyClient().API().GetBlockMetas(ctx, in)
}

// GetActions is the GetActions interface
func (s *getInfoService) GetActions(ctx context.Context, in *iotexapi.GetActionsRequest) (*iotexapi.GetActionsResponse, error) {
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	return s.ReadOnlyClient().API().GetActions(ctx, in)
}

// GetStakingBuckets is the GetStakingBuckets interface
func (s *getInfoService) GetStakingBuckets(ctx context.Context, height uint64) (voteBucketListAll *iotextypes.VoteBucketList, err error) {
	err = s.Connect()
	if err != nil {
		return nil, err
	}
	voteBucketListAll = &iotextypes.VoteBucketList{}
	for i := uint32(0); ; i++ {
		offset := i * readBucketsLimit
		size := uint32(readBucketsLimit)
		voteBucketList, err := s.getStakingBuckets(ctx, offset, size, height)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get bucket")
		}
		voteBucketListAll.Buckets = append(voteBucketListAll.Buckets, voteBucketList.Buckets...)
		if len(voteBucketList.Buckets) < readBucketsLimit {
			break
		}
	}
	return
}

// GetStakingCandidates is the GetStakingCandidates interface
func (s *getInfoService) GetStakingCandidates(ctx context.Context, height uint64) (candidateListAll *iotextypes.CandidateListV2, err error) {
	err = s.Connect()
	if err != nil {
		return nil, err
	}
	candidateListAll = &iotextypes.CandidateListV2{}
	for i := uint32(0); ; i++ {
		offset := i * readCandidatesLimit
		size := uint32(readCandidatesLimit)
		candidateList, err := s.getStakingCandidates(ctx, offset, size, height)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get candidates")
		}
		candidateListAll.Candidates = append(candidateListAll.Candidates, candidateList.Candidates...)
		if len(candidateList.Candidates) < readCandidatesLimit {
			break
		}
	}
	return
}

func (s *getInfoService) getStakingBuckets(ctx context.Context, offset, limit uint32, height uint64) (voteBucketList *iotextypes.VoteBucketList, err error) {
	methodName, err := proto.Marshal(&iotexapi.ReadStakingDataMethod{
		Method: iotexapi.ReadStakingDataMethod_BUCKETS,
	})
	if err != nil {
		return nil, err
	}
	arg, err := proto.Marshal(&iotexapi.ReadStakingDataRequest{
		Request: &iotexapi.ReadStakingDataRequest_Buckets{
			Buckets: &iotexapi.ReadStakingDataRequest_VoteBuckets{
				Pagination: &iotexapi.PaginationParam{
					Offset: offset,
					Limit:  limit,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	readStateRequest := &iotexapi.ReadStateRequest{
		ProtocolID: []byte(protocolID),
		MethodName: methodName,
		Arguments:  [][]byte{arg},
		Height:     fmt.Sprintf("%d", height),
	}
	readStateRes, err := s.ReadOnlyClient().API().ReadState(ctx, readStateRequest)
	if err != nil {
		return
	}
	voteBucketList = &iotextypes.VoteBucketList{}
	if err := proto.Unmarshal(readStateRes.GetData(), voteBucketList); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal VoteBucketList")
	}
	return
}

func (s *getInfoService) getStakingCandidates(ctx context.Context, offset, limit uint32, height uint64) (candidateList *iotextypes.CandidateListV2, err error) {
	methodName, err := proto.Marshal(&iotexapi.ReadStakingDataMethod{
		Method: iotexapi.ReadStakingDataMethod_CANDIDATES,
	})
	if err != nil {
		return nil, err
	}
	arg, err := proto.Marshal(&iotexapi.ReadStakingDataRequest{
		Request: &iotexapi.ReadStakingDataRequest_Candidates_{
			Candidates: &iotexapi.ReadStakingDataRequest_Candidates{
				Pagination: &iotexapi.PaginationParam{
					Offset: offset,
					Limit:  limit,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	readStateRequest := &iotexapi.ReadStateRequest{
		ProtocolID: []byte(protocolID),
		MethodName: methodName,
		Arguments:  [][]byte{arg},
		Height:     fmt.Sprintf("%d", height),
	}
	readStateRes, err := s.ReadOnlyClient().API().ReadState(ctx, readStateRequest)
	if err != nil {
		return
	}
	candidateList = &iotextypes.CandidateListV2{}
	if err := proto.Unmarshal(readStateRes.GetData(), candidateList); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal VoteBucketList")
	}
	return
}
