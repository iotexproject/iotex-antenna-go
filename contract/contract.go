// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/iotexproject/iotex-core/action"
)

type CustomOptions struct {
	Address  string
	Abi      string
	Data     string
	From     string
	GasPrice *big.Int
	GasLimit uint64
}

type contractOptions struct {
	CustomOptions
}
type Contract struct {
	options contractOptions
}

func NewContract(options CustomOptions) *Contract {
	return &Contract{contractOptions{options}}
}
func (c *Contract) ABI() string {
	return c.options.Abi
}
func (c *Contract) Address() string {
	return c.options.Address
}
func (c *Contract) Deploy(args ...interface{}) (*action.Execution, error) {
	data, err := hex.DecodeString(c.options.Data)
	if err != nil {
		return nil, err
	}
	arg, err := c.EncodeArguments("", args...)
	if err != nil {
		return nil, err
	}
	data = append(data, arg...)

	execution, err := action.NewExecution("", 0, big.NewInt(0), c.options.GasLimit, c.options.GasPrice, data)
	return execution, err
}

func (c *Contract) EncodeArguments(method string, args ...interface{}) ([]byte, error) {
	reader := strings.NewReader(c.options.Abi)
	abiParam, err := abi.JSON(reader)
	if err != nil {
		return nil, err
	}
	return abiParam.Pack(method, args...)
}

func GetFuncHash(fun string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(fun))[:4])
}
