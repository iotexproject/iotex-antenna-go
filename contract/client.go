// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package contract

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"

	"go.uber.org/zap"

	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

// ExpectedBalance defines an account-balance pair
type ExpectedBalance struct {
	Account    string `json:"account"`
	RawBalance string `json:"rawBalance"`
}

func (eb *ExpectedBalance) Balance() *big.Int {
	balance, ok := new(big.Int).SetString(eb.RawBalance, 10)
	if !ok {
		log.L().Panic("invalid balance", zap.String("balance", eb.RawBalance))
	}

	return balance
}

type Log struct {
	Topics []string `json:"topics"`
	Data   string   `json:"data"`
}

type ExecutionConfig struct {
	Comment                 string            `json:"comment"`
	ContractIndex           int               `json:"contractIndex"`
	AppendContractAddress   bool              `json:"appendContractAddress"`
	ContractIndexToAppend   int               `json:"contractIndexToAppend"`
	ContractAddressToAppend string            `json:"contractAddressToAppend"`
	ReadOnly                bool              `json:"readOnly"`
	RawPrivateKey           string            `json:"rawPrivateKey"`
	RawByteCode             string            `json:"rawByteCode"`
	RawAmount               string            `json:"rawAmount"`
	RawGasLimit             uint              `json:"rawGasLimit"`
	RawGasPrice             string            `json:"rawGasPrice"`
	Failed                  bool              `json:"failed"`
	RawReturnValue          string            `json:"rawReturnValue"`
	RawExpectedGasConsumed  uint              `json:"rawExpectedGasConsumed"`
	ExpectedBalances        []ExpectedBalance `json:"expectedBalances"`
	ExpectedLogs            []Log             `json:"expectedLogs"`
}

func (cfg *ExecutionConfig) PrivateKey() (priKey keypair.PrivateKey) {
	priKey, err := keypair.HexStringToPrivateKey(cfg.RawPrivateKey)
	if err != nil {
		return
	}
	return
}

func (cfg *ExecutionConfig) Executor() address.Address {
	priKey := cfg.PrivateKey()
	addr, err := address.FromBytes(priKey.PublicKey().Hash())
	if err != nil {
		log.L().Panic(
			"invalid private key",
			zap.String("privateKey", cfg.RawPrivateKey),
			zap.Error(err),
		)
	}

	return addr
}

func (cfg *ExecutionConfig) ByteCode() []byte {
	byteCode, err := hex.DecodeString(cfg.RawByteCode)
	if err != nil {
		log.L().Panic(
			"invalid byte code",
			zap.String("byteCode", cfg.RawByteCode),
			zap.Error(err),
		)
	}
	if cfg.AppendContractAddress {
		addr, err := address.FromString(cfg.ContractAddressToAppend)
		if err != nil {
			log.L().Panic(
				"invalid contract address to append",
				zap.String("contractAddressToAppend", cfg.ContractAddressToAppend),
				zap.Error(err),
			)
		}
		ba := addr.Bytes()
		ba = append(make([]byte, 12), ba...)
		byteCode = append(byteCode, ba...)
	}

	return byteCode
}

func (cfg *ExecutionConfig) Amount() *big.Int {
	amount, ok := new(big.Int).SetString(cfg.RawAmount, 10)
	if !ok {
		log.L().Panic("invalid amount", zap.String("amount", cfg.RawAmount))
	}

	return amount
}

func (cfg *ExecutionConfig) GasPrice() *big.Int {
	price, ok := new(big.Int).SetString(cfg.RawGasPrice, 10)
	if !ok {
		log.L().Panic("invalid gas price", zap.String("gasPrice", cfg.RawGasPrice))
	}

	return price
}

func (cfg *ExecutionConfig) GasLimit() uint64 {
	return uint64(cfg.RawGasLimit)
}

func (cfg *ExecutionConfig) ExpectedGasConsumed() uint64 {
	return uint64(cfg.RawExpectedGasConsumed)
}

func (cfg *ExecutionConfig) ExpectedReturnValue() []byte {
	retval, err := hex.DecodeString(cfg.RawReturnValue)
	if err != nil {
		log.L().Panic(
			"invalid return value",
			zap.String("returnValue", cfg.RawReturnValue),
			zap.Error(err),
		)
	}
	return retval
}

type SmartContract struct {
	// the order matters
	InitBalances      []ExpectedBalance `json:"initBalances"`
	Deployments       []ExecutionConfig `json:"deployments"`
	Executions        []ExecutionConfig `json:"executions"`
	rpc               *rpcmethod.RPCMethod
	deployActionHash  []*hash.Hash256
	contractAddresses []string
}

