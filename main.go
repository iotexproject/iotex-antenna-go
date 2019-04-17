// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/iotexproject/iotex-antenna-go/antenna"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
)

func main() {
	var blk iotextypes.Block
	fmt.Printf("Hello Antenna %s\n", blk.Header.ProducerPubkey)

	svr, err := antenna.NewRPCMethod("api.iotex.one:80")
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := svr.SuggestGasPrice(&antenna.SuggestGasPriceRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
