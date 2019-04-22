// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"math/big"

	"github.com/iotexproject/iotex-core/action"
)

// TransferRequest defines transfer request parameters
type TransferRequest struct {
	From     string
	To       string
	Value    string
	Payload  string
	GasLimit string
	GasPrice string
}

// ContractRequest defines contract request parameters
type ContractRequest struct {
	From   string
	Amount string
	// contract bytecode
	Data     string
	Abi      string
	GasLimit string
	GasPrice string
}

// NewTransferEnvelop return action envelop
func NewTransferEnvelop(
	nonce uint64,
	amount *big.Int,
	recipient string,
	payload string,
	gasLimit uint64,
	gasPrice *big.Int) (action.Envelope, error) {
	tx, err := action.NewTransfer(nonce, amount,
		recipient, []byte(payload), gasLimit, gasPrice)
	if err != nil {
		return action.Envelope{}, err
	}
	bd := &action.EnvelopeBuilder{}
	return bd.SetNonce(nonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(tx).Build(), nil
}
