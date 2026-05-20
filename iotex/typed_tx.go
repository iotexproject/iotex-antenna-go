// Copyright (c) 2026 IoTeX
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

// BlobData carries the EIP-4844 blob fields a caller supplies via
// ExecuteContractCaller.SetBlobTxData. The Sidecar may be nil only when the
// caller already has the network broadcasting the blobs separately.
type BlobData struct {
	BlobFeeCap *big.Int
	BlobHashes []common.Hash
	Sidecar    *types.BlobTxSidecar
}

func (b *BlobData) toProto() *iotextypes.BlobTxData {
	if b == nil {
		return nil
	}
	out := &iotextypes.BlobTxData{}
	if b.BlobFeeCap != nil {
		out.BlobFeeCap = b.BlobFeeCap.String()
	}
	if n := len(b.BlobHashes); n > 0 {
		out.BlobHashes = make([][]byte, n)
		for i, h := range b.BlobHashes {
			out.BlobHashes[i] = h.Bytes()
		}
	}
	if b.Sidecar != nil {
		out.BlobTxSidecar = sidecarToProto(b.Sidecar)
	}
	return out
}

func sidecarToProto(sc *types.BlobTxSidecar) *iotextypes.BlobTxSidecar {
	if sc == nil || len(sc.Blobs) == 0 {
		return nil
	}
	out := &iotextypes.BlobTxSidecar{
		Blobs:       make([][]byte, len(sc.Blobs)),
		Commitments: make([][]byte, len(sc.Commitments)),
		Proofs:      make([][]byte, len(sc.Proofs)),
	}
	for i := range sc.Blobs {
		out.Blobs[i] = sc.Blobs[i][:]
	}
	for i := range sc.Commitments {
		out.Commitments[i] = sc.Commitments[i][:]
	}
	for i := range sc.Proofs {
		out.Proofs[i] = sc.Proofs[i][:]
	}
	return out
}

// authListToProto converts go-ethereum SetCodeAuthorization tuples to the
// proto form. The proto ChainID field is uint32, so an authorization whose
// ChainID does not fit is rejected rather than silently truncated.
func authListToProto(list []types.SetCodeAuthorization) ([]*iotextypes.SetCodeAuthorization, error) {
	if len(list) == 0 {
		return nil, nil
	}
	out := make([]*iotextypes.SetCodeAuthorization, len(list))
	for i, a := range list {
		if !a.ChainID.IsUint64() || a.ChainID.Uint64() > math.MaxUint32 {
			return nil, fmt.Errorf("setcode authorization %d: chainID %s overflows uint32", i, a.ChainID.String())
		}
		out[i] = &iotextypes.SetCodeAuthorization{
			ChainID: uint32(a.ChainID.Uint64()),
			Address: a.Address.Bytes(),
			Nonce:   a.Nonce,
			V:       uint64(a.V),
			R:       a.R.Bytes(),
			S:       a.S.Bytes(),
		}
	}
	return out, nil
}
