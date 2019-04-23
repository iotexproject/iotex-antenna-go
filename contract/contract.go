// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/iotex-antenna-go/action"
)

// Contract defines contract
type Contract struct {
	// contract address
	Address string

	// contract abi for invoke contract
	ABI string

	// contract bytecode for deploy
	Data []byte
}

// New construct Contract instrance
func New(address, abi, data string) (*Contract, error) {
	if len(abi) == 0 {
		return nil, errors.New("must set contract abi")
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return &Contract{
		Address: address,
		ABI:     abi,
		Data:    dataBytes,
	}, nil
}

// DeployAction returns deploy contract Execution ActionCore
func (c *Contract) DeployAction(nonce uint64, gasLimit uint64, gasPrice *big.Int, args ...interface{}) (*action.ActionCore, error) {
	data := c.Data
	if len(c.Data) == 0 {
		return nil, errors.New("contract bytecode can not empty for deploy")
	}
	if len(args) > 0 {
		ab, err := c.EncodeArguments("", args...)
		if err != nil {
			return nil, err
		}
		data = append(data, ab...)
	}
	return action.NewExecution(nonce, gasLimit, gasPrice, big.NewInt(0), "", data)
}

// ExecuteAction returns invoke contract Execution ActionCore
func (c *Contract) ExecuteAction(nonce uint64, gasLimit uint64, gasPrice *big.Int, amount *big.Int, method string, args ...interface{}) (*action.ActionCore, error) {
	data, err := c.EncodeArguments(method, args...)
	if err != nil {
		return nil, err
	}
	return action.NewExecution(nonce, gasLimit, gasPrice, amount, c.Address, data)
}

// EncodeArguments encode method arguments to bytes.
func (c *Contract) EncodeArguments(method string, args ...interface{}) ([]byte, error) {
	reader := strings.NewReader(c.ABI)
	abiParam, err := abi.JSON(reader)
	if err != nil {
		return nil, err
	}
	return abiParam.Pack(method, args...)
}

// GetFuncHash returns contract method
func GetFuncHash(fun string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(fun))[:4])
}
