// Copyright (c) 2026 IoTeX
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

// Example: send a typed (EIP-1559 DynamicFee) transfer on the IoTeX testnet.
//
//	IOTEX_KEY=<funded-private-key-hex> go run ./examples/typed_transfer
//
// The same chained API works for AccessList (TxType=1) by calling
// SetTxType(1) and SetGasPrice instead of SetGasTipCap/SetGasFeeCap.
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
)

const (
	endpoint  = "api.testnet.iotex.one:443"
	recipient = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
	chainID   = 2

	// DynamicFee in iotex's EIP-1559 model: TxType=2 selects DynamicFee,
	// the tip is what the sender is willing to pay above the base fee, and
	// the fee cap is the maximum total per-gas price the sender will pay.
	dynamicFeeTxType = 2
)

func main() {
	key := os.Getenv("IOTEX_KEY")
	if key == "" {
		log.Fatal("set IOTEX_KEY to a funded testnet private key (hex, no 0x prefix)")
	}

	acc, err := account.HexStringToAccount(key)
	if err != nil {
		log.Fatalf("invalid private key: %v", err)
	}
	to, err := address.FromString(recipient)
	if err != nil {
		log.Fatalf("invalid recipient: %v", err)
	}

	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("grpc dial: %v", err)
	}
	defer conn.Close()
	cli := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), chainID, acc)

	tip := new(big.Int).SetUint64(1_000_000_000_000) // 1 Qev (10^12)
	feeCap := new(big.Int).Mul(big.NewInt(2), tip)

	hash, err := cli.Transfer(to, big.NewInt(1)).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(tip).
		SetGasFeeCap(feeCap).
		SetGasLimit(50_000).
		Call(context.Background())
	if err != nil {
		log.Fatalf("send: %v", err)
	}
	fmt.Printf("sent DynamicFee transfer: %x\n", hash[:])
}
