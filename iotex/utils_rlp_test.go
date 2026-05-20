package iotex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
)

func accountFromHex(key string) (account.Account, error) {
	return account.HexStringToAccount(key)
}

const (
	_rlpTestKey     = "73c7b4a62bf165dccf8ebdea8278db811efd5b8638e2ed9683d2d94889450426"
	_rlpTestTo      = "io1emxf8zzqckhgjde6dqd97ts0y3q496gm3fdrl6"
	_rlpTestChainID = uint32(4690)
)

// fixedSig is a syntactically valid 65-byte secp256k1 signature; the test
// pipes the same bytes through both paths so the V/R/S values don't need to
// actually verify against a key — we only care that actionToRLP and a
// hand-built types.Transaction produce identical hashes.
var fixedSig = func() []byte {
	out := make([]byte, 65)
	for i := range out[:64] {
		out[i] = byte(i + 1)
	}
	out[64] = 27
	return out
}()

func toEthAddr(t *testing.T, ioAddr string) common.Address {
	t.Helper()
	addr, err := address.FromString(ioAddr)
	require.NoError(t, err)
	return common.BytesToAddress(addr.Bytes())
}

func signAndHash(t *testing.T, tx *types.Transaction) []byte {
	t.Helper()
	h, err := rlpSignedHash(tx, _rlpTestChainID, fixedSig)
	require.NoError(t, err)
	return h[:]
}

func TestActionToRLP_Legacy(t *testing.T) {
	to := toEthAddr(t, _rlpTestTo)
	core := &iotextypes.ActionCore{
		ChainID:  _rlpTestChainID,
		Nonce:    7,
		GasLimit: 21000,
		GasPrice: "1000000000000",
		Action: &iotextypes.ActionCore_Transfer{
			Transfer: &iotextypes.Transfer{
				Amount:    "160000000000000000",
				Recipient: _rlpTestTo,
			},
		},
	}
	gotTx, err := actionToRLP(core)
	require.NoError(t, err)
	require.Equal(t, uint8(types.LegacyTxType), gotTx.Type())

	wantTx := types.NewTransaction(7, to, big.NewInt(160000000000000000), 21000, big.NewInt(1000000000000), nil)
	require.Equal(t, signAndHash(t, wantTx), signAndHash(t, gotTx))
}

func TestActionToRLP_AccessList(t *testing.T) {
	to := toEthAddr(t, _rlpTestTo)
	addrHex := "1234567890123456789012345678901234567890"
	keyHex := "00000000000000000000000000000000000000000000000000000000000000ff"

	core := &iotextypes.ActionCore{
		ChainID:  _rlpTestChainID,
		Nonce:    9,
		GasLimit: 50000,
		GasPrice: "2000000000000",
		TxType:   accessListTxType,
		AccessList: []*iotextypes.AccessTuple{
			{Address: addrHex, StorageKeys: []string{keyHex}},
		},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{
				Amount:   "0",
				Contract: _rlpTestTo,
				Data:     []byte{0xde, 0xad},
			},
		},
	}
	gotTx, err := actionToRLP(core)
	require.NoError(t, err)
	require.Equal(t, uint8(types.AccessListTxType), gotTx.Type())

	wantTx := types.NewTx(&types.AccessListTx{
		ChainID:  big.NewInt(int64(_rlpTestChainID)),
		Nonce:    9,
		GasPrice: big.NewInt(2000000000000),
		Gas:      50000,
		To:       &to,
		Value:    big.NewInt(0),
		Data:     []byte{0xde, 0xad},
		AccessList: types.AccessList{
			{
				Address:     common.HexToAddress(addrHex),
				StorageKeys: []common.Hash{common.HexToHash(keyHex)},
			},
		},
	})
	require.Equal(t, signAndHash(t, wantTx), signAndHash(t, gotTx))
}

