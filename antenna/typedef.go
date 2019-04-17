// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package antenna

import (
	"github.com/iotexproject/iotex-antenna-go/iotx"
)

// request
type GetAccountRequest = iotx.GetAccountRequest

type GetActionsRequest = iotx.GetActionsRequest
type GetActionsRequest_ByHash = iotx.GetActionsRequest_ByHash
type GetActionByHashRequest = iotx.GetActionByHashRequest
type GetActionsRequest_ByIndex = iotx.GetActionsRequest_ByIndex
type GetActionsByIndexRequest = iotx.GetActionsByIndexRequest
type SendActionRequest = iotx.SendActionRequest
type GetActionsRequest_ByAddr = iotx.GetActionsRequest_ByAddr
type GetActionsByAddressRequest = iotx.GetActionsByAddressRequest
type GetActionsRequest_UnconfirmedByAddr = iotx.GetActionsRequest_UnconfirmedByAddr
type GetUnconfirmedActionsByAddressRequest = iotx.GetUnconfirmedActionsByAddressRequest
type GetActionsRequest_ByBlk = iotx.GetActionsRequest_ByBlk
type GetActionsByBlockRequest = iotx.GetActionsByBlockRequest
type GetBlockMetasRequest = iotx.GetBlockMetasRequest
type GetBlockMetasRequest_ByIndex = iotx.GetBlockMetasRequest_ByIndex
type GetBlockMetasByIndexRequest = iotx.GetBlockMetasByIndexRequest
type GetBlockMetasRequest_ByHash = iotx.GetBlockMetasRequest_ByHash
type GetBlockMetaByHashRequest = iotx.GetBlockMetaByHashRequest
type GetChainMetaRequest = iotx.GetChainMetaRequest
type GetServerMetaRequest = iotx.GetServerMetaRequest
type ReadStateRequest = iotx.ReadStateRequest
type GetReceiptByActionRequest = iotx.GetReceiptByActionRequest
type ReadContractRequest = iotx.ReadContractRequest
type SuggestGasPriceRequest = iotx.SuggestGasPriceRequest
type EstimateGasForActionRequest = iotx.EstimateGasForActionRequest
type GetEpochMetaRequest = iotx.GetEpochMetaRequest

// response
type GetAccountResponse = iotx.GetAccountResponse
type GetActionsResponse = iotx.GetActionsResponse
type GetBlockMetasResponse = iotx.GetBlockMetasResponse
type GetChainMetaResponse = iotx.GetChainMetaResponse
type GetServerMetaResponse = iotx.GetServerMetaResponse
type SendActionResponse = iotx.SendActionResponse
type GetReceiptByActionResponse = iotx.GetReceiptByActionResponse
type ReadContractResponse = iotx.ReadContractResponse
type SuggestGasPriceResponse = iotx.SuggestGasPriceResponse
type EstimateGasForActionResponse = iotx.EstimateGasForActionResponse
type GetEpochMetaResponse = iotx.GetEpochMetaResponse
type ReadStateResponse = iotx.ReadStateResponse
