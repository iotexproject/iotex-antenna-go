// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

func main() {
	s := NewGetInfoService("", "api.testnet.iotex.one:443", true)

	r, err := s.GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	out, _ := json.MarshalIndent(r, "", "\t")
	fmt.Println("chain meta", string(out), err)

	blockMetasRequest := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: 10000,
				Count: 1,
			},
		},
	}
	BlockMetasResponse, err := s.GetBlockMetas(context.Background(), blockMetasRequest)
	out, _ = json.MarshalIndent(BlockMetasResponse, "", "\t")
	fmt.Println("block metas", string(out), err)

	getActionsRequest := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByIndex{
			ByIndex: &iotexapi.GetActionsByIndexRequest{
				Start: 1000000,
				Count: 1,
			},
		},
	}
	getActionsResponse, err := s.GetActions(context.Background(), getActionsRequest)
	out, _ = json.MarshalIndent(getActionsResponse, "", "\t")
	fmt.Println("action", string(out), err)

	getCandidatesResponse, err := s.GetStakingCandidates(context.Background(), 7060000)
	out, _ = json.MarshalIndent(getCandidatesResponse, "", "\t")
	fmt.Println("candidates", string(out), err)

	getBucketsResponse, err := s.GetStakingBuckets(context.Background(), 7060000)
	out, _ = json.MarshalIndent(getBucketsResponse, "", "\t")
	fmt.Println("buckets", string(out), err)
}
