// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/iotex-antenna-go/account"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

type (
	// IotexAction iotextypes Action
	IotexAction struct {
		*iotextypes.Action
	}

	// IotexActionCore iotextypes ActionCore
	IotexActionCore struct {
		*iotextypes.ActionCore
	}
)

// Sign signs the ActionCore
func (ac *IotexActionCore) Sign(act account.Account) (*IotexAction, error) {
	msg, err := proto.Marshal(ac.ActionCore)
	if err != nil {
		return nil, err
	}
	sig, err := act.Sign(msg)
	if err != nil {
		return nil, err
	}
	return &IotexAction{
		Action: &iotextypes.Action{
			Core:         ac.ActionCore,
			SenderPubKey: act.PublicKey().Bytes(),
			Signature:    sig,
		},
	}, nil
}
