package iotex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

type authedClient struct {
	client
	chainID uint32
	account account.Account
}

// NewAuthedClient creates an AuthedClient using given account's credentials.
func NewAuthedClient(api iotexapi.APIServiceClient, chainID uint32, a account.Account) AuthedClient {
	return &authedClient{
		client: client{
			api: api,
		},
		chainID: chainID,
		account: a,
	}
}

func (c *authedClient) Contract(co address.Address, abi abi.ABI) Contract {
	return &contract{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
		},
		address: co,
		abi:     &abi,
	}
}

func (c *authedClient) Transfer(to address.Address, value *big.Int) SendActionCaller {
	return &transferCaller{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
		},
		amount:    value,
		recipient: to,
	}
}

func (c *authedClient) ClaimReward(value *big.Int) ClaimRewardCaller {
	return &claimRewardCaller{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
		},
		amount: value,
	}
}

func (c *authedClient) DeployContract(data []byte) DeployContractCaller {
	return &deployContractCaller{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
			payload: data,
		},
	}
}

//Staking interface
func (c *authedClient) Staking() StakingCaller {
	return &stakingCaller{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
		}}
}

//Candidate interface
func (c *authedClient) Candidate() CandidateCaller {
	return &stakingCaller{
		sendActionCaller: &sendActionCaller{
			chainID: c.chainID,
			account: c.account,
			api:     c.api,
		}}
}

func (c *authedClient) Account() account.Account { return c.account }

func (c *authedClient) ChainID() uint32 { return c.chainID }

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
