// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/contract"
	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
	"github.com/iotexproject/iotex-antenna-go/utils"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

// Iotx ...
type Iotx struct {
	*rpcmethod.RPCMethod
	Accounts *account.Accounts
}

// New new Iotx
func New(host string) (*Iotx, error) {
	rpc, err := rpcmethod.NewRPCMethod(host)
	if err != nil {
		return nil, err
	}
	iotx := &Iotx{rpc, account.NewAccounts()}
	return iotx, nil
}

func (i *Iotx) normalizeGas(acc *account.Account, elp action.Envelope, gasLimit, gasPrice string) (uint64, *big.Int, error) {
	var limit uint64
	var price *big.Int
	if gasLimit == "" {
		sealed, err := action.Sign(elp, *acc.Private())
		if err != nil {
			return 0, nil, err
		}
		selp := sealed.Proto()
		request := &iotexapi.EstimateGasForActionRequest{Action: selp}
		response, err := i.EstimateGasForAction(request)
		if err != nil {
			return 0, nil, err
		}
		limit = response.Gas
	} else {
		ul, err := strconv.ParseUint(gasLimit, 10, 64)
		if err != nil {
			return 0, nil, err
		}
		limit = ul
	}
	if gasPrice == "" {
		response, err := i.SuggestGasPrice(&iotexapi.SuggestGasPriceRequest{})
		if err != nil {
			return 0, nil, err
		}
		price = big.NewInt(0).SetUint64(response.GasPrice)
	} else {
		p, ok := big.NewInt(0).SetString(utils.ToRau(gasPrice, "Qev"), 10)
		if !ok {
			return 0, nil, fmt.Errorf("gas price %s error", gasPrice)
		}
		price = p
	}

	return limit, price, nil
}

// SendTransfer ...
func (i *Iotx) SendTransfer(req *TransferRequest) (string, error) {
	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", fmt.Errorf("account:%s not exist", req.From)
	}

	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	amount, ok := big.NewInt(0).SetString(req.Value, 10)
	if !ok {
		return "", fmt.Errorf("amount:%s error", req.Value)
	}

	tx, err := action.NewTransfer(nonce, amount,
		req.To, []byte(req.Payload), 0, big.NewInt(0))
	if err != nil {
		return "", err
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(big.NewInt(0)).
		SetGasLimit(0).
		SetAction(tx).Build()

	gasLimit, gasPrice, err := i.normalizeGas(sender, elp, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	tx, _ = action.NewTransfer(nonce, amount,
		req.To, []byte(req.Payload), gasLimit, gasPrice)
	elp = bd.SetNonce(nonce).
		SetGasPrice(gasPrice).
		SetGasLimit(gasLimit).
		SetAction(tx).Build()

	return i.sendAction(sender, elp)
}

// DeployContract ...
func (i *Iotx) DeployContract(req *ContractRequest, args ...interface{}) (string, error) {
	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", fmt.Errorf("account:%s not exist", req.From)
	}

	conOptions := &contract.ContractOptions{}
	conOptions.From = req.From
	conOptions.Data = req.Data
	conOptions.Abi = req.Abi
	limit, err := strconv.ParseUint(req.GasLimit, 10, 64)
	if err != nil {
		return "", err
	}
	price, ok := new(big.Int).SetString(req.GasPrice, 10)
	if !ok {
		return "", errors.New("gas price convert err")
	}
	conOptions.GasLimit = limit
	conOptions.GasPrice = price
	contract, err := contract.NewContract(conOptions)
	if err != nil {
		return "", err
	}
	exec, err := contract.Deploy(args...)
	if err != nil {
		return "", err
	}
	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(exec.GasPrice()).
		SetGasLimit(exec.GasLimit()).
		SetAction(exec).Build()
	selp, err := action.Sign(elp, *sender.Private())
	if err != nil {
		return "", err
	}
	request := &rpcmethod.SendActionRequest{Action: selp.Proto()}
	_, err = i.SendAction(request)
	if err != nil {
		return "", err
	}
	shash := hash.Hash256b(byteutil.Must(proto.Marshal(selp.Proto())))
	return hex.EncodeToString(shash[:]), nil
}

func (i *Iotx) sendAction(acc *account.Account, elp action.Envelope) (string, error) {
	sealed, err := action.Sign(elp, *acc.Private())
	if err != nil {
		return "", err
	}
	selp := sealed.Proto()
	request := &iotexapi.SendActionRequest{Action: selp}
	_, err = i.SendAction(request)
	if err != nil {
		return "", err
	}
	shash := hash.Hash256b(byteutil.Must(proto.Marshal(selp)))
	return hex.EncodeToString(shash[:]), nil
}
