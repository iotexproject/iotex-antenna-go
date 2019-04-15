// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"context"
	"math/big"

	"github.com/iotexproject/iotex-core/protogen/iotextypes"

	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
	"github.com/iotexproject/iotex-core/testutil"
)

// RPCMethod provides simple interface tp invoke rpc method
type RPCMethod struct {
	Endpoint string
	conn     *grpc.ClientConn
	cli      iotexapi.APIServiceClient
}

// NewRPCMethod returns RPCMethod interacting with endpoint
func NewRPCMethod(endpoint string) (*RPCMethod, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cli := iotexapi.NewAPIServiceClient(conn)

	return &RPCMethod{
		Endpoint: endpoint,
		conn:     conn,
		cli:      cli,
	}, nil
}

// Close closes the underlaying connection, after Close, no method should be
// invoked anymore
func (r *RPCMethod) Close() {
	r.conn.Close()
}

// GetAccount gets the address detail of an address
func (r *RPCMethod) GetAccount(accountAddress string) (*GetAccountResponse, error) {
	in := &iotexapi.GetAccountRequest{Address: accountAddress}
	ctx := context.Background()
	ret, err := r.cli.GetAccount(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetAccountResponse{ret}
	return res, nil
}

// GetActions gets action(s) by:
// 1. start index and action count
// 2. action hash
// 3. address with start index and action count
// 4. get unconfirmed actions by address with start index and action count
// 5. block hash with start index and action count
func (r *RPCMethod) getActions(in *iotexapi.GetActionsRequest) (*GetActionsResponse, error) {
	ctx := context.Background()
	ret, err := r.cli.GetActions(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetActionsResponse{ret}
	return res, nil
}

// GetActionsByIndex
func (r *RPCMethod) GetActionsByIndex(start, count uint64) (*GetActionsResponse, error) {
	in := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByIndex{
			ByIndex: &iotexapi.GetActionsByIndexRequest{
				Start: start,
				Count: count,
			},
		},
	}
	return r.getActions(in)
}

// GetActionsByHash
func (r *RPCMethod) GetActionsByHash(hash string, checkPending bool) (*GetActionsResponse, error) {
	in := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash:   hash,
				CheckPending: checkPending,
			},
		},
	}
	return r.getActions(in)
}

// GetActionsByAddress
func (r *RPCMethod) GetActionsByAddress(address string, start, count uint64) (*GetActionsResponse, error) {
	in := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByAddr{
			ByAddr: &iotexapi.GetActionsByAddressRequest{
				Address: address,
				Start:   start,
				Count:   count,
			},
		},
	}
	return r.getActions(in)
}

// GetUnconfirmedActionsByAddress
func (r *RPCMethod) GetUnconfirmedActionsByAddress(address string, start, count uint64) (*GetActionsResponse, error) {
	in := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_UnconfirmedByAddr{
			UnconfirmedByAddr: &iotexapi.GetUnconfirmedActionsByAddressRequest{
				Address: address,
				Start:   start,
				Count:   count,
			},
		},
	}
	return r.getActions(in)
}

// GetActionsByBlock
func (r *RPCMethod) GetActionsByBlock(blockHash string, start, count uint64) (*GetActionsResponse, error) {
	in := &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByBlk{
			ByBlk: &iotexapi.GetActionsByBlockRequest{
				BlkHash: blockHash,
				Start:   start,
				Count:   count,
			},
		},
	}
	return r.getActions(in)
}

// GetBlockMetas gets block metadata(s) by:
// 1. start index and block count
// 2. block hash
// GetBlockMetasByIndexAndCount
func (r *RPCMethod) GetBlockMetasByIndexAndCount(index, count uint64) (*GetBlockMetasResponse, error) {
	in := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByIndex{
			ByIndex: &iotexapi.GetBlockMetasByIndexRequest{
				Start: index,
				Count: count,
			},
		},
	}
	return r.getBlockMetas(in)
}

