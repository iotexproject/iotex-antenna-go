package main

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
)

const (
	createDID  = "createDID"
	deleteDID  = "deleteDID"
	updateHash = "updateHash"
	updateURI  = "updateURI"
	getHash    = "getHash"
	getURI     = "getURI"
)

type DID interface {
	CreateDID(did, didHash, url string) (hash string, err error)
	DeleteDID(did string) (hash string, err error)
	UpdateHash(did, didHash string) (hash string, err error)
	UpdateUri(did, uri string) (hash string, err error)
	GetHash(did string) (hash string, err error)
	GetUri(did string) (uri string, err error)
}

type did struct {
	endpoint string
	account  account.Account
	contract address.Address
	abi      abi.ABI
	gasPrice *big.Int
	gasLimit uint64
}

func NewDID(endpoint, privateKey, contract, abiString string, gasPrice *big.Int, gasLimit uint64) (d DID, err error) {
	abi, err := abi.JSON(strings.NewReader(abiString)) // note,this is IoTeXDID_abi
	if err != nil {
		return
	}
	var acc account.Account
	if privateKey != "" {
		acc, err = account.HexStringToAccount(privateKey)
		if err != nil {
			return
		}
	}

	addr, err := address.FromString(contract)
	if err != nil {
		return
	}
	d = &did{endpoint, acc, addr, abi, gasPrice, gasLimit}
	return
}

func (d *did) CreateDID(id, didHash, url string) (hash string, err error) {
	if len(didHash) != 64 {
		err = errors.New("hash should be 32 bytes")
		return
	}
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), d.account)
	hashSlice, err := hex.DecodeString(didHash)
	if err != nil {
		return
	}
	var hashArray [32]byte
	copy(hashArray[:], hashSlice)
	h, err := cli.Contract(d.contract, d.abi).Execute(createDID, id, hashArray, url).SetGasPrice(d.gasPrice).SetGasLimit(d.gasLimit).Call(context.Background())
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}

func (d *did) DeleteDID(did string) (hash string, err error) {
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), d.account)
	h, err := cli.Contract(d.contract, d.abi).Execute(deleteDID, did).SetGasPrice(d.gasPrice).SetGasLimit(d.gasLimit).Call(context.Background())
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}

func (d *did) UpdateHash(did, didHash string) (hash string, err error) {
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), d.account)
	hashSlice, err := hex.DecodeString(didHash)
	if err != nil {
		return
	}
	var hashArray [32]byte
	copy(hashArray[:], hashSlice)
	h, err := cli.Contract(d.contract, d.abi).Execute(updateHash, did, hashArray).SetGasPrice(d.gasPrice).SetGasLimit(d.gasLimit).Call(context.Background())
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}

func (d *did) UpdateUri(did, uri string) (hash string, err error) {
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(conn), d.account)
	h, err := cli.Contract(d.contract, d.abi).Execute(updateURI, did, uri).SetGasPrice(d.gasPrice).SetGasLimit(d.gasLimit).Call(context.Background())
	if err != nil {
		return
	}
	hash = hex.EncodeToString(h[:])
	return
}

func (d *did) GetHash(did string) (hash string, err error) {
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	ret, err := cli.ReadOnlyContract(d.contract, d.abi).Read(getHash, did).Call(context.Background())
	if err != nil {
		return
	}
	hashBytes := [32]uint8{}
	err = ret.Unmarshal(&hashBytes)
	if err != nil {
		return
	}
	hash = hex.EncodeToString(hashBytes[:])
	return
}

func (d *did) GetUri(did string) (uri string, err error) {
	conn, err := iotex.NewDefaultGRPCConn(d.endpoint)
	if err != nil {
		return
	}
	defer conn.Close()
	cli := iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(conn))
	ret, err := cli.ReadOnlyContract(d.contract, d.abi).Read(getURI, did).Call(context.Background())
	if err != nil {
		return
	}
	err = ret.Unmarshal(&uri)
	if err != nil {
		return
	}
	return
}
