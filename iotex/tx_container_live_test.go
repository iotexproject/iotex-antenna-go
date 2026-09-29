package iotex

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

// TestTxContainerLive sends a typed (dynamic-fee) transfer to a running chain.
//
// signContainer stopped rebuilding a bare ActionCore and now reuses the
// caller's, so the wire message carries nonce, gasLimit, gasTipCap and
// gasFeeCap alongside the container. TestCall_TypedTx_ActionCoreFieldsPreserved
// asserts that client-side; it cannot say whether a node accepts the result.
// A node reads only ChainID and TxContainer.Raw out of that message and
// ignores the rest, so the extra fields are inert -- but "inert" is a claim
// about the node, and only a node can settle it.
func TestTxContainerLive(t *testing.T) {
	endpoint := os.Getenv("IOTEX_LIVE_GRPC")
	key := os.Getenv("IOTEX_LIVE_KEY")
	if endpoint == "" || key == "" {
		t.Skip("set IOTEX_LIVE_GRPC and IOTEX_LIVE_KEY to run")
	}
	r := require.New(t)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	r.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(key)
	r.NoError(err)
	api := iotexapi.NewAPIServiceClient(conn)
	cli := NewAuthedClient(api, 3, acc)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	to, err := address.FromString("io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02")
	r.NoError(err)
	amount := big.NewInt(1234567890)

	h, err := cli.Transfer(to, amount).
		SetTxType(dynamicFeeTxType).
		SetGasTipCap(big.NewInt(1000000000000)).
		SetGasFeeCap(big.NewInt(2000000000000)).
		SetGasLimit(100000).
		Call(ctx)
	r.NoError(err)
	t.Logf("typed transfer  %x", h)
	r.EqualValues(1, waitStatus(ctx, t, cli, h))

	// Read it back. TX_CONTAINER is a transport encoding and does not survive:
	// Unfold() rebuilds the real envelope and rewrites the encoding through
	// typeToEncoding(), which reports ETHEREUM_EIP155 for every typed tx. What
	// has to survive is the content -- if the extra ActionCore fields had
	// disturbed the node's read of the container, it would be this that came
	// back wrong.
	resp, err := api.GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{ActionHash: hexHash(h)},
		},
	})
	r.NoError(err)
	r.Len(resp.GetActionInfo(), 1)
	act := resp.GetActionInfo()[0].GetAction()
	r.Equal(iotextypes.Encoding_ETHEREUM_EIP155, act.GetEncoding())

	core := act.GetCore()
	tr := core.GetTransfer()
	r.NotNil(tr, "node must unfold the container into a transfer")
	r.Equal(amount.String(), tr.GetAmount())
	r.Equal(to.String(), tr.GetRecipient())
	// The typed fields live inside the signed raw tx, so their surviving the
	// round trip is what says the container was read rather than guessed at.
	r.EqualValues(dynamicFeeTxType, core.GetTxType())
	r.Equal("1000000000000", core.GetGasTipCap())
	r.Equal("2000000000000", core.GetGasFeeCap())
	r.EqualValues(100000, core.GetGasLimit())
	t.Logf("  encoding=%s txType=%d tip=%s cap=%s gasLimit=%d -> %s amount=%s",
		act.GetEncoding(), core.GetTxType(), core.GetGasTipCap(), core.GetGasFeeCap(),
		core.GetGasLimit(), tr.GetRecipient(), tr.GetAmount())
}

func hexHash(h [32]byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, 64)
	for i, b := range h {
		out[i*2] = hexdigits[b>>4]
		out[i*2+1] = hexdigits[b&0x0f]
	}
	return string(out)
}