func TestActionToRLP_DynamicFee(t *testing.T) {
	to := toEthAddr(t, _rlpTestTo)
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     11,
		GasLimit:  60000,
		GasTipCap: "1500000000",
		GasFeeCap: "3000000000",
		TxType:    dynamicFeeTxType,
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{
				Amount:   "100",
				Contract: _rlpTestTo,
				Data:     []byte{0xbe, 0xef},
			},
		},
	}
	gotTx, err := actionToRLP(core)
	require.NoError(t, err)
	require.Equal(t, uint8(types.DynamicFeeTxType), gotTx.Type())

	wantTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(int64(_rlpTestChainID)),
		Nonce:     11,
		GasTipCap: big.NewInt(1500000000),
		GasFeeCap: big.NewInt(3000000000),
		Gas:       60000,
		To:        &to,
		Value:     big.NewInt(100),
		Data:      []byte{0xbe, 0xef},
	})
	require.Equal(t, signAndHash(t, wantTx), signAndHash(t, gotTx))
}

func TestActionToRLP_Blob(t *testing.T) {
	to := toEthAddr(t, _rlpTestTo)
	blobHash := common.HexToHash("0x01abcdef00000000000000000000000000000000000000000000000000000000")
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     13,
		GasLimit:  70000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    blobTxType,
		BlobTxData: &iotextypes.BlobTxData{
			BlobFeeCap: "500000000",
			BlobHashes: [][]byte{blobHash.Bytes()},
		},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{
				Amount:   "0",
				Contract: _rlpTestTo,
			},
		},
	}
	gotTx, err := actionToRLP(core)
	require.NoError(t, err)
	require.Equal(t, uint8(types.BlobTxType), gotTx.Type())

	wantTx := types.NewTx(&types.BlobTx{
		ChainID:    uint256.NewInt(uint64(_rlpTestChainID)),
		Nonce:      13,
		GasTipCap:  uint256.NewInt(1000000000),
		GasFeeCap:  uint256.NewInt(2000000000),
		Gas:        70000,
		To:         to,
		Value:      uint256.NewInt(0),
		BlobFeeCap: uint256.NewInt(500000000),
		BlobHashes: []common.Hash{blobHash},
	})
	require.Equal(t, signAndHash(t, wantTx), signAndHash(t, gotTx))
}

func TestActionToRLP_SetCode(t *testing.T) {
	to := toEthAddr(t, _rlpTestTo)
	authAddr := common.HexToAddress("0xabababababababababababababababababababab")
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     17,
		GasLimit:  80000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    setCodeTxType,
		SetCodeAuthList: []*iotextypes.SetCodeAuthorization{
			{
				ChainID: _rlpTestChainID,
				Address: authAddr.Bytes(),
				Nonce:   3,
				V:       1,
				R:       []byte{0xaa, 0xbb},
				S:       []byte{0xcc, 0xdd},
			},
		},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{
				Amount:   "0",
				Contract: _rlpTestTo,
			},
		},
	}
	gotTx, err := actionToRLP(core)
	require.NoError(t, err)
	require.Equal(t, uint8(types.SetCodeTxType), gotTx.Type())

	wantTx := types.NewTx(&types.SetCodeTx{
		ChainID:   uint256.NewInt(uint64(_rlpTestChainID)),
		Nonce:     17,
		GasTipCap: uint256.NewInt(1000000000),
		GasFeeCap: uint256.NewInt(2000000000),
		Gas:       80000,
		To:        to,
		Value:     uint256.NewInt(0),
		AuthList: []types.SetCodeAuthorization{{
			ChainID: *uint256.NewInt(uint64(_rlpTestChainID)),
			Address: authAddr,
			Nonce:   3,
			V:       1,
			R:       *new(uint256.Int).SetBytes([]byte{0xaa, 0xbb}),
			S:       *new(uint256.Int).SetBytes([]byte{0xcc, 0xdd}),
		}},
	})
	require.Equal(t, signAndHash(t, wantTx), signAndHash(t, gotTx))
}

