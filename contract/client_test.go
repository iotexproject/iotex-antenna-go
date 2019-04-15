// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/iotexproject/iotex-core/address"

	"github.com/stretchr/testify/require"
)

const (
	host = "api.iotex.one:80"
)

func TestContract(t *testing.T) {
	require := require.New(t)
	bin := "608060405234801561001057600080fd5b5060405160208061022e833981016040525160008054600160a060020a03909216600160a060020a03199092169190911790556101dc806100526000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630423a132811461005b578063252bd4d314610085578063bfe43b4c146100c3575b600080fd5b34801561006757600080fd5b50610073600435610191565b60408051918252519081900360200190f35b34801561009157600080fd5b5061009a610194565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156100cf57600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261011c9436949293602493928401919081908401838280828437509497506101919650505050505050565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561015657818101518382015260200161013e565b50505050905090810190601f1680156101835780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b90565b60005473ffffffffffffffffffffffffffffffffffffffff16905600a165627a7a723058203fec9e60ea1eeb408d1fb0dcb9dc38c32b513302f817b39548a3d2c36b0772430029"
	abi := `[{"constant":true,"inputs":[{"name":"x","type":"uint256"}],"name":"bar","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"getaddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"y","type":"string"}],"name":"barstring","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"a","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`

	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(9000000000000)
	sct, err := NewContract(host, bin, abi, gasLimit, gasPrice)
	require.NoError(err)
	accountPrivateKey := "8c379a71721322d16912a88b1602c5596ca9e99a5f70777561c3029efa71a435"
	accountAddress := "io1ns7y0pxmklk8ceattty6n7makpw76u770u5avy"
	sct.SetOwner(accountAddress, accountPrivateKey)
	addr, err := address.FromString("io1ns7y0pxmklk8ceattty6n7makpw76u770u5avy")
	require.NoError(err)
	var evmContractAddrHash common.Address
	copy(evmContractAddrHash[:], addr.Bytes())

	hash, err := sct.Deploy(evmContractAddrHash)
	require.NoError(err)
	receipt, err := sct.CheckCallResult(hash)
	require.NoError(err)
	sct.SetContractAddress(receipt.ContractAddress)
	sct.SetExecutor(accountAddress, accountPrivateKey)

	ret, err := sct.CallMethod("bar", big.NewInt(10))
	require.NoError(err)
	require.Equal("*big.Int", reflect.TypeOf(ret).String())
	require.Equal(0, ret.(*big.Int).Cmp(big.NewInt(10)))

	ret, err = sct.CallMethod("barstring", "foobar")
	require.NoError(err)
	require.Equal("string", reflect.TypeOf(ret).String())
	require.Equal("foobar", ret.(string))

	ret, err = sct.CallMethod("getaddress")
	require.NoError(err)
	require.Equal("common.Address", reflect.TypeOf(ret).String())
	retAddr, ok := ret.(common.Address)
	require.True(ok)
	require.Equal(evmContractAddrHash.String(), retAddr.String())
}
