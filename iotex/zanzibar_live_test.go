package iotex

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

// TestZanzibarLive submits both IIP-59 actions to a running chain.
//
//	IOTEX_LIVE_GRPC=127.0.0.1:14125 \
//	IOTEX_LIVE_KEY=<hex> IOTEX_LIVE_CHAINID=3 \
//	go test ./iotex -run TestZanzibarLive -v
func TestZanzibarLive(t *testing.T) {
	endpoint := os.Getenv("IOTEX_LIVE_GRPC")
	key := os.Getenv("IOTEX_LIVE_KEY")
	if endpoint == "" || key == "" {
		t.Skip("set IOTEX_LIVE_GRPC and IOTEX_LIVE_KEY to run")
	}
	chainID := uint32(3)
	r := require.New(t)

	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	r.NoError(err)
	defer conn.Close()

	acc, err := account.HexStringToAccount(key)
	r.NoError(err)
	t.Log("sender:", acc.Address().String())

	cli := NewAuthedClient(iotexapi.NewAPIServiceClient(conn), chainID, acc)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	gas := big.NewInt(1000000000000)

	h1, err := cli.Candidate().SetVoterRewardOptIn().SetGasPrice(gas).SetGasLimit(400000).Call(ctx)
	r.NoError(err)
	t.Logf("setVoterRewardOptIn  %x", h1)
	waitOK(ctx, t, cli, h1)

	dest, err := address.FromString("io1ph0u2psnd7muq5xv9623rmxdsxc4uapxhzpg02")
	r.NoError(err)
	h2, err := cli.Candidate().SetVoterRewardDestination(dest).SetGasPrice(gas).SetGasLimit(400000).Call(ctx)
	r.NoError(err)
	t.Logf("setVoterRewardDestination  %x", h2)
	waitOK(ctx, t, cli, h2)
}

func waitOK(ctx context.Context, t *testing.T, cli AuthedClient, h [32]byte) {
	t.Helper()
	for i := 0; i < 40; i++ {
		time.Sleep(2 * time.Second)
		resp, err := cli.GetReceipt(h).Call(ctx)
		if err != nil {
			continue
		}
		rec := resp.GetReceiptInfo().GetReceipt()
		t.Logf("  status=%d height=%d gas=%d", rec.GetStatus(), rec.GetBlkHeight(), rec.GetGasConsumed())
		require.EqualValues(t, 1, rec.GetStatus(), "receipt not successful")
		return
	}
	t.Fatal("no receipt")
}
