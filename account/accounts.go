package account

import (
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/ethereum/go-ethereum/crypto"
)

type Accounts struct {
	acts []Account
}

func (acts *Accounts) Create() (Account, error) {
	private, err := crypto.GenerateKey()
	if err != nil {
		return Account{}, err
	}
	return privateToAccount(private)
}

func (acts *Accounts) PrivateKeyToAccount(privateKey string) (Account, error) {
	private, err := keypair.DecodePrivateKey(privateKey)
	if err != nil {
		return Account{}, nil
	}

	return privateToAccount(private)
}

func (acts *Accounts) Sign(data []byte, privateKey string) ([]byte, error) {
	act, err := acts.PrivateKeyToAccount(privateKey)
	if err != nil {
		return nil, err
	}

	return act.Sign(data)
}
