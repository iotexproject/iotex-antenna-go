// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	host = "api.iotex.one:80"
)

func TestServer_Deploy(t *testing.T) {
	require := require.New(t)
	bin := "6080604052348015600f57600080fd5b5060a18061001e6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063febb0f7e146044575b600080fd5b348015604f57600080fd5b506056606c565b6040518082815260200191505060405180910390f35b600060649050905600a165627a7a72305820631492f34cc1852dd0cfbfaa0631b86e9c06d7ba577caf330b4535651fd1408a0029"
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(9000000000000)
	sct, err := NewContract(host, bin, gasLimit, gasPrice)
	require.NoError(err)
	accountPrivateKey := "8c379a71721322d16912a88b1602c5596ca9e99a5f70777561c3029efa71a435"
	accountAddress := "io1ns7y0pxmklk8ceattty6n7makpw76u770u5avy"
	sct.SetExecutor(accountAddress, accountPrivateKey)
	hash, err := sct.Deploy()
	require.NoError(err)
	receipt, err := sct.CheckCallResult(hash)
	require.NoError(err)
	fmt.Println("receipt contract:", receipt.ContractAddress)
	sct.SetContractAddress(receipt.ContractAddress)
	fmt.Println("contract:", sct.ContractAddress())
	ret, err := sct.CallMethod("febb0f7e")
	require.NoError(err)
	fmt.Println(ret)
}
