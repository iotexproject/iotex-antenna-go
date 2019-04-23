// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpc

import (
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

type (
	// request
	GetAccountRequest                     = iotexapi.GetAccountRequest
	GetActionsRequest                     = iotexapi.GetActionsRequest
	GetActionsRequest_ByHash              = iotexapi.GetActionsRequest_ByHash
	GetActionByHashRequest                = iotexapi.GetActionByHashRequest
	GetActionsRequest_ByIndex             = iotexapi.GetActionsRequest_ByIndex
	GetActionsByIndexRequest              = iotexapi.GetActionsByIndexRequest
	SendActionRequest                     = iotexapi.SendActionRequest
	GetActionsRequest_ByAddr              = iotexapi.GetActionsRequest_ByAddr
	GetActionsByAddressRequest            = iotexapi.GetActionsByAddressRequest
	GetActionsRequest_UnconfirmedByAddr   = iotexapi.GetActionsRequest_UnconfirmedByAddr
	GetUnconfirmedActionsByAddressRequest = iotexapi.GetUnconfirmedActionsByAddressRequest
	GetActionsRequest_ByBlk               = iotexapi.GetActionsRequest_ByBlk
	GetActionsByBlockRequest              = iotexapi.GetActionsByBlockRequest
	GetBlockMetasRequest                  = iotexapi.GetBlockMetasRequest
	GetBlockMetasRequest_ByIndex          = iotexapi.GetBlockMetasRequest_ByIndex
	GetBlockMetasByIndexRequest           = iotexapi.GetBlockMetasByIndexRequest
	GetBlockMetasRequest_ByHash           = iotexapi.GetBlockMetasRequest_ByHash
	GetBlockMetaByHashRequest             = iotexapi.GetBlockMetaByHashRequest
	GetChainMetaRequest                   = iotexapi.GetChainMetaRequest
	GetServerMetaRequest                  = iotexapi.GetServerMetaRequest
	ReadStateRequest                      = iotexapi.ReadStateRequest
	GetReceiptByActionRequest             = iotexapi.GetReceiptByActionRequest
	ReadContractRequest                   = iotexapi.ReadContractRequest
	SuggestGasPriceRequest                = iotexapi.SuggestGasPriceRequest
	EstimateGasForActionRequest           = iotexapi.EstimateGasForActionRequest
	GetEpochMetaRequest                   = iotexapi.GetEpochMetaRequest

	// response
	GetAccountResponse           = iotexapi.GetAccountResponse
	GetActionsResponse           = iotexapi.GetActionsResponse
	GetBlockMetasResponse        = iotexapi.GetBlockMetasResponse
	GetChainMetaResponse         = iotexapi.GetChainMetaResponse
	GetServerMetaResponse        = iotexapi.GetServerMetaResponse
	SendActionResponse           = iotexapi.SendActionResponse
	GetReceiptByActionResponse   = iotexapi.GetReceiptByActionResponse
	ReadContractResponse         = iotexapi.ReadContractResponse
	SuggestGasPriceResponse      = iotexapi.SuggestGasPriceResponse
	EstimateGasForActionResponse = iotexapi.EstimateGasForActionResponse
	GetEpochMetaResponse         = iotexapi.GetEpochMetaResponse
	ReadStateResponse            = iotexapi.ReadStateResponse
)
