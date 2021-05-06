package iotex

import (
	"fmt"

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
	default:
		return hash.ZeroHash256, fmt.Errorf("invalid encoding type = %v", act.Encoding)
	}
}
