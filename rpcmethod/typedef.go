// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

// request
type GetAccountRequest = iotexapi.GetAccountRequest

type GetActionsRequest = iotexapi.GetActionsRequest
type GetActionsRequest_ByHash = iotexapi.GetActionsRequest_ByHash
type GetActionByHashRequest = iotexapi.GetActionByHashRequest
type GetActionsRequest_ByIndex = iotexapi.GetActionsRequest_ByIndex
type GetActionsByIndexRequest = iotexapi.GetActionsByIndexRequest
type SendActionRequest = iotexapi.SendActionRequest
type GetActionsRequest_ByAddr = iotexapi.GetActionsRequest_ByAddr
type GetActionsByAddressRequest = iotexapi.GetActionsByAddressRequest
type GetActionsRequest_UnconfirmedByAddr = iotexapi.GetActionsRequest_UnconfirmedByAddr
type GetUnconfirmedActionsByAddressRequest = iotexapi.GetUnconfirmedActionsByAddressRequest
type GetActionsRequest_ByBlk = iotexapi.GetActionsRequest_ByBlk
type GetActionsByBlockRequest = iotexapi.GetActionsByBlockRequest
type GetBlockMetasRequest = iotexapi.GetBlockMetasRequest
type GetBlockMetasRequest_ByIndex = iotexapi.GetBlockMetasRequest_ByIndex
type GetBlockMetasByIndexRequest = iotexapi.GetBlockMetasByIndexRequest
type GetBlockMetasRequest_ByHash = iotexapi.GetBlockMetasRequest_ByHash
type GetBlockMetaByHashRequest = iotexapi.GetBlockMetaByHashRequest
type GetChainMetaRequest = iotexapi.GetChainMetaRequest
type GetServerMetaRequest = iotexapi.GetServerMetaRequest
type ReadStateRequest = iotexapi.ReadStateRequest
type GetReceiptByActionRequest = iotexapi.GetReceiptByActionRequest
type ReadContractRequest = iotexapi.ReadContractRequest
type SuggestGasPriceRequest = iotexapi.SuggestGasPriceRequest
type EstimateGasForActionRequest = iotexapi.EstimateGasForActionRequest
type GetEpochMetaRequest = iotexapi.GetEpochMetaRequest

// response
type GetAccountResponse = iotexapi.GetAccountResponse
type GetActionsResponse = iotexapi.GetActionsResponse
type GetBlockMetasResponse = iotexapi.GetBlockMetasResponse
type GetChainMetaResponse = iotexapi.GetChainMetaResponse
type GetServerMetaResponse = iotexapi.GetServerMetaResponse
type SendActionResponse = iotexapi.SendActionResponse
type GetReceiptByActionResponse = iotexapi.GetReceiptByActionResponse
type ReadContractResponse = iotexapi.ReadContractResponse
type SuggestGasPriceResponse = iotexapi.SuggestGasPriceResponse
type EstimateGasForActionResponse = iotexapi.EstimateGasForActionResponse
type GetEpochMetaResponse = iotexapi.GetEpochMetaResponse
type ReadStateResponse = iotexapi.ReadStateResponse
