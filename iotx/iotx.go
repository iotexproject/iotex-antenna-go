// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

// Iotx
type Iotx struct {
	*rpcmethod.RPCMethod
	Accounts account.Accounts
}

// New new Iotx
func New(host string) (*Iotx, error) {
	rpc, err := rpcmethod.NewRPCMethod(host)
	if err != nil {
		return nil, err
	}
	iotx := &Iotx{rpc, account.Accounts{}}
	return iotx, nil
}

// SendTransfer ...
func (i *Iotx) SendTransfer(request *TransferRequest) (string, error) {
	sender, exist := i.Accounts.GetAccount(request.From)
	if !exist {
		return "", fmt.Errorf("account:%s not exist", request.From)
	}

	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: request.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.Nonce
	amount, ok := new(big.Int).SetString(request.Value, 10)
	if !ok {
		return "", fmt.Errorf("amount:%s error", request.Value)
	}
	gasLimit, err := strconv.ParseUint(request.GasLimit, 10, 64)
	gasPrice, ok := new(big.Int).SetString(request.GasPrice, 10)
	if !ok {
		return "", fmt.Errorf("gasPrice:%s error", request.GasPrice)
	}

	tx, err := action.NewTransfer(nonce, amount,
		request.To, []byte(request.Payload), gasLimit, gasPrice)
	if err != nil {
		return "", err
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(tx).Build()

	return i.sendAction(sender, elp)
}

// DeployContract deploy contract
func (i *Iotx) DeployContract(request *ContractRequest) error {
	// TODO
	return nil
}

func (i *Iotx) sendAction(acc *account.Account, elp action.Envelope) (string, error) {
	sealed, err := action.Sign(elp, *acc.Private())
	if err != nil {
		return "", err
	}
	selp := sealed.Proto()
	request := &iotexapi.SendActionRequest{Action: selp}
	_, err = i.SendAction(request)
	shash := hash.Hash256b(byteutil.Must(proto.Marshal(selp)))
	return hex.EncodeToString(shash[:]), nil
}
