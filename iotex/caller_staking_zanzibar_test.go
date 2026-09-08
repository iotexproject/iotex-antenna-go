package iotex

import (
	"math/big"
	"testing"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// buildCore runs the same action->envelope mapping Call() does, without a
// network. The mapping is the part that can silently go wrong: a wrapper that
// builds the right message but files it under the wrong oneof field produces a
// well-formed action the node rejects.
func buildCore(t *testing.T, set func(c *stakingCaller)) *iotextypes.ActionCore {
	t.Helper()
	c := &stakingCaller{}
	set(c)
	core := &iotextypes.ActionCore{Version: ProtocolVersion}
	switch a := c.action.(type) {
	case *iotextypes.SetVoterRewardOptIn:
		core.Action = &iotextypes.ActionCore_SetVoterRewardOptIn{SetVoterRewardOptIn: a}
	case *iotextypes.SetVoterRewardDestination:
		core.Action = &iotextypes.ActionCore_SetVoterRewardDestination{SetVoterRewardDestination: a}
	case *iotextypes.CandidateRegister:
		core.Action = &iotextypes.ActionCore_CandidateRegister{CandidateRegister: a}
	case *iotextypes.CandidateBasicInfo:
		core.Action = &iotextypes.ActionCore_CandidateUpdate{CandidateUpdate: a}
	default:
		t.Fatalf("unmapped action type %T", c.action)
	}
	return core
}

// fieldNumber reports which oneof field the core actually populated, read via
// protobuf reflection rather than by walking the wire format -- field numbers
// above 15 use a multi-byte key, and a hand-rolled walker that assumes one byte
// silently misreads exactly the range these actions live in.
func fieldNumber(t *testing.T, core *iotextypes.ActionCore) int {
	t.Helper()
	m := core.ProtoReflect()
	od := m.Descriptor().Oneofs().ByName("action")
	require.NotNil(t, od, "ActionCore has no `action` oneof")
	fd := m.WhichOneof(od)
	require.NotNil(t, fd, "no action payload set")
	return int(fd.Number())
}

// The field numbers the node dispatches on. Pinned here because a wrapper that
// files an action under the wrong one still marshals, still signs, and is only
// rejected once it reaches a node.
func TestZanzibarActionsUseTheProtocolFieldNumbers(t *testing.T) {
	r := require.New(t)
	recipient, err := address.FromString("io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02")
	r.NoError(err)

	r.Equal(57, fieldNumber(t, buildCore(t, func(c *stakingCaller) { c.SetVoterRewardOptIn() })))
	r.Equal(58, fieldNumber(t, buildCore(t, func(c *stakingCaller) { c.SetVoterRewardDestination(recipient) })))
}

// SetVoterRewardOptIn carries no payload at all -- the sender is the delegate.
// An empty message still has to occupy its field, or the node sees no action.
func TestOptInIsAnEmptyButPresentPayload(t *testing.T) {
	core := buildCore(t, func(c *stakingCaller) { c.SetVoterRewardOptIn() })
	require.NotNil(t, core.GetSetVoterRewardOptIn())
	raw, err := proto.Marshal(core)
	require.NoError(t, err)
	require.NotEmpty(t, raw)
}

func TestSetVoterRewardDestinationCarriesRawAddressBytes(t *testing.T) {
	r := require.New(t)
	recipient, err := address.FromString("io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02")
	r.NoError(err)
	core := buildCore(t, func(c *stakingCaller) { c.SetVoterRewardDestination(recipient) })
	// 20 raw bytes, not the bech32 string: the protocol reads them with
	// address.FromBytes.
	r.Equal(recipient.Bytes(), core.GetSetVoterRewardDestination().GetRecipient())
	r.Len(core.GetSetVoterRewardDestination().GetRecipient(), 20)
}

// WithBLS has to reach both Register and Update, and the proof must travel
// with the key -- core verifies the pair against the owner address and fails
// the action when only one is present.
func TestWithBLSReachesRegisterAndUpdate(t *testing.T) {
	r := require.New(t)
	owner, err := address.FromString("io1k7rqg4ksjx93f467mdfjzm6un6ezskwcppls80")
	r.NoError(err)
	pub, pop := []byte("pubkey-48-bytes-placeholder"), []byte("pop-96-bytes-placeholder")

	reg := buildCore(t, func(c *stakingCaller) {
		c.WithBLS(pub, pop)
		c.Register("cand", owner, owner, owner, bigOne(), 91, true, nil)
	})
	r.Equal(pub, reg.GetCandidateRegister().GetCandidate().GetBlsPubKey())
	r.Equal(pop, reg.GetCandidateRegister().GetCandidate().GetBlsPop())

	upd := buildCore(t, func(c *stakingCaller) {
		c.WithBLS(pub, pop)
		c.Update("cand", owner, owner)
	})
	r.Equal(pub, upd.GetCandidateUpdate().GetBlsPubKey())
	r.Equal(pop, upd.GetCandidateUpdate().GetBlsPop())
}

// Without WithBLS the fields stay absent rather than becoming empty bytes:
// a pre-Zanzibar registration must marshal exactly as it did before.
func TestRegisterWithoutBLSLeavesTheFieldsUnset(t *testing.T) {
	r := require.New(t)
	owner, err := address.FromString("io1k7rqg4ksjx93f467mdfjzm6un6ezskwcppls80")
	r.NoError(err)
	reg := buildCore(t, func(c *stakingCaller) {
		c.Register("cand", owner, owner, owner, bigOne(), 91, true, nil)
	})
	r.Nil(reg.GetCandidateRegister().GetCandidate().GetBlsPubKey())
	r.Nil(reg.GetCandidateRegister().GetCandidate().GetBlsPop())
}

func bigOne() *big.Int { return big.NewInt(1) }