func TestActionToRLP_UnsupportedType(t *testing.T) {
	core := &iotextypes.ActionCore{
		ChainID:  _rlpTestChainID,
		Nonce:    1,
		GasLimit: 21000,
		GasPrice: "1",
		TxType:   99,
		Action: &iotextypes.ActionCore_Transfer{
			Transfer: &iotextypes.Transfer{Amount: "0", Recipient: _rlpTestTo},
		},
	}
	_, err := actionToRLP(core)
	require.Error(t, err)
}

func TestRlpSignedHash_BadSig(t *testing.T) {
	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)
	_, err := rlpSignedHash(tx, _rlpTestChainID, []byte{1, 2, 3})
	require.Error(t, err)
}

func TestSignContainer_DynamicFee(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     5,
		GasLimit:  21000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    dynamicFeeTxType,
		Action: &iotextypes.ActionCore_Transfer{
			Transfer: &iotextypes.Transfer{
				Amount:    "100",
				Recipient: _rlpTestTo,
			},
		},
	}
	sealed, err := signContainer(acc, core)
	require.NoError(t, err)
	require.Equal(t, iotextypes.Encoding_TX_CONTAINER, sealed.Encoding)
	require.Len(t, sealed.Signature, 65)
	require.Equal(t, acc.PublicKey().Bytes(), sealed.SenderPubKey)

	// the container carries a raw eth tx that round-trips to a DynamicFee tx
	raw := sealed.GetCore().GetTxContainer().GetRaw()
	require.NotEmpty(t, raw)
	decoded := new(types.Transaction)
	require.NoError(t, decoded.UnmarshalBinary(raw))
	require.Equal(t, uint8(types.DynamicFeeTxType), decoded.Type())
	require.Equal(t, uint64(5), decoded.Nonce())

	// ActionHash on the sealed action equals the eth tx hash
	h, err := ActionHash(sealed, _rlpTestChainID)
	require.NoError(t, err)
	require.Equal(t, decoded.Hash().Bytes(), h[:])
}

func TestSignContainer_RejectsUnsupportedBody(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	core := &iotextypes.ActionCore{
		ChainID:  _rlpTestChainID,
		Nonce:    1,
		GasLimit: 21000,
		GasPrice: "1",
		// no transfer/execution body — actionToRLP must error
	}
	_, err = signContainer(acc, core)
	require.Error(t, err)
}

// TestSignContainer_SignatureRecoversSigner proves the Action.Signature that
// accompanies a TX_CONTAINER lets the node recover the actual signer — i.e.
// extractEthTxSig produces the V/R/S form the node's recovery expects.
func TestSignContainer_SignatureRecoversSigner(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     2,
		GasLimit:  21000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    dynamicFeeTxType,
		Action: &iotextypes.ActionCore_Transfer{
			Transfer: &iotextypes.Transfer{Amount: "1", Recipient: _rlpTestTo},
		},
	}
	sealed, err := signContainer(acc, core)
	require.NoError(t, err)

	decoded := new(types.Transaction)
	require.NoError(t, decoded.UnmarshalBinary(sealed.GetCore().GetTxContainer().GetRaw()))
	signer := types.LatestSignerForChainID(big.NewInt(int64(_rlpTestChainID)))
	pub, err := crypto.RecoverPubkey(signer.Hash(decoded).Bytes(), sealed.Signature)
	require.NoError(t, err)
	require.Equal(t, acc.PublicKey().Bytes(), pub.Bytes())
}

// signContainerRoundTrip runs signContainer, then asserts the sealed action
// is a TX_CONTAINER whose raw bytes UnmarshalBinary back to the expected tx
// type and whose Action signature recovers the signer. It returns the decoded
// tx so callers can assert type-specific fields.
func signContainerRoundTrip(t *testing.T, acc account.Account, core *iotextypes.ActionCore, wantType uint8) *types.Transaction {
	t.Helper()
	sealed, err := signContainer(acc, core)
	require.NoError(t, err)
	require.Equal(t, iotextypes.Encoding_TX_CONTAINER, sealed.Encoding)
	require.Len(t, sealed.Signature, 65)
	require.Equal(t, acc.PublicKey().Bytes(), sealed.SenderPubKey)

	raw := sealed.GetCore().GetTxContainer().GetRaw()
	require.NotEmpty(t, raw)
	decoded := new(types.Transaction)
	require.NoError(t, decoded.UnmarshalBinary(raw))
	require.Equal(t, wantType, decoded.Type())

	signer := types.LatestSignerForChainID(big.NewInt(int64(core.GetChainID())))
	pub, err := crypto.RecoverPubkey(signer.Hash(decoded).Bytes(), sealed.Signature)
	require.NoError(t, err)
	require.Equal(t, acc.PublicKey().Bytes(), pub.Bytes())
	return decoded
}

