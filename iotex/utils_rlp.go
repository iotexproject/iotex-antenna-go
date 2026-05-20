package iotex

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

const (
	legacyTxType     = 0
	accessListTxType = 1
	dynamicFeeTxType = 2
	blobTxType       = 3
	setCodeTxType    = 4
)

func actionToRLP(core *iotextypes.ActionCore) (*types.Transaction, error) {
	to, value, payload, err := extractBody(core)
	if err != nil {
		return nil, err
	}

	gasPrice := new(big.Int)
	gasPrice.SetString(core.GetGasPrice(), 10)

	switch core.GetTxType() {
	case legacyTxType:
		if to == nil {
			return types.NewContractCreation(core.GetNonce(), value, core.GetGasLimit(), gasPrice, payload), nil
		}
		return types.NewTransaction(core.GetNonce(), *to, value, core.GetGasLimit(), gasPrice, payload), nil
	case accessListTxType:
		return types.NewTx(&types.AccessListTx{
			ChainID:    big.NewInt(int64(core.GetChainID())),
			Nonce:      core.GetNonce(),
			GasPrice:   gasPrice,
			Gas:        core.GetGasLimit(),
			To:         to,
			Value:      value,
			Data:       payload,
			AccessList: accessListFromProto(core.GetAccessList()),
		}), nil
	case dynamicFeeTxType:
		tipCap, feeCap, err := parseFeeCaps(core)
		if err != nil {
			return nil, err
		}
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:    big.NewInt(int64(core.GetChainID())),
			Nonce:      core.GetNonce(),
			GasTipCap:  tipCap,
			GasFeeCap:  feeCap,
			Gas:        core.GetGasLimit(),
			To:         to,
			Value:      value,
			Data:       payload,
			AccessList: accessListFromProto(core.GetAccessList()),
		}), nil
	case blobTxType:
		if to == nil {
			return nil, fmt.Errorf("blob tx requires non-empty recipient")
		}
		tipCap, feeCap, err := parseFeeCaps(core)
		if err != nil {
			return nil, err
		}
		blobFeeCap, blobHashes, sidecar, err := blobDataFromProto(core.GetBlobTxData())
		if err != nil {
			return nil, err
		}
		valU256, overflow := uint256.FromBig(value)
		if overflow {
			return nil, fmt.Errorf("blob tx value overflows uint256")
		}
		tipCapU, overflow := uint256.FromBig(tipCap)
		if overflow {
			return nil, fmt.Errorf("blob tx gas tip cap overflows uint256")
		}
		feeCapU, overflow := uint256.FromBig(feeCap)
		if overflow {
			return nil, fmt.Errorf("blob tx gas fee cap overflows uint256")
		}
		return types.NewTx(&types.BlobTx{
			ChainID:    uint256.NewInt(uint64(core.GetChainID())),
			Nonce:      core.GetNonce(),
			GasTipCap:  tipCapU,
			GasFeeCap:  feeCapU,
			Gas:        core.GetGasLimit(),
			To:         *to,
			Value:      valU256,
			Data:       payload,
			AccessList: accessListFromProto(core.GetAccessList()),
			BlobFeeCap: blobFeeCap,
			BlobHashes: blobHashes,
			Sidecar:    sidecar,
		}), nil
	case setCodeTxType:
		if to == nil {
			return nil, fmt.Errorf("setcode tx cannot create contract")
		}
		tipCap, feeCap, err := parseFeeCaps(core)
		if err != nil {
			return nil, err
		}
		authList, err := authListFromProto(core.GetSetCodeAuthList())
		if err != nil {
			return nil, err
		}
		valU256, overflow := uint256.FromBig(value)
		if overflow {
			return nil, fmt.Errorf("setcode tx value overflows uint256")
		}
		tipCapU, overflow := uint256.FromBig(tipCap)
		if overflow {
			return nil, fmt.Errorf("setcode tx gas tip cap overflows uint256")
		}
		feeCapU, overflow := uint256.FromBig(feeCap)
		if overflow {
			return nil, fmt.Errorf("setcode tx gas fee cap overflows uint256")
		}
		return types.NewTx(&types.SetCodeTx{
			ChainID:    uint256.NewInt(uint64(core.GetChainID())),
			Nonce:      core.GetNonce(),
			GasTipCap:  tipCapU,
			GasFeeCap:  feeCapU,
			Gas:        core.GetGasLimit(),
			To:         *to,
			Value:      valU256,
			Data:       payload,
			AccessList: accessListFromProto(core.GetAccessList()),
			AuthList:   authList,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported tx type %d", core.GetTxType())
	}
}

