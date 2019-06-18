package iotex

import (
	"github.com/gogo/protobuf/proto"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
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