// GetBlockMetasByBlockHash
func (r *RPCMethod) GetBlockMetasByBlockHash(blockHash string) (*GetBlockMetasResponse, error) {
	in := &iotexapi.GetBlockMetasRequest{
		Lookup: &iotexapi.GetBlockMetasRequest_ByHash{
			ByHash: &iotexapi.GetBlockMetaByHashRequest{
				BlkHash: blockHash,
			},
		},
	}
	return r.getBlockMetas(in)
}
func (r *RPCMethod) getBlockMetas(in *iotexapi.GetBlockMetasRequest) (*GetBlockMetasResponse, error) {
	ctx := context.Background()
	ret, err := r.cli.GetBlockMetas(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetBlockMetasResponse{ret}
	return res, nil
}

// GetChainMeta gets chain metadata
func (r *RPCMethod) GetChainMeta() (*GetChainMetaResponse, error) {
	ctx := context.Background()
	in := &iotexapi.GetChainMetaRequest{}
	ret, err := r.cli.GetChainMeta(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetChainMetaResponse{ret}
	return res, nil
}

// GetServerMeta gets server metadata
func (r *RPCMethod) GetServerMeta() (*GetServerMetaResponse, error) {
	ctx := context.Background()
	in := &iotexapi.GetServerMetaRequest{}
	ret, err := r.cli.GetServerMeta(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetServerMetaResponse{ret}
	return res, nil
}
func (r *RPCMethod) sendAction(act action.SealedEnvelope) (*SendActionResponse, error) {
	in := &iotexapi.SendActionRequest{Action: act.Proto()}
	ctx := context.Background()
	ret, err := r.cli.SendAction(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &SendActionResponse{ret}
	return res, nil
}

// SendTransfer sends transfer to svr
func (r *RPCMethod) SendTransfer(recipientAddr, senderPriKey string, nonce uint64, amount *big.Int, payload []byte, gasLimit uint64, gasPrice *big.Int) (*SendActionResponse, error) {
	priKey, err := keypair.HexStringToPrivateKey(senderPriKey)
	if err != nil {
		return nil, err
	}
	transfer, err := testutil.SignedTransfer(recipientAddr,
		priKey, nonce, amount, payload, gasLimit, gasPrice)
	if err != nil {
		return nil, err
	}
	return r.sendAction(transfer)
}

// SendVote sends vote to svr
func (r *RPCMethod) SendVote(voteeAddr, voterPriKey string, nonce uint64, gasLimit uint64, gasPrice *big.Int) (*SendActionResponse, error) {
	priKey, err := keypair.HexStringToPrivateKey(voterPriKey)
	if err != nil {
		return nil, err
	}
	vote, err := testutil.SignedVote(voteeAddr, priKey, nonce, gasLimit, gasPrice)
	if err != nil {
		return nil, err
	}
	return r.sendAction(vote)
}

// SendExecution sends Execution to svr
func (r *RPCMethod) SendExecution(contractAddr, executorPriKey string, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) (*SendActionResponse, error) {
	priKey, err := keypair.HexStringToPrivateKey(executorPriKey)
	if err != nil {
		return nil, err
	}
	transfer, err := testutil.SignedExecution(contractAddr, priKey, nonce, amount, gasLimit, gasPrice, data)
	if err != nil {
		return nil, err
	}
	return r.sendAction(transfer)
}

// GetReceiptByAction gets receipt by action hash
func (r *RPCMethod) GetReceiptByAction(hash string) (*GetReceiptByActionResponse, error) {
	ctx := context.Background()
	in := &iotexapi.GetReceiptByActionRequest{ActionHash: hash}
	ret, err := r.cli.GetReceiptByAction(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetReceiptByActionResponse{ret}
	return res, nil
}

// ReadContract reads contract
func (r *RPCMethod) ReadContract(action *iotextypes.Action, checkPending bool) (*ReadContractResponse, error) {
	ctx := context.Background()
	in := &iotexapi.ReadContractRequest{Action: action}
	ret, err := r.cli.ReadContract(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &ReadContractResponse{ret}
	return res, nil
}

// SuggestGasPrice suggests gas price
func (r *RPCMethod) SuggestGasPrice() (*SuggestGasPriceResponse, error) {
	ctx := context.Background()
	in := &iotexapi.SuggestGasPriceRequest{}
	ret, err := r.cli.SuggestGasPrice(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &SuggestGasPriceResponse{ret}
	return res, nil
}

// EstimateGasForAction estimates gas for action
func (r *RPCMethod) EstimateGasForAction(recipientAddr, executorPriKey string, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, payload []byte) (*EstimateGasForActionResponse, error) {
	ctx := context.Background()
	priKey, err := keypair.HexStringToPrivateKey(executorPriKey)
	if err != nil {
		return nil, err
	}
	transfer, err := testutil.SignedTransfer(recipientAddr,
		priKey, nonce, amount, payload, gasLimit, gasPrice)
	if err != nil {
		return nil, err
	}
	in := &iotexapi.EstimateGasForActionRequest{Action: transfer.Proto()}
	ret, err := r.cli.EstimateGasForAction(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &EstimateGasForActionResponse{ret}
	return res, nil
}

// ReadState reads state from blockchain
func (r *RPCMethod) ReadState(protoID, methodName string, arguments [][]byte) (*ReadStateResponse, error) {
	ctx := context.Background()
	in := &iotexapi.ReadStateRequest{
		ProtocolID: []byte(protoID),
		MethodName: []byte(methodName),
		Arguments:  arguments,
	}
	ret, err := r.cli.ReadState(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &ReadStateResponse{ret}
	return res, nil
}

// GetEpochMeta get epoch meta
func (r *RPCMethod) GetEpochMeta(epochNumber uint64) (*GetEpochMetaResponse, error) {
	ctx := context.Background()
	in := &iotexapi.GetEpochMetaRequest{EpochNumber: epochNumber}
	ret, err := r.cli.GetEpochMeta(ctx, in)
	if err != nil {
		return nil, err
	}
	res := &GetEpochMetaResponse{ret}
	return res, nil
}
