package iotex

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

// TestZanzibarBLSLive sends candidateUpdate carrying a BLS public key and a
// proof of possession to a running chain, once with a proof bound to the
// candidate and once with a proof bound to somebody else.
//
// The negative case is the point. The unit tests prove the two byte slices
// reach the right proto fields, but a wrapper that swapped them, truncated one,
// or sent the same bytes twice would still satisfy them -- and the chain would
// reject every such action identically. Only the pair of results separates
// "the chain accepted this" from "the chain verified this": a correctly bound
// proof settles with status 1, and one bound to another candidate settles with
// ErrUnauthorizedOperator (203), which is what VerifyBLSPop surfaces.
//
// Generate each pair with local-dev/iip59-harness/bin/blspop <candidate-io-addr>.
func TestZanzibarBLSLive(t *testing.T) {
	endpoint := os.Getenv("IOTEX_LIVE_GRPC")
	key := os.Getenv("IOTEX_LIVE_KEY")
	if endpoint == "" || key == "" || os.Getenv("IOTEX_LIVE_BLS_PUB") == "" {
		t.Skip("set IOTEX_LIVE_GRPC, IOTEX_LIVE_KEY and the IOTEX_LIVE_BLS_* pairs to run")
	}
	r := require.New(t)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	r.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(key)
	r.NoError(err)
	cli := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), 3, acc)

	// candidateUpdate rewrites name, operator and reward address wholesale, so
	// each has to be re-sent at its current value or the update silently
	// changes it.
	name := os.Getenv("IOTEX_LIVE_CAND_NAME")
	operator, err := address.FromString(os.Getenv("IOTEX_LIVE_CAND_OPERATOR"))
	r.NoError(err)
	reward, err := address.FromString(os.Getenv("IOTEX_LIVE_CAND_REWARD"))
	r.NoError(err)

	cases := []struct {
		name   string
		pub    string
		pop    string
		status uint64
	}{
		{"bound to this candidate", os.Getenv("IOTEX_LIVE_BLS_PUB"), os.Getenv("IOTEX_LIVE_BLS_POP"), 1},
		{"bound to another candidate", os.Getenv("IOTEX_LIVE_BLS_PUB_WRONG"), os.Getenv("IOTEX_LIVE_BLS_POP_WRONG"), 203},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.pub == "" || c.pop == "" {
				t.Skip("key pair not supplied")
			}
			r := require.New(t)
			pub, err := hex.DecodeString(strings.TrimPrefix(c.pub, "0x"))
			r.NoError(err)
			pop, err := hex.DecodeString(strings.TrimPrefix(c.pop, "0x"))
			r.NoError(err)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			h, err := cli.Candidate().WithBLS(pub, pop).
				Update(name, operator, reward).
				SetGasPrice(big.NewInt(1000000000000)).SetGasLimit(1000000).Call(ctx)
			r.NoError(err)
			t.Logf("candidateUpdate+BLS  %x", h)
			r.EqualValues(c.status, waitStatus(ctx, t, cli, h))
		})
	}
}

func waitStatus(ctx context.Context, t *testing.T, cli AuthedClient, h [32]byte) uint64 {
	t.Helper()
	for i := 0; i < 40; i++ {
		time.Sleep(2 * time.Second)
		resp, err := cli.GetReceipt(h).Call(ctx)
		if err != nil {
			continue
		}
		rec := resp.GetReceiptInfo().GetReceipt()
		t.Logf("  status=%d height=%d gas=%d", rec.GetStatus(), rec.GetBlkHeight(), rec.GetGasConsumed())
		return rec.GetStatus()
	}
	t.Fatal("no receipt")
	return 0
}
