package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-antenna-go/v2/utils/unit"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

func main() {
	conn, err := iotex.NewDefaultGRPCConn("api.testnet.iotex.one:443")
	if err != nil {
		log.Fatalf("connection error : %v", err)
	}
	defer conn.Close()

	creator, _ := account.HexStringToAccount("73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426")
	c := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), creator)

	abi, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"x","type":"uint256"}],"name":"set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_x","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`))
	if err != nil {
		log.Fatalf("JSON error : %v", err)
	}
	bytecode := "608060405234801561001057600080fd5b5060df8061001f6000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a7230582043be766a6a271867521090c3e12130ea8891a8f59d4833bc205a3e2e2c70b4730029"
	data, err := hex.DecodeString(bytecode)
	if err != nil {
		log.Fatalf("Hex Decoding error : %v", err)
	}

	actionHash, err := c.DeployContract(data).SetGasPrice(big.NewInt(int64(unit.Qev))).SetGasLimit(1000000).SetArgs(abi, big.NewInt(10)).Call(context.Background())
	if err != nil {
		log.Fatalf("deploy error: %v", err)
	}
	time.Sleep(10 * time.Second)
	receiptResponse, err := c.GetReceipt(actionHash).Call(context.Background())
	contractAddress := receiptResponse.GetReceiptInfo().GetReceipt().GetContractAddress()
	fmt.Println("Contract Address:", contractAddress)
}
