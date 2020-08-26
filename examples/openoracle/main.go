// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

var (
	gasPrice, _ = big.NewInt(0).SetString("1000000000000", 10)
	gasLimit    = uint64(1000000)
)

func main() {
	PrivateKey := os.Getenv("PrivateKey")
	if PrivateKey == "" {
		return
	}
	bin, err := ioutil.ReadFile("OpenOraclePriceData.bin")
	if err != nil {
		return
	}
	abi, err := ioutil.ReadFile("OpenOraclePriceData.abi")
	if err != nil {
		return
	}

	s, err := NewOpenOracleService(PrivateKey, string(abi), string(bin), "", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}
	r, err := s.Deploy(context.Background(), true)
	fmt.Println("hash", r, err)

	authClient, err := NewOpenOracleService(PrivateKey, string(abi), string(bin), "io1s50xy46vjtneh5m8jv6ync75m9vlj28qe0pr26", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}
	// get some data from coinbase,curl https://api.pro.coinbase.com/oracle, have to be latest data or it will fail
	mes := "0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000005f45b50000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000002a1ec56780000000000000000000000000000000000000000000000000000000000000006707269636573000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034254430000000000000000000000000000000000000000000000000000000000"
	sig := "187a858734a49fd7ffaa5ba57a44cf12dd4dbb73e142e6f70a3b3e27d9d00a8bb8e430c650eb58b33a79c5ef9ec8d3bf1b3937015a3dbbb00f50cce07b90e8a6000000000000000000000000000000000000000000000000000000000000001c"
	testMessage, err := hex.DecodeString(mes)
	if err != nil {
		return
	}
	testSignature, err := hex.DecodeString(sig)
	if err != nil {
		return
	}
	hash, err := authClient.Put(context.Background(), testMessage, testSignature)
	fmt.Println("put", hash, err)

	readClient, err := NewOpenOracleService("", string(abi), string(bin), "io1s50xy46vjtneh5m8jv6ync75m9vlj28qe0pr26", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}
	// io1ln4d4743f4rwyq2y7jyzf5xqnvdq8u4u9hma6k - 0xfCEAdAFab14d46e20144F48824d0C09B1a03F2BC
	ret, err := readClient.Get(context.Background(), "io1ln4d4743f4rwyq2y7jyzf5xqnvdq8u4u9hma6k", "BTC")
	fmt.Println("get", ret, err)
}
