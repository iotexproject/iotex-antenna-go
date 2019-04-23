// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iotexproject/iotex-antenna-go/action"
)

// Contract defines contract
type Contract struct {
	// contract abi for invoke contract
	ABI string
}

// DeployAction returns deploy contract Execution ActionCore
func DeployAction(nonce uint64, gasLimit uint64, gasPrice *big.Int, data string) (*action.ActionCore, error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return action.NewExecution(nonce, gasLimit, gasPrice, big.NewInt(0), "", dataBytes)
}

/*
func NewContract(options *ContractOptions) (Contract, error) {
	err := validate(options)
	if err != nil {
		return nil, err
	}
	return &contract{&contractOptions{options}}, nil
}

func (c *contract) ABI() string {
	return c.options.Abi
}

func (c *contract) Address() string {
	return c.options.Address
}

// Deploy args is used for this contract's constructor
func (c *contract) Deploy(args ...interface{}) (*action.Execution, error) {
	data, err := hex.DecodeString(c.options.Data)
	if err != nil {
		return nil, err
	}
	arg, err := c.encodeArguments("", args...)
	if err != nil {
		return nil, err
	}
	data = append(data, arg...)

	execution, err := action.NewExecution("", 0, big.NewInt(0), c.options.GasLimit, c.options.GasPrice, data)
	return execution, err
}

func (c *contract) encodeArguments(method string, args ...interface{}) ([]byte, error) {
	reader := strings.NewReader(c.ContractOptions.Abi)
	abiParam, err := abi.JSON(reader)
	if err != nil {
		return nil, err
	}
	return abiParam.Pack(method, args...)
}

func GetFuncHash(fun string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(fun))[:4])
}

func validate(options *ContractOptions) error {
	if options.Abi == "" || options.Data == "" || options.From == "" {
		return errors.New("some params is empty")
	}
	return nil
}

*/

// GetFuncHash returns contract method
func GetFuncHash(fun string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(fun))[:4])
}
