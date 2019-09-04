package iotex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
)

type authedClient struct {
	client

	account account.Account
}

// NewAuthedClient creates an AuthedClient using given account's credentials.
func NewAuthedClient(api iotexapi.APIServiceClient, a account.Account) AuthedClient {
	return &authedClient{
		client: client{
			api: api,
		},
		account: a,
	}
}

func (c *authedClient) Contract(co address.Address, abi abi.ABI) Contract {
	return &contract{
		address: co,
		abi:     &abi,
		api:     c.api,
		account: c.account,
	}
}

func (c *authedClient) Transfer(to address.Address, value *big.Int) TransferCaller {
	return &transferCaller{
		account:   c.account,
		api:       c.api,
		amount:    value,
		recipient: to,
	}
}

func (c *authedClient) ClaimReward(value *big.Int) ClaimRewardCaller {
	return &claimRewardCaller{
		account: c.account,
		api:     c.api,
		amount:  value,
	}
}

func (c *authedClient) DeployContract(data []byte) DeployContractCaller {
	return &deployContractCaller{
		account: c.account,
		api:     c.api,
		data:    data,
	}
}

func (c *authedClient) Account() account.Account { return c.account }

// NewReadOnlyClient creates a ReadOnlyClient.
func NewReadOnlyClient(c iotexapi.APIServiceClient) ReadOnlyClient {
	return &client{api: c}
}

type client struct {
	api iotexapi.APIServiceClient
}

func (c *client) ReadOnlyContract(contract address.Address, abi abi.ABI) ReadOnlyContract {
	return &readOnlyContract{
		address: contract,
		abi:     &abi,
		api:     c.api,
	}
}

func (c *client) GetReceipt(actionHash hash.Hash256) GetReceiptCaller {
	return &getReceiptCaller{
		api:        c.api,
		actionHash: actionHash,
	}
}

func (c *client) GetLogs(request *iotexapi.GetLogsRequest) GetLogsCaller {
	return &getLogsCaller{
		api:     c.api,
		Request: request,
	}
}

func (c *client) API() iotexapi.APIServiceClient { return c.api }
