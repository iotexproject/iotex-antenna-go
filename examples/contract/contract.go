package main

import (
	"context"
	"fmt"
	"math/big"
)

var (
	gasPrice, _ = big.NewInt(0).SetString("1000000000000", 10)
	gasLimit    = uint64(1000000)
)

func main() {
	s, err := NewIotexService("583aa7b02dbba44d257cad116e7e427d4b2040a0079c348d83636e100a4a4039", IotexTokenABI, IotexTokenBin, "", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}

	initialSupply := big.NewInt(2000000000)
	tokenName := "IOTX"
	tokenSymbol := "IOTX"
	r, err := s.Deploy(context.Background(), true, initialSupply, tokenName, tokenSymbol)
	fmt.Println("hash", r, err)

	readOnly, err := NewIotexService("", IotexTokenABI, "", "io1eq786nwuu6ygw4ct075gfp3u2f6xgmp8f5hygq", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}

	b, err := readOnly.BalanceOf(context.Background(), "io1tdfyk5gqrfas22am6sw732twxyjcnl6xqe850s")
	fmt.Println("balance", b, err)
}
