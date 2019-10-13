package wait

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Wait waits on a send action caller to finish (wait 20 second), then fetch the receipt of the action,
// to make sure the action is on chain.
func Wait(ctx context.Context, caller iotex.SendActionCaller, opts ...grpc.CallOption) error {
	h, err := caller.Call(ctx, opts...)
	if err != nil {
		return err
	}

	time.Sleep(25 * time.Second)

	var response *iotexapi.GetReceiptByActionResponse
	rerr := backoff.Retry(func() error {
		response, err = caller.API().GetReceiptByAction(ctx, &iotexapi.GetReceiptByActionRequest{
			ActionHash: hex.EncodeToString(h[:]),
		}, opts...)
		return err
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(30*time.Second), 3))
	if rerr != nil {
		return errors.Errorf("execution get receipt timed out: %v", rerr)
	}
	if response.ReceiptInfo.Receipt.Status != 1 {
		return errors.Errorf("execution failed: %x", h)
	}
	return nil
}
