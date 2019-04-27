# iotex-antenna-go

[![CircleCI](https://circleci.com/gh/iotexproject/iotex-antenna-go.svg?style=svg)](https://circleci.com/gh/iotexproject/iotex-antenna-go)
[![Go version](https://img.shields.io/badge/go-1.11.5-blue.svg)](https://github.com/moovweb/gvm)
[![LICENSE](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Welcome to the official Go implementation of IoTeX Golang SDK! IoTeX is building the next generation of the decentralized 
network for IoT powered by scalability- and privacy-centric blockchains. Please refer to IoTeX
[whitepaper](https://iotex.io/academics) for details.

## Get started

### Minimum requirements

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |
| [Dep](https://golang.github.io/dep/) | &ge; 0.5.0 | Dependency management tool, required only when you update dependencies |

### Add to your project

```
// dep
dep ensure -add github.com/iotexproject/iotex-antenna-go

// go mod
go get github.com/iotexproject/iotex-antenna-go
```

### Sample

```
package main

import (
	"log"

	"github.com/iotexproject/iotex-antenna-go/antenna"
)

const (
	host = "api.testnet.iotex.one:80"
)

func main() {
	antenna, err := antenna.NewAntenna(host)

	if err != nil {
		log.Fatalf("New antenna error: %v", err)
	}

	// Add account by private key
	antenna.Iotx.Accounts.PrivateKeyToAccount("...")

	// transfer
	antenna.Iotx.SendTransfer(...)

	// deploy contract
	antenna.Iotx.DeployContract(...)
}
```