# iotex-antenna-go

[![CircleCI](https://circleci.com/gh/iotexproject/iotex-antenna-go.svg?style=svg)](https://circleci.com/gh/iotexproject/iotex-antenna-go)
[![Go version](https://img.shields.io/badge/go-1.11.5-blue.svg)](https://github.com/moovweb/gvm)
[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

This is the the official Go implementation of IoTeX SDK! Please refer to IoTeX [whitepaper](https://iotex.io/research) and the [protocol](https://github.com/iotexproject/iotex-core) for details.

## Get Started

### Minimum Requirements

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |

### Add Dependency

```
// go mod
go get github.com/iotexproject/iotex-antenna-go/v2
```

### Code It Up
The below example code shows the 4 easy steps to send a transaction to IoTeX blockchain
1. connect to the chain's RPC endpoint
2. create an account by importing a private key
3. create a client and generate an action sender
4. send the transaction to the chain

```
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

const (
	mainnetRPC     = "api.iotex.one:443"
	testnetRPC     = "api.testnet.iotex.one:443"
	mainnetChainID = 1
	testnetChainID = 2
)

func main() {
	// Create grpc connection
	conn, err := iotex.NewDefaultGRPCConn(testnetRPC)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Add account by private key
	acc, err := account.HexStringToAccount("...")
	if err != nil {
		log.Fatal(err)
	}

	// create client
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), acc)
	
	// send the transfer to chain
	to, err := address.FromString("io1zq5g9c5c3hqw9559ks4anptkpumxgsjfn2e4ke")
	if err != nil {
		log.Fatal(err)
	}
	hash, err := c.Transfer(to, big.NewInt(10)).SetChainID(testnetChainID).SetGasPrice(big.NewInt(100000000000)).SetGasLimit(20000).Call(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("transaction hash = %x\n", hash)
}
```

### More Examples
There are three examples demostrating the use of this SDK on Testnet. You can `make examples` to build and try:
- `./examples/chaininfo` shows **how to use the SDK to pull chain, block, action and delegates info**
- `./examples/openoracle` shows **how to deploy and invoke [Open Oracle Contracts](https://github.com/compound-finance/open-oracle)**
- `./examples/xrc20tokens` shows **how to deploy and invoke XRC20 tokens**