func NewSmartContract(file, endpoint string) (sct *SmartContract, err error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return
	}
	sctBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}
	sct = &SmartContract{}
	err = json.Unmarshal(sctBytes, sct)
	if err != nil {
		return
	}
	sct.rpc, err = rpcmethod.NewRPCMethod(endpoint)
	return
}

//
func (sct *SmartContract) runExecution(
	ecfg *ExecutionConfig,
	contractAddr string,
) (err error) {
	log.S().Info(ecfg.Comment)
	request := &iotexapi.GetAccountRequest{Address: ecfg.Executor().String()}
	res, err := sct.rpc.GetAccount(request)
	nonce := res.AccountMeta.Nonce
	if err != nil {
		return
	}
	exec, err := action.NewExecution(
		contractAddr,
		nonce+1,
		ecfg.Amount(),
		ecfg.GasLimit(),
		ecfg.GasPrice(),
		ecfg.ByteCode(),
	)
	if err != nil {
		return
	}

	builder := &action.EnvelopeBuilder{}
	elp := builder.SetAction(exec).
		SetNonce(exec.Nonce()).
		SetGasLimit(ecfg.GasLimit()).
		SetGasPrice(ecfg.GasPrice()).
		Build()
	selp, err := action.Sign(elp, ecfg.PrivateKey())
	if err != nil {
		return
	}
	request2 := &iotexapi.SendActionRequest{Action: selp.Proto()}
	_, err = sct.rpc.SendAction(request2)
	if err != nil {
		return
	}
	hash := exec.Hash()
	sct.deployActionHash = append(sct.deployActionHash, &hash)
	return
}

func (sct *SmartContract) DeployContracts() (err error) {
	for _, contract := range sct.Deployments {
		if contract.AppendContractAddress {
			contract.ContractAddressToAppend = sct.contractAddresses[contract.ContractIndexToAppend]
		}
		err = sct.runExecution(&contract, action.EmptyAddress)
		if err != nil {
			return
		}
	}
	return
}
func (sct *SmartContract) GetContractAddresses() []string {
	return sct.contractAddresses
}
func (sct *SmartContract) GetContractAddressFromChain() {
	for _, hash := range sct.deployActionHash {
		hashString := hex.EncodeToString(hash[:])
		request3 := &iotexapi.GetReceiptByActionRequest{ActionHash: hashString}
		res3, err := sct.rpc.GetReceiptByAction(request3)
		if err != nil {
			return
		}
		cd := res3.ReceiptInfo.Receipt.ContractAddress
		sct.contractAddresses = append(sct.contractAddresses, cd)
	}
}

//func (sct *SmartContract) run(r *require.Assertions) {
//	// prepare blockchain
//	ctx := context.Background()
//	bc := sct.prepareBlockchain(ctx, r)
//	defer r.NoError(bc.Stop(ctx))
//
//	// deploy smart contract
//	contractAddresses := sct.deployContracts(bc, r)
//	if len(contractAddresses) == 0 {
//		return
//	}
//
//	// run executions
//	for _, exec := range sct.Executions {
//		contractAddr := contractAddresses[exec.ContractIndex]
//		if exec.AppendContractAddress {
//			exec.ContractAddressToAppend = contractAddresses[exec.ContractIndexToAppend]
//		}
//		retval, receipt, err := runExecution(bc, &exec, contractAddr)
//		r.NoError(err)
//		r.NotNil(receipt)
//		if exec.Failed {
//			r.Equal(action.FailureReceiptStatus, receipt.Status)
//		} else {
//			r.Equal(action.SuccessReceiptStatus, receipt.Status)
//		}
//		if exec.ExpectedGasConsumed() != 0 {
//			r.Equal(exec.ExpectedGasConsumed(), receipt.GasConsumed)
//		}
//		if exec.ReadOnly {
//			expected := exec.ExpectedReturnValue()
//			if len(expected) == 0 {
//				r.Equal(0, len(retval))
//			} else {
//				r.Equal(expected, retval)
//			}
//			return
//		}
//		for _, expectedBalance := range exec.ExpectedBalances {
//			account := expectedBalance.Account
//			if account == "" {
//				account = contractAddr
//			}
//			balance, err := bc.Balance(account)
//			r.NoError(err)
//			r.Equal(
//				0,
//				balance.Cmp(expectedBalance.Balance()),
//				"balance of account %s is different from expectation, %d vs %d",
//				account,
//				balance,
//				expectedBalance.Balance(),
//			)
//		}
//		r.Equal(len(exec.ExpectedLogs), len(receipt.Logs))
//		// TODO: check value of logs
//	}
//}