func extractBody(core *iotextypes.ActionCore) (*common.Address, *big.Int, []byte, error) {
	var (
		to      string
		value   = new(big.Int)
		payload []byte
	)
	switch {
	case core.GetTransfer() != nil:
		tsf := core.GetTransfer()
		to = tsf.Recipient
		value.SetString(tsf.Amount, 10)
		payload = tsf.Payload
	case core.GetExecution() != nil:
		exec := core.GetExecution()
		to = exec.Contract
		value.SetString(exec.Amount, 10)
		payload = exec.Data
	default:
		return nil, nil, nil, fmt.Errorf("invalid action type %T not supported", core)
	}

	if to == "" {
		return nil, value, payload, nil
	}
	addr, err := address.FromString(to)
	if err != nil {
		return nil, nil, nil, err
	}
	ethAddr := common.BytesToAddress(addr.Bytes())
	return &ethAddr, value, payload, nil
}

func parseFeeCaps(core *iotextypes.ActionCore) (*big.Int, *big.Int, error) {
	tipCap, ok := new(big.Int).SetString(core.GetGasTipCap(), 10)
	if !ok {
		return nil, nil, fmt.Errorf("invalid gasTipCap %q", core.GetGasTipCap())
	}
	feeCap, ok := new(big.Int).SetString(core.GetGasFeeCap(), 10)
	if !ok {
		return nil, nil, fmt.Errorf("invalid gasFeeCap %q", core.GetGasFeeCap())
	}
	return tipCap, feeCap, nil
}

// accessListFromProto mirrors iotex-core's encoding: AccessTuple.Address and
// StorageKeys are hex-encoded strings (no 0x prefix).
func accessListFromProto(list []*iotextypes.AccessTuple) types.AccessList {
	if len(list) == 0 {
		return nil
	}
	out := make(types.AccessList, len(list))
	for i, v := range list {
		out[i].Address = common.HexToAddress(v.Address)
		if n := len(v.StorageKeys); n > 0 {
			out[i].StorageKeys = make([]common.Hash, n)
			for j, k := range v.StorageKeys {
				out[i].StorageKeys[j] = common.HexToHash(k)
			}
		} else {
			out[i].StorageKeys = []common.Hash{}
		}
	}
	return out
}

func blobDataFromProto(pb *iotextypes.BlobTxData) (*uint256.Int, []common.Hash, *types.BlobTxSidecar, error) {
	if pb == nil {
		return nil, nil, nil, fmt.Errorf("blob tx requires blobTxData")
	}
	blobFeeCap := new(uint256.Int)
	if s := pb.GetBlobFeeCap(); len(s) > 0 {
		if err := blobFeeCap.SetFromDecimal(s); err != nil {
			return nil, nil, nil, fmt.Errorf("invalid blob fee cap: %w", err)
		}
	}
	hashes := make([]common.Hash, len(pb.GetBlobHashes()))
	for i, h := range pb.GetBlobHashes() {
		hashes[i] = common.BytesToHash(h)
	}
	sidecar, err := sidecarFromProto(pb.GetBlobTxSidecar())
	if err != nil {
		return nil, nil, nil, err
	}
	return blobFeeCap, hashes, sidecar, nil
}

