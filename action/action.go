// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/keypair"
	"github.com/iotexproject/iotex-core/protogen/iotextypes"
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

// Hash returns the ActionCore's hash
func (ac *IotexActionCore) Hash() ([]byte, error) {
	msg, err := proto.Marshal(ac.ActionCore)
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(msg)
	return h[:], nil
}

// Sign signs the ActionCore
func (ac *IotexActionCore) Sign(sk keypair.PrivateKey) (*IotexAction, error) {
	h, err := ac.Hash()
	if err != nil {
		return nil, err
	}
	sig, err := sk.Sign(h[:])
	if err != nil {
		return nil, err
	}
	return &IotexAction{
		Action: &iotextypes.Action{
			Core:         ac.ActionCore,
			SenderPubKey: sk.PublicKey().Bytes(),
			Signature:    sig,
		},
	}, nil
}

// Hash returns the Action's hash
func (a *IotexAction) Hash() ([]byte, error) {
	msg, err := proto.Marshal(a.Action)
	if err != nil {
		return nil, err
	}
	h := hash.Hash256b(msg)
	return h[:], nil
}
