# iotex-antenna-go

[![CircleCI](https://circleci.com/gh/iotexproject/iotex-antenna-go.svg?style=svg)](https://circleci.com/gh/iotexproject/iotex-antenna-go)
[![Go version](https://img.shields.io/badge/go-1.11.5-blue.svg)](https://github.com/moovweb/gvm)
[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

This is the the official Go implementation of IoTeX SDK! Please refer to IoTeX [whitepaper](https://iotex.io/research) and the [protocol](https://github.com/iotexproject/iotex-core) for details.

## Get started

### Minimum requirements

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |

### Add to your project

```
// go mod
go get github.com/iotexproject/iotex-antenna-go/v2
```

### Example

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
	host = "api.testnet.iotex.one:443"
)

func main() {
	// Create grpc connection
	conn, err := iotex.NewDefaultGRPCConn(host)
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
	
	// transfer
	to, err := address.FromString("to...")
	if err != nil {
		log.Fatal(err)
	}
	hash, err := c.Transfer(to, big.NewInt(10)).Call(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
```
The other examples are under the examples folder.
1. The chaininfo folder shows that how to get chain info, block info, tx info and delegate info,etc.
2. The tokens folder shows that how to transfer XRC20 tokens in the xrc20 contract.
3. The contract folder shows that how to deploy contract and call contract function.