func sidecarFromProto(pb *iotextypes.BlobTxSidecar) (*types.BlobTxSidecar, error) {
	if pb == nil || len(pb.GetBlobs()) == 0 {
		return nil, nil
	}
	sc := types.BlobTxSidecar{
		Blobs:       make([]kzg4844.Blob, len(pb.GetBlobs())),
		Commitments: make([]kzg4844.Commitment, len(pb.GetCommitments())),
		Proofs:      make([]kzg4844.Proof, len(pb.GetProofs())),
	}
	for i, b := range pb.GetBlobs() {
		if len(b) != len(kzg4844.Blob{}) {
			return nil, fmt.Errorf("invalid blob length at %d", i)
		}
		copy(sc.Blobs[i][:], b)
	}
	for i, c := range pb.GetCommitments() {
		if len(c) != len(kzg4844.Commitment{}) {
			return nil, fmt.Errorf("invalid commitment length at %d", i)
		}
		copy(sc.Commitments[i][:], c)
	}
	for i, p := range pb.GetProofs() {
		if len(p) != len(kzg4844.Proof{}) {
			return nil, fmt.Errorf("invalid proof length at %d", i)
		}
		copy(sc.Proofs[i][:], p)
	}
	return &sc, nil
}

func authListFromProto(list []*iotextypes.SetCodeAuthorization) ([]types.SetCodeAuthorization, error) {
	if len(list) == 0 {
		return nil, fmt.Errorf("setcode tx requires non-empty auth list")
	}
	out := make([]types.SetCodeAuthorization, len(list))
	for i, a := range list {
		r := new(uint256.Int).SetBytes(a.GetR())
		s := new(uint256.Int).SetBytes(a.GetS())
		out[i] = types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(uint64(a.GetChainID())),
			Address: common.BytesToAddress(a.GetAddress()),
			Nonce:   a.GetNonce(),
			V:       uint8(a.GetV()),
			R:       *r,
			S:       *s,
		}
	}
	return out, nil
}

func rlpSignedHash(tx *types.Transaction, chainID uint32, sig []byte) (hash.Hash256, error) {
	if len(sig) != 65 {
		return hash.ZeroHash256, fmt.Errorf("invalid signature length = %d, expecting 65", len(sig))
	}
	sc := make([]byte, 65)
	copy(sc, sig)
	if sc[64] >= 27 {
		sc[64] -= 27
	}

	signer := types.LatestSignerForChainID(big.NewInt(int64(chainID)))
	signedTx, err := tx.WithSignature(signer, sc)
	if err != nil {
		return hash.ZeroHash256, err
	}
	h := signedTx.Hash()
	return hash.BytesToHash256(h[:]), nil
}

// extractEthTxSig pulls the 65-byte [R||S||V] signature out of a signed eth
// tx in the exact form iotex-core's ExtractTypeSigPubkey produces, so the
// Action.Signature accompanying a TX_CONTAINER matches what the node derives
// from the raw tx. Typed txs carry recovery id 0/1 and are normalized to
// 27/28; legacy protected txs have the EIP-155 chain offset removed.
func extractEthTxSig(tx *types.Transaction) ([]byte, error) {
	V, R, S := tx.RawSignatureValues()
	switch tx.Type() {
	case types.LegacyTxType:
		if tx.Protected() {
			chainIDMul := tx.ChainId()
			V = new(big.Int).Sub(V, new(big.Int).Lsh(chainIDMul, 1))
			V = V.Sub(V, big.NewInt(8))
		}
	case types.AccessListTxType, types.DynamicFeeTxType, types.BlobTxType, types.SetCodeTxType:
		V = new(big.Int).Add(V, big.NewInt(27))
	default:
		return nil, fmt.Errorf("unsupported tx type %d", tx.Type())
	}
	if V.BitLen() > 8 {
		return nil, fmt.Errorf("invalid signature V value")
	}
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = byte(V.Uint64())
	return sig, nil
}
