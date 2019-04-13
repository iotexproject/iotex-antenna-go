// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"log"
	"math/big"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
	"github.com/pkg/errors"
)

type (
	// Contract is contract interface
	Contract interface {
		Deploy(...[]byte) (string, error)
		CallMethod(string, ...[]byte) (string, error)
		SendToChain([]byte, bool) (string, error)
		CheckCallResult(string) (*iotextypes.Receipt, error)
		ContractAddress() string
		SetContractAddress(string) Contract
		SetOwner(string, string) Contract
		SetExecutor(string, string) Contract
		RunAsOwner() Contract
	}

	contract struct {
		endpoint        string // blockchain service endpoint
		ownerAddress    string // owner of the smart contract
		ownerPk         string // owner's private key
		contractAddress string // address of the smart contract
		executorAddress string // executor of the smart contract
		executorPk      string // private key of executor
		rpc             *rpcmethod.RPCMethod
		codeBin         string   // code of the smart contract
		gasLimit        uint64   //gas limit
		gasPrice        *big.Int //gas price
	}
)

// NewContract creates a new contract
func NewContract(endpoint, bin string, gasLimit uint64, gasPrice *big.Int) (Contract, error) {
	ret := &contract{endpoint: endpoint, codeBin: bin, gasLimit: gasLimit, gasPrice: gasPrice}
	rpcmethod, err := rpcmethod.NewRPCMethod(endpoint)
	if err != nil {
		return nil, err
	}
	ret.rpc = rpcmethod
	return ret, nil
}
func (c *contract) Deploy(args ...[]byte) (string, error) {
	data, err := hex.DecodeString(c.codeBin)
	if err != nil {
		return "", err
	}
	for _, arg := range args {
		if arg != nil {
			if len(arg) < 32 {
				value := hash.BytesToHash256(arg)
				data = append(data, value[:]...)
			} else {
				data = append(data, arg...)
			}
		}
	}
	// deploy send to empty address
	return c.SetContractAddress("").SendToChain(data, false)
}
func (c *contract) method(method string, args ...[]byte) ([]byte, error) {
	data, err := hex.DecodeString(method)
	if err != nil {
		return nil, err
	}
	if len(data) != 4 {
		return nil, errors.Errorf("invalid method id format, length = %d", len(data))
	}
	for _, arg := range args {
		if arg != nil {
			if len(arg) < 32 {
				value := hash.BytesToHash256(arg)
				data = append(data, value[:]...)
			} else {
				data = append(data, arg...)
			}
		}
	}
	return data, nil
}
func (c *contract) CallMethod(method string, args ...[]byte) (string, error) {
	data, err := c.method(method, args...)
	if err != nil {
		return "", err
	}
	return c.SendToChain(data, true)
}
func (c *contract) ExecMethod(method string, args ...[]byte) (string, error) {
	data, err := c.method(method, args...)
	if err != nil {
		return "", err
	}
	return c.SendToChain(data, false)
}

func (c *contract) SendToChain(data []byte, readOnly bool) (string, error) {
	response, err := c.rpc.GetAccount(c.executorAddress)
	if err != nil {
		return "", err
	}
	nonce := response.AccountMeta.PendingNonce
	tx, err := action.NewExecution(
		c.contractAddress,
		nonce,
		big.NewInt(0),
		c.gasLimit,
		c.gasPrice,
		data)
	if err != nil {
		return "", err
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(nonce).
		SetGasPrice(c.gasPrice).
		SetGasLimit(c.gasLimit).
		SetAction(tx).Build()

	prvKey, err := keypair.HexStringToPrivateKey(c.executorPk)
	if err != nil {
		return "", keypair.ErrInvalidKey
	}
	selp, err := action.Sign(elp, prvKey)
	if err != nil {
		return "", err
	}
	if readOnly {
		response, err := c.rpc.ReadContract(selp.Proto(), true)
		if err != nil {
			return "", err
		}
		return response.Data, nil
	}
	_, err = c.rpc.SendExecution(c.contractAddress, c.executorPk, nonce, big.NewInt(0), c.gasLimit, c.gasPrice, data)
	h := selp.Hash()
	hex := hex.EncodeToString(h[:])
	if err != nil {
		return hex, errors.Wrapf(err, "tx 0x%s failed to send to Blockchain", hex)
	}
	return hex, nil
}

func (c *contract) CheckCallResult(h string) (*iotextypes.Receipt, error) {
	var rec *iotextypes.Receipt
	// max retry 120 times with interval = 500ms
	num := 120
	err := backoff.Retry(func() error {
		var err error
		rec, err = c.checkCallResult(h)
		log.Printf("Hash: %s times: %d", h, num)
		num--
		return err
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Millisecond*500), uint64(num)))
	return rec, err
}

func (c *contract) checkCallResult(h string) (*iotextypes.Receipt, error) {
	response, err := c.rpc.GetReceiptByAction(h)
	if err != nil {
		return nil, err
	}
	if response.ReceiptInfo.Receipt.Status != 1 {
		return nil, errors.Errorf("tx 0x%s execution on Blockchain failed", h)
	}
	return response.ReceiptInfo.Receipt, nil
}

func (c *contract) ContractAddress() string {
	return c.contractAddress
}

func (c *contract) SetContractAddress(addr string) Contract {
	c.contractAddress = addr
	return c
}

func (c *contract) SetOwner(owner, pk string) Contract {
	c.ownerAddress = owner
	c.ownerPk = pk
	return c
}

func (c *contract) SetExecutor(exec, pk string) Contract {
	c.executorAddress = exec
	c.executorPk = pk
	return c
}

func (c *contract) RunAsOwner() Contract {
	c.executorAddress = c.ownerAddress
	c.executorPk = c.ownerPk
	return c
}
