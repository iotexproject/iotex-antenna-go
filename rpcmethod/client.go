// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"context"

	"google.golang.org/grpc"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
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
func (r *RPCMethod) GetAccount(in *iotexapi.GetAccountRequest) (*iotexapi.GetAccountResponse, error) {
	ctx := context.Background()
	return r.cli.GetAccount(ctx, in)
}

// GetActions gets action(s) by:
// 1. start index and action count
// 2. action hash
// 3. address with start index and action count
// 4. get unconfirmed actions by address with start index and action count
// 5. block hash with start index and action count
func (r *RPCMethod) GetActions(in *iotexapi.GetActionsRequest) (*iotexapi.GetActionsResponse, error) {
	ctx := context.Background()
	return r.cli.GetActions(ctx, in)
}

// GetBlockMetas gets block metadata(s) by:
// 1. start index and block count
// 2. block hash
func (r *RPCMethod) GetBlockMetas(in *iotexapi.GetBlockMetasRequest) (*iotexapi.GetBlockMetasResponse, error) {
	ctx := context.Background()
	return r.cli.GetBlockMetas(ctx, in)
}

// GetChainMeta gets chain metadata
func (r *RPCMethod) GetChainMeta(in *iotexapi.GetChainMetaRequest) (*iotexapi.GetChainMetaResponse, error) {
	ctx := context.Background()
	return r.cli.GetChainMeta(ctx, in)
}

// GetServerMeta gets server metadata
func (r *RPCMethod) GetServerMeta(in *iotexapi.GetServerMetaRequest) (*iotexapi.GetServerMetaResponse, error) {
	ctx := context.Background()
	return r.cli.GetServerMeta(ctx, in)
}

// SendAction sends atcion to svr
func (r *RPCMethod) SendAction(in *iotexapi.SendActionRequest) (*iotexapi.SendActionResponse, error) {
	ctx := context.Background()
	return r.cli.SendAction(ctx, in)
}

// GetReceiptByAction gets receipt by action hash
func (r *RPCMethod) GetReceiptByAction(in *iotexapi.GetReceiptByActionRequest) (*iotexapi.GetReceiptByActionResponse, error) {
	ctx := context.Background()
	return r.cli.GetReceiptByAction(ctx, in)
}

// ReadContract reads contract
func (r *RPCMethod) ReadContract(in *iotexapi.ReadContractRequest) (*iotexapi.ReadContractResponse, error) {
	ctx := context.Background()
	return r.cli.ReadContract(ctx, in)
}

// SuggestGasPrice suggests gas price
func (r *RPCMethod) SuggestGasPrice(in *iotexapi.SuggestGasPriceRequest) (*iotexapi.SuggestGasPriceResponse, error) {
	ctx := context.Background()
	return r.cli.SuggestGasPrice(ctx, in)
}

// EstimateGasForAction estimates gas for action
func (r *RPCMethod) EstimateGasForAction(in *iotexapi.EstimateGasForActionRequest) (*iotexapi.EstimateGasForActionResponse, error) {
	ctx := context.Background()
	return r.cli.EstimateGasForAction(ctx, in)
}

// ReadState reads state from blockchain
func (r *RPCMethod) ReadState(in *iotexapi.ReadStateRequest) (*iotexapi.ReadStateResponse, error) {
	ctx := context.Background()
	return r.cli.ReadState(ctx, in)
}

// GetEpochMeta get epoch meta
func (r *RPCMethod) GetEpochMeta(in *iotexapi.GetEpochMetaRequest) (*iotexapi.GetEpochMetaResponse, error) {
	ctx := context.Background()
	return r.cli.GetEpochMeta(ctx, in)
}
