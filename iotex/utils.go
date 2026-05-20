package iotex

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

func sign(a account.Account, act *iotextypes.ActionCore) (*iotextypes.Action, error) {
	msg, err := proto.Marshal(act)
	if err != nil {
		return nil, err
	}
	sig, err := a.Sign(msg)
	if err != nil {
		return nil, err
	}
	return &iotextypes.Action{
		Core:         act,
		SenderPubKey: a.PublicKey().Bytes(),
		Signature:    sig,
	}, nil
}

// signContainer builds, signs and packages an Ethereum typed transaction
// (AccessList/DynamicFee/Blob/SetCode) into an iotex Action using the
// TX_CONTAINER encoding: the signed tx is RLP-marshaled and carried verbatim
// inside iotextypes.TxContainer, so the node decodes it with go-ethereum
// without reconstructing it from proto fields.
func signContainer(a account.Account, act *iotextypes.ActionCore) (*iotextypes.Action, error) {
	tx, err := actionToRLP(act)
	if err != nil {
		return nil, err
	}
	signer := types.LatestSignerForChainID(big.NewInt(int64(act.GetChainID())))
	h := signer.Hash(tx)
	sig, err := a.PrivateKey().Sign(h[:])
	if err != nil {
		return nil, err
	}
	// WithSignature wants the signature in recovery-id form (V = 0/1).
	sc := make([]byte, 65)
	copy(sc, sig)
	if sc[64] >= 27 {
		sc[64] -= 27
	}
	signedTx, err := tx.WithSignature(signer, sc)
	if err != nil {
		return nil, err
	}
	raw, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// The Action signature must match what the node extracts from the raw tx.
	actionSig, err := extractEthTxSig(signedTx)
	if err != nil {
		return nil, err
	}
	return &iotextypes.Action{
		Core: &iotextypes.ActionCore{
			ChainID: act.GetChainID(),
			Action: &iotextypes.ActionCore_TxContainer{
				TxContainer: &iotextypes.TxContainer{Raw: raw},
			},
		},
		SenderPubKey: a.PublicKey().Bytes(),
		Signature:    actionSig,
		Encoding:     iotextypes.Encoding_TX_CONTAINER,
	}, nil
}

// ActionHash computes the hash of an action
func ActionHash(act *iotextypes.Action, chainid uint32) (hash.Hash256, error) {
	switch act.Encoding {
	case iotextypes.Encoding_IOTEX_PROTOBUF:
		ser, err := proto.Marshal(act)
		if err != nil {
			return hash.ZeroHash256, err
		}
		return hash.Hash256b(ser), nil
	case iotextypes.Encoding_ETHEREUM_RLP:
		tx, err := actionToRLP(act.Core)
		if err != nil {
			return hash.ZeroHash256, err
		}
		h, err := rlpSignedHash(tx, chainid, act.GetSignature())
		if err != nil {
			return hash.ZeroHash256, err
		}
		return h, nil
	case iotextypes.Encoding_TX_CONTAINER:
		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(act.GetCore().GetTxContainer().GetRaw()); err != nil {
			return hash.ZeroHash256, err
		}
		h := tx.Hash()
		return hash.BytesToHash256(h[:]), nil
	default:
		return hash.ZeroHash256, fmt.Errorf("invalid encoding type = %v", act.Encoding)
	}
}
