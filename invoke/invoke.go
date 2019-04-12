package invoke

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/address"
	"github.com/iotexproject/iotex-core/protogen/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/rpcmethod"
)

// Errors
var (
	ErrNoAbiName   = errors.New("No such ABI name")
	ErrInvalidABI  = errors.New("Invilid ABI")
	ErrInvalidAddr = errors.New("Invilid address")
)

var (
	getRPCMethod      getRPCMethodFn
	getAbstractAction getAbstractActionFn
	getAmount         getAmountFn
	sign              signFn
)

type (
	getRPCMethodFn      func() (*rpcmethod.RPCMethod, error)
	getAbstractActionFn func() (*action.AbstractAction, error)
	getAmountFn         func() (*big.Int, error)
	signFn              func(action.Envelope) (*action.SealedEnvelope, error)
)

// ABI provides simple interface for ABI
type ABI struct {
	ctt     *Contract
	abiName string
}

// Contract provides simple interface for contract
type Contract struct {
	abi     abi.ABI
	address string
}

// NewContract generate a new contract by abi JSON
func NewContract(abiJSON string, addr string) (ctt *Contract, err error) {
	reader := strings.NewReader(abiJSON)
	ctt.abi, err = abi.JSON(reader)
	if err != nil {
		return nil, ErrInvalidABI
	}
	if _, err := address.FromString(addr); err != nil {
		return nil, ErrInvalidAddr
	}
	ctt.address = addr
	return
}

// ABI returns an ABI by abi name
func (ctt *Contract) ABI(abiName string) (abi *ABI, err error) {
	for _, method := range ctt.abi.Methods {
		if method.Name == abiName {
			abi = &ABI{ctt: ctt, abiName: abiName}
			return
		}
	}
	err = ErrNoAbiName
	return
}

// Call send an invoke tx to blockchain
func (abi *ABI) Call(args ...interface{}) (err error) {
	bytecode, err := abi.ctt.abi.Pack(abi.abiName, args)
	if err != nil {
		return
	}
	abs, err := getAbstractAction()
	if err != nil {
		return
	}
	amount, err := getAmount()
	if err != nil {
		return
	}
	tx, err := action.NewExecution(abi.ctt.address, abs.Nonce(), amount, abs.GasLimit(), abs.GasPrice(), bytecode)
	if err != nil {
		return
	}
	bd := &action.EnvelopeBuilder{}
	elp := bd.SetNonce(abs.Nonce()).
		SetGasPrice(abs.GasPrice()).
		SetGasLimit(abs.GasLimit()).
		SetAction(tx).Build()
	sealed, err := sign(elp)
	if err != nil {
		return
	}
	rpcMethod, err := getRPCMethod()
	if err != nil {
		return
	}
	request := &iotexapi.SendActionRequest{Action: sealed.Proto()}
	_, err = rpcMethod.SendAction(request)
	if err != nil {
		return
	}
	return nil
}