func TestSignContainer_AccessList(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	core := &iotextypes.ActionCore{
		ChainID:  _rlpTestChainID,
		Nonce:    3,
		GasLimit: 50000,
		GasPrice: "1000000000000",
		TxType:   accessListTxType,
		AccessList: []*iotextypes.AccessTuple{{
			Address:     "1234567890123456789012345678901234567890",
			StorageKeys: []string{"00000000000000000000000000000000000000000000000000000000000000ff"},
		}},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{Amount: "0", Contract: _rlpTestTo, Data: []byte{0xde, 0xad}},
		},
	}
	decoded := signContainerRoundTrip(t, acc, core, types.AccessListTxType)
	require.Equal(t, uint64(3), decoded.Nonce())
	require.Len(t, decoded.AccessList(), 1)
	require.Len(t, decoded.AccessList()[0].StorageKeys, 1)
}

func TestSignContainer_Blob(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	blobHash := common.HexToHash("0x01abcdef00000000000000000000000000000000000000000000000000000000")
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     4,
		GasLimit:  70000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    blobTxType,
		BlobTxData: &iotextypes.BlobTxData{
			BlobFeeCap: "500000000",
			BlobHashes: [][]byte{blobHash.Bytes()},
			// structurally-valid dummy sidecar: marshaling does not verify KZG,
			// so correctly-sized zero blobs/commitments/proofs round-trip.
			BlobTxSidecar: &iotextypes.BlobTxSidecar{
				Blobs:       [][]byte{make([]byte, 131072)},
				Commitments: [][]byte{make([]byte, 48)},
				Proofs:      [][]byte{make([]byte, 48)},
			},
		},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{Amount: "0", Contract: _rlpTestTo},
		},
	}
	decoded := signContainerRoundTrip(t, acc, core, types.BlobTxType)
	require.Equal(t, uint64(4), decoded.Nonce())
	require.Len(t, decoded.BlobHashes(), 1)
	require.Equal(t, blobHash, decoded.BlobHashes()[0])
	require.NotNil(t, decoded.BlobTxSidecar())
	require.Len(t, decoded.BlobTxSidecar().Blobs, 1)
}

func TestSignContainer_SetCode(t *testing.T) {
	acc, err := accountFromHex(_rlpTestKey)
	require.NoError(t, err)
	authAddr := common.HexToAddress("0xabababababababababababababababababababab")
	core := &iotextypes.ActionCore{
		ChainID:   _rlpTestChainID,
		Nonce:     6,
		GasLimit:  80000,
		GasTipCap: "1000000000",
		GasFeeCap: "2000000000",
		TxType:    setCodeTxType,
		SetCodeAuthList: []*iotextypes.SetCodeAuthorization{{
			ChainID: _rlpTestChainID,
			Address: authAddr.Bytes(),
			Nonce:   1,
			V:       1,
			R:       []byte{0xaa, 0xbb},
			S:       []byte{0xcc, 0xdd},
		}},
		Action: &iotextypes.ActionCore_Execution{
			Execution: &iotextypes.Execution{Amount: "0", Contract: _rlpTestTo},
		},
	}
	decoded := signContainerRoundTrip(t, acc, core, types.SetCodeTxType)
	require.Equal(t, uint64(6), decoded.Nonce())
	require.Len(t, decoded.SetCodeAuthorizations(), 1)
	require.Equal(t, authAddr, decoded.SetCodeAuthorizations()[0].Address)
}
