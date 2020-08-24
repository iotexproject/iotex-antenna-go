package main

import (
	"context"
	"fmt"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

func main() {
	s := NewIotexService("", "api.testnet.iotex.one:443", true)

	r, err := s.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	fmt.Println("chain meta", r, err)

	blockMetasRequest := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: 10000,
				Count: 1,
			},
		},
	}
	BlockMetasResponse, err := s.GetBlockMetas(context.Background(), blockMetasRequest)
	fmt.Println("block metas", BlockMetasResponse, err)

	getActionsRequest := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByIndex{
			ByIndex: &iotexapi.GetActionsByIndexRequest{
				Start: 1000000,
				Count: 1,
			},
		},
	}
	getActionsResponse, err := s.GetActions(context.Background(), getActionsRequest)
	fmt.Println("action", getActionsResponse, err)

	getCandidatesResponse, err := s.GetStakingCandidates(context.Background(), 7060000)
	fmt.Println("candidates", getCandidatesResponse, err)

	getBucketsResponse, err := s.GetStakingBuckets(context.Background(), 7060000)
	fmt.Println("buckets", getBucketsResponse, err)
}
