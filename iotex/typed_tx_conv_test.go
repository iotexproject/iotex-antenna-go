// Copyright (c) 2026 IoTeX
// This source code is provided 'as is' and no warranties are given as to title or non-infringement, merchantability
// or fitness for purpose and, to the extent permitted by law, all liability for your use of the code is disclaimed.
// This source code is governed by Apache License 2.0 that can be found in the LICENSE file.

package iotex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBlobData_toProto(t *testing.T) {
	hashA := common.HexToHash("0x01aa")
	hashB := common.HexToHash("0x02bb")
	b := &BlobData{
		BlobFeeCap: big.NewInt(987654321),
		BlobHashes: []common.Hash{hashA, hashB},
		Sidecar:    nil,
	}
	pb := b.toProto()
	require.Equal(t, "987654321", pb.GetBlobFeeCap())
	require.Equal(t, [][]byte{hashA.Bytes(), hashB.Bytes()}, pb.GetBlobHashes())
	require.Nil(t, pb.GetBlobTxSidecar())
}

func TestBlobData_toProto_NilSafe(t *testing.T) {
	require.Nil(t, (*BlobData)(nil).toProto())
}

func TestAuthListToProto(t *testing.T) {
	addr := common.HexToAddress("0xababababababababababababababababababab01")
	auths := []SetCodeAuthorization{{
		ChainID: 2, // IoTeX testnet → maps to EVM 4690
		Address: addr,
		Nonce:   7,
		V:       1,
		R:       new(big.Int).SetBytes([]byte{0xaa, 0xbb}),
		S:       new(big.Int).SetBytes([]byte{0xcc, 0xdd}),
	}}
	pb, err := authListToProto(auths)
	require.NoError(t, err)
	require.Len(t, pb, 1)
	require.Equal(t, uint32(4690), pb[0].GetChainID())
	require.Equal(t, addr.Bytes(), pb[0].GetAddress())
	require.Equal(t, uint64(7), pb[0].GetNonce())
	require.Equal(t, uint64(1), pb[0].GetV())
	require.Equal(t, []byte{0xaa, 0xbb}, pb[0].GetR())
	require.Equal(t, []byte{0xcc, 0xdd}, pb[0].GetS())
}

func TestAuthListToProto_Empty(t *testing.T) {
	pb, err := authListToProto(nil)
	require.NoError(t, err)
	require.Nil(t, pb)
	pb, err = authListToProto([]SetCodeAuthorization{})
	require.NoError(t, err)
	require.Nil(t, pb)
}

func TestAuthListToProto_InvalidChainID(t *testing.T) {
	auths := []SetCodeAuthorization{{
		ChainID: 4, // invalid IoTeX chain ID (must be 1, 2, or 3)
		Address: common.HexToAddress("0x01"),
	}}
	_, err := authListToProto(auths)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid chain id")
}
