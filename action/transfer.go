// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"

	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/pkg/errors"
)

// NewTransfer return new Transfer ActionCore
func NewTransfer(
	nonce uint64, gasLimit uint64, gasPrice *big.Int, amount *big.Int, recipient string, payload []byte,
) (*IotexActionCore, error) {
	if amount.Sign() == -1 || gasPrice.Sign() == -1 || recipient == "" {
		return nil, errors.New("invalid input for NewTransfer()")
	}
	return &IotexActionCore{
		ActionCore: &iotextypes.ActionCore{
			Version:  1,
			Nonce:    nonce,
			GasLimit: gasLimit,
			GasPrice: gasPrice.String(),
			Action: &iotextypes.ActionCore_Transfer{
				Transfer: &iotextypes.Transfer{
					Amount:    amount.String(),
					Recipient: recipient,
					Payload:   payload,
				},
			},
		},
	}, nil
}
