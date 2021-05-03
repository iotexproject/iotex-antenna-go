package iotex

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

// Data is the data returned from read contract.
type Data struct {
	method string
	abi    *abi.ABI
	Raw    []byte
}

// Unmarshal unmarshals data into a data holder object.
func (d Data) Unmarshal() ([]interface{}, error) { return d.abi.Unpack(d.method, d.Raw) }

type contract struct {
	address address.Address
	abi     *abi.ABI
	api     iotexapi.APIServiceClient
	account account.Account
}

func (c *contract) Read(method string, args ...interface{}) ReadContractCaller {
	return &readContractCaller{
		api: c.api,
		contractArgs: contractArgs{
			contract: c.address,
			abi:      c.abi,
			method:   method,
			args:     args,
		},
	}
}

func (c *contract) Execute(method string, args ...interface{}) ExecuteContractCaller {
	return &executeContractCaller{
		sendActionCaller: sendActionCaller{
			account: c.account,
			api:     c.api,
		},
		contractArgs: contractArgs{
			contract: c.address,
			abi:      c.abi,
			method:   method,
			args:     args,
		},
	}
}

type readOnlyContract struct {
	address address.Address
	abi     *abi.ABI
	api     iotexapi.APIServiceClient
}

func (c *readOnlyContract) Read(method string, args ...interface{}) ReadContractCaller {
	return &readContractCaller{
		api: c.api,
		contractArgs: contractArgs{
			contract: c.address,
			abi:      c.abi,
			method:   method,
			args:     args,
		},
	}
}
