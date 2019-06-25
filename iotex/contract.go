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
func (d Data) Unmarshal(v interface{}) error { return d.abi.Unpack(v, d.method, d.Raw) }

type contract struct {
	address address.Address
	abi     *abi.ABI
	api     iotexapi.APIServiceClient
	account account.Account
}

func (c *contract) Read(method string, args ...interface{}) ReadContractCaller {
	return &readContractCaller{
		method: method,
		args:   args,
		rc: &readOnlyContract{
			address: c.address,
			abi:     c.abi,
			api:     c.api,
		},
		sender: c.account.Address(),
	}
}
func (c *contract) Execute(method string, args ...interface{}) ExecuteContractCaller {
	return &executeContractCaller{
		abi:      c.abi,
		api:      c.api,
		contract: c.address,
		account:  c.account,
		method:   method,
		args:     args,
	}
}

type readOnlyContract struct {
	address address.Address
	abi     *abi.ABI
	api     iotexapi.APIServiceClient
}

func (c *readOnlyContract) Read(method string, args ...interface{}) ReadContractCaller {
	sender, _ := address.FromString("io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6")
	return &readContractCaller{
		method: method,
		args:   args,
		rc:     c,
		sender: sender,
	}
}
