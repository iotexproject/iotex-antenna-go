// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/iotexproject/iotex-antenna-go/contract"
	"github.com/iotexproject/iotex-antenna-go/rpcmethod"

	"github.com/iotexproject/iotex-antenna-go/action"

	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/utils"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

// Error strings
var (
	// ErrAccountNotExist indicates error for not exist account
	ErrAccountNotExist = fmt.Errorf("account not exist")
	// ErrAmount indicates error for error amount convert
	ErrAmount = fmt.Errorf("error amount")
)

// Iotx service RPCMethod and Accounts
type Iotx struct {
	*rpcmethod.RPCMethod
	Accounts *account.Accounts
}

// New return Iotx instance
func New(host string) (*Iotx, error) {
	rpc, err := rpcmethod.NewRPCMethod(host)
	if err != nil {
		return nil, err
	}
	iotx := &Iotx{rpc, account.NewAccounts()}
	return iotx, nil
}

func (i *Iotx) normalizeGas(acc *account.Account, ac *action.ActionCore, gasLimit, gasPrice string) (uint64, *big.Int, error) {
	var limit uint64
	var price *big.Int
	if gasLimit == "" {
		sealed, err := ac.Sign(*acc.Private())
		if err != nil {
			return 0, nil, err
		}
		request := &iotexapi.EstimateGasForActionRequest{Action: sealed.Action}
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

// SendTransfer invoke send transfer action by rpc
func (i *Iotx) SendTransfer(req *TransferRequest) (string, error) {
	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", ErrAccountNotExist
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
		return "", ErrAmount
	}

	act, err := action.NewTransfer(nonce, 0, big.NewInt(0), amount, req.To, req.Payload)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = action.NewTransfer(nonce, gasLimit, gasPrice, amount, req.To, req.Payload)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

// DeployContract invoke execution action for deploy contract
func (i *Iotx) DeployContract(req *ContractRequest, args ...interface{}) (string, error) {
	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", ErrAccountNotExist
	}

	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	ctr, err := contract.New("", req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	act, err := ctr.DeployAction(nonce, 0, big.NewInt(0), args...)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = ctr.DeployAction(nonce, gasLimit, gasPrice, args...)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

func (i *Iotx) sendAction(acc *account.Account, ta *action.ActionCore) (string, error) {
	sealed, err := ta.Sign(*acc.Private())
	if err != nil {
		return "", err
	}
	request := &iotexapi.SendActionRequest{Action: sealed.Action}
	response, err := i.SendAction(request)
	if err != nil {
		return "", err
	}
	return response.ActionHash, nil
}

// ExecuteContract returns execute contract method action hash
func (i *Iotx) ExecuteContract(req *ContractRequest, args ...interface{}) (string, error) {
	if req.Address == "" || req.Method == "" {
		return "", errors.New("contract address and method can not empty")
	}

	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return "", ErrAmount
	}

	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", ErrAccountNotExist
	}

	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	ctr, err := contract.New(req.Address, req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	act, err := ctr.ExecuteAction(nonce, 0, big.NewInt(0), amount, req.Method, args...)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = ctr.ExecuteAction(nonce, gasLimit, gasPrice, amount, req.Method, args...)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

// ReadContractByHash returns execute contract method result by action hash
func (i *Iotx) ReadContractByHash(hash string) (string, error) {
	actionResponse, err := i.GetActions(&iotexapi.GetActionsRequest{Lookup: &iotexapi.GetActionsRequest_ByHash{
		ByHash: &iotexapi.GetActionByHashRequest{
			ActionHash:   hash,
			CheckPending: true,
		},
	}})
	if err != nil {
		return "", err
	}

	request := &iotexapi.ReadContractRequest{Action: actionResponse.ActionInfo[0].Action}
	response, err := i.ReadContract(request)
	if err != nil {
		return "", err
	}

	return response.Data, nil
}

// ReadContractByMethod returns execute contract view method result
func (i *Iotx) ReadContractByMethod(req *ContractRequest, args ...interface{}) (string, error) {
	if req.Address == "" || req.Method == "" {
		return "", errors.New("contract address and method can not empty")
	}

	sender, ok := i.Accounts.GetAccount(req.From)
	if !ok {
		return "", ErrAccountNotExist
	}

	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	ctr, err := contract.New(req.Address, req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	act, err := ctr.ExecuteAction(nonce, 0, big.NewInt(0), big.NewInt(0), req.Method, args...)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = ctr.ExecuteAction(nonce, gasLimit, gasPrice, big.NewInt(0), req.Method, args...)
	if err != nil {
		return "", err
	}
	sealed, err := act.Sign(*sender.Private())
	if err != nil {
		return "", err
	}
	request := &iotexapi.ReadContractRequest{Action: sealed.Action}
	response, err := i.ReadContract(request)
	if err != nil {
		return "", err
	}

	return response.Data, nil
}
