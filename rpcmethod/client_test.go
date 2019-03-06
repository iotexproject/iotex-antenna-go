package rpcmethod

import (
	"testing"

	"github.com/iotexproject/iotex-core/protogen/iotexapi"
)

func Test_Client(t *testing.T) {
	r, err := NewRPCMethod("128.199.211.107:14014")
	if err != nil {
		t.Fatal("rpc conn", err)
	}
	actionReq := iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{
			ByHash: &iotexapi.GetActionByHashRequest{
				ActionHash: "4ab56a16a0a904c92cb05afcdcdd962e4535b5f318dfc2d63684a02727d5b5da",
			},
		},
	}
	actionResp, err := r.GetActions(&actionReq)
	if err != nil {
		t.Fatal("get actions", err)
	}
	t.Logf("%+v\n", actionResp)

	accReq := iotexapi.GetAccountRequest{
		Address: "io14maafjgdxazyqluwl9yur85rfd6kn59l9zvc57",
	}
	accResp, err := r.GetAccount(&accReq)
	if err != nil {
		t.Fatal("get account", err)
	}
	t.Logf("%+v\n", accResp)
}
