// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/crypto"
)

// Contract defines contract interface
type Contract interface {
	// Address returns contract address
	Address() string

	// ABI returns contract abi
	ABI() string

	// DeployData returns deploy contract Execution packed data
	DeployData(args ...interface{}) ([]byte, error)

	// ExecuteData returns invoke contract Execution packed data
	ExecuteData(method string, args ...interface{}) ([]byte, error)
}

type contract struct {
	// contract address
	address string

	// contract abi for invoke contract
	abi string

	// contract bytecode for deploy
	data []byte
}

// NewContract returns construct Contract instance
func NewContract(address, abi, data string) (Contract, error) {
	if len(abi) == 0 {
		return nil, errors.New("must set contract abi")
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return &contract{
		address: address,
		abi:     abi,
		data:    dataBytes,
	}, nil
}

// Address returns contract address
func (c *contract) Address() string {
	return c.address
}

// ABI returns contract abi
func (c *contract) ABI() string {
	return c.abi
}

// DeployAction returns deploy contract Execution ActionCore
func (c *contract) DeployData(args ...interface{}) ([]byte, error) {
	data := c.data
	if len(c.data) == 0 {
		return nil, errors.New("contract bytecode can not empty for deploy")
	}
	if len(args) > 0 {
		ab, err := c.EncodeArguments("", args...)
		if err != nil {
			return nil, err
		}
		data = append(data, ab...)
	}
	return data, nil
}

// ExecuteData returns invoke contract Execution ActionCore
func (c *contract) ExecuteData(method string, args ...interface{}) ([]byte, error) {
	return c.EncodeArguments(method, args...)
}

// EncodeArguments encode method arguments to bytes.
func (c *contract) EncodeArguments(method string, args ...interface{}) ([]byte, error) {
	reader := strings.NewReader(c.abi)
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
