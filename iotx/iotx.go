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

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-antenna-go/action"
	"github.com/iotexproject/iotex-antenna-go/contract"
	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
	"github.com/iotexproject/iotex-antenna-go/utils"
)

// Error strings
var (
	// ErrAmount indicates error for error amount convert
	ErrAmount = fmt.Errorf("error amount")
)

// Iotx service RPCMethod and Accounts
type Iotx struct {
	*rpcmethod.RPCMethod
	Accounts *account.Accounts
}

// NewIotx return Iotx instance
func NewIotx(host string, tls bool) (*Iotx, error) {
	var rpc *rpcmethod.RPCMethod
	var err error

	if tls {
		rpc, err = rpcmethod.NewRPCWithTLSEnabled(host)
	} else {
		rpc, err = rpcmethod.NewRPCMethod(host)
	}
	if err != nil {
		return nil, err
	}
	return &Iotx{
		rpc,
		account.NewAccounts(),
	}, nil
}

func (i *Iotx) normalizeGas(acc account.Account, ac *action.IotexActionCore, gasLimit, gasPrice string) (uint64, *big.Int, error) {
	var limit uint64
	var price *big.Int
	if gasLimit == "" {
		sealed, err := ac.Sign(acc)
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
	sender, err := i.Accounts.GetAccount(req.From)
	if err != nil {
		return "", err
	}

	amount, ok := big.NewInt(0).SetString(req.Value, 10)
	if !ok {
		return "", ErrAmount
	}
	data, err := hex.DecodeString(req.Payload)
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

	act, err := action.NewTransfer(nonce, 0, big.NewInt(0), amount, req.To, data)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = action.NewTransfer(nonce, gasLimit, gasPrice, amount, req.To, data)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

// DeployContract invoke execution action for deploy contract
func (i *Iotx) DeployContract(req *ContractRequest, args ...interface{}) (string, error) {
	sender, err := i.Accounts.GetAccount(req.From)
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

	ctr, err := contract.NewContract("", req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	actData, err := ctr.DeployData(args...)
	if err != nil {
		return "", err
	}
	act, err := action.NewExecution(nonce, 0, big.NewInt(0), big.NewInt(0), "", actData)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = action.NewExecution(nonce, gasLimit, gasPrice, big.NewInt(0), "", actData)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

func (i *Iotx) sendAction(acc account.Account, ta *action.IotexActionCore) (string, error) {
	sealed, err := ta.Sign(acc)
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
	sender, err := i.Accounts.GetAccount(req.From)
	if err != nil {
		return "", err
	}

	if req.Address == "" || req.Method == "" {
		return "", errors.New("contract address and method can not empty")
	}
	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return "", ErrAmount
	}
	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	ctr, err := contract.NewContract(req.Address, req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	actData, err := ctr.ExecuteData(req.Method, args...)
	if err != nil {
		return "", err
	}
	act, err := action.NewExecution(nonce, 0, big.NewInt(0), amount, req.Address, actData)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = action.NewExecution(nonce, gasLimit, gasPrice, amount, req.Address, actData)
	if err != nil {
		return "", err
	}
	return i.sendAction(sender, act)
}

// ReadContractByHash returns execute contract method result by action hash
func (i *Iotx) ReadContractByHash(hash string) (string, error) {
	actionResponse, err := i.GetActions(&iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash:   hash,
				CheckPending: true,
			},
		},
	})
	if err != nil {
		return "", err
	}

	action := actionResponse.ActionInfo[0].Action
	if action == nil {
		return "", errors.Errorf("action %x is not valid", hash)
	}
	pk, _ := crypto.BytesToPublicKey(action.SenderPubKey)
	caller, _ := address.FromBytes(pk.Hash())

	execution := action.GetCore().GetExecution()
	if execution == nil {
		return "", errors.Errorf("action %x is not an execution", hash)
	}

	response, err := i.ReadContract(&iotexapi.ReadContractRequest{
		Execution:     execution,
		CallerAddress: caller.String(),
	})
	if err != nil {
		return "", err
	}
	return response.Data, nil
}

// ReadContractByMethod returns execute contract view method result
func (i *Iotx) ReadContractByMethod(req *ContractRequest, args ...interface{}) (string, error) {
	sender, err := i.Accounts.GetAccount(req.From)
	if err != nil {
		return "", err
	}

	if req.Address == "" || req.Method == "" {
		return "", errors.New("contract address and method can not empty")
	}
	// get account nonce
	accountReq := &rpcmethod.GetAccountRequest{Address: req.From}
	res, err := i.GetAccount(accountReq)
	if err != nil {
		return "", err
	}
	nonce := res.AccountMeta.PendingNonce

	ctr, err := contract.NewContract(req.Address, req.Abi, req.Data)
	if err != nil {
		return "", err
	}

	actData, err := ctr.ExecuteData(req.Method, args...)
	if err != nil {
		return "", err
	}
	act, err := action.NewExecution(nonce, 0, big.NewInt(0), big.NewInt(0), req.Address, actData)
	if err != nil {
		return "", err
	}
	gasLimit, gasPrice, err := i.normalizeGas(sender, act, req.GasLimit, req.GasPrice)
	if err != nil {
		return "", err
	}

	act, err = action.NewExecution(nonce, gasLimit, gasPrice, big.NewInt(0), req.Address, actData)
	if err != nil {
		return "", err
	}
	request := &iotexapi.ReadContractRequest{
		Execution:     act.GetExecution(),
		CallerAddress: sender.Address(),
	}
	response, err := i.ReadContract(request)
	if err != nil {
		return "", err
	}
	return response.Data, nil
}
