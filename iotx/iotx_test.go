// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-address/address"
)

const (
	host              = "api.testnet.iotex.one:80"
	accountPrivateKey = "9cdf22c5caa8a4d99eb674da27756b438c05c6b1e8995f4a0586745e2071b115"
	accountAddress    = "io14gnqxf9dpkn05g337rl7eyt2nxasphf5m6n0rd"
	to                = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
)

func TestTransfer(t *testing.T) {
	require := require.New(t)
	iotx, err := NewIotx(host)
	require.NoError(err)
	err = iotx.Accounts.AddAccount(accountPrivateKey)
	require.NoError(err)

	req := &TransferRequest{From: accountAddress, To: to, Value: "1000000000000000000", Payload: "", GasLimit: "1000000", GasPrice: "1000000000000"}
	err = iotx.SendTransfer(req)
	require.NoError(err)
}
func TestDeployContract(t *testing.T) {
	require := require.New(t)
	iotx, err := NewIotx(host)
	require.NoError(err)
	err = iotx.Accounts.AddAccount(accountPrivateKey)
	require.NoError(err)

	req := &ContractRequest{From: accountAddress, Data: "6080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630423a132811461005b578063252bd4d314610085578063bfe43b4c146100c3575b600080fd5b34801561006757600080fd5b50610073600435610191565b60408051918252519081900360200190f35b34801561009157600080fd5b5061009a610194565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100cf57600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261011c9436949293602493928401919081908401838280828437509497506101919650505050505050565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561015657818101518382015260200161013e565b50505050905090810190601f1680156101835780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b90565b60005473ffffffffffffffffffffffffffffffffffffffff16905600a165627a7a723058203fec9e60ea1eeb408d1fb0dcb9dc38c32b513302f817b39548a3d2c36b0772430029", Abi: `[{"constant":true,"inputs":[{"name":"x","type":"uint256"}],"name":"bar","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"getaddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"y","type":"string"}],"name":"barstring","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"a","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]"`, GasLimit: "1000000", GasPrice: "9000000000000000"}
	addr, err := address.FromString(to)
	require.NoError(err)
	var evmContractAddrHash common.Address
	copy(evmContractAddrHash[:], addr.Bytes())

	hash, err := iotx.DeployContract(req, evmContractAddrHash)
	require.NoError(err)
	require.NotNil(hash)
}
