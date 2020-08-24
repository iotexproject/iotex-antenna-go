package main

import (
	"context"
	"fmt"
	"math/big"
)

const xrc20abi = `[
    {
        "constant": false,
        "inputs": [
            {
                "name": "_to",
                "type": "address"
            },
            {
                "name": "_value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`

var (
	gasPrice, _ = big.NewInt(0).SetString("1000000000000", 10)
	gasLimit    = uint64(1000000)
)

func main() {
	s, err := NewIotexService("583aa7b02dbba44d257cad116e7e427d4b2040a0079c348d83636e100a4a4039", xrc20abi, "io10gjq3edy6n953q0j3lm2p7zh3a8d3s4em0w6cj", gasPrice, gasLimit, "api.testnet.iotex.one:80", false)
	if err != nil {
		return
	}
	r, err := s.Transfer(context.Background(), "io1zk6gqq0m2z9ytlu77t76e3632ezy39fa83xjnn", big.NewInt(10))
	fmt.Println("hash", r, err)
}
