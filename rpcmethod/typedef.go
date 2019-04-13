// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package rpcmethod

import (
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

type GetAccountResponse struct {
	*iotexapi.GetAccountResponse
}
type GetActionsResponse struct {
	*iotexapi.GetActionsResponse
}
type GetBlockMetasResponse struct {
	*iotexapi.GetBlockMetasResponse
}
type GetChainMetaResponse struct {
	*iotexapi.GetChainMetaResponse
}
type GetServerMetaResponse struct {
	*iotexapi.GetServerMetaResponse
}
type SendActionResponse struct {
	*iotexapi.SendActionResponse
}
type GetReceiptByActionResponse struct {
	*iotexapi.GetReceiptByActionResponse
}
type ReadContractResponse struct {
	*iotexapi.ReadContractResponse
}
type SuggestGasPriceResponse struct {
	*iotexapi.SuggestGasPriceResponse
}
type EstimateGasForActionResponse struct {
	*iotexapi.EstimateGasForActionResponse
}
type GetEpochMetaResponse struct {
	*iotexapi.GetEpochMetaResponse
}
type ReadStateResponse struct {
	*iotexapi.ReadStateResponse
}
