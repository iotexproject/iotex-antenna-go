// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package jwt

import (
	"encoding/hex"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/go-pkgs/crypto"
)

// const
const (
	CREATE = "Create"
	READ   = "Read"
	UPDATE = "Update"
	DELETE = "Delete"
)

type (
	// JWT is a JWT object
	JWT struct {
		IssuedAt   int64
		ExpiresAt  int64
		Issuer     string
		Subject    string
		Scope      string
		SignMethod string
		SigHex     string
	}

	claimWithScope struct {
		jwt.StandardClaims
		Scope string
	}
)

// SignJWT creates a JWT
func SignJWT(issue, expire int64, subject, scope string, key crypto.PrivateKey) (string, error) {
	c := &claimWithScope{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			IssuedAt:  issue,
			Issuer:    "0x" + key.PublicKey().HexString(),
			Subject:   subject,
		},
		Scope: scope,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, c)
	return token.SignedString(key.EcdsaPrivateKey())
}

// VerifyJWT verifies the JWT
func VerifyJWT(jwtString string) (*JWT, error) {
	claim := &claimWithScope{}
	token, err := jwt.ParseWithClaims(jwtString, claim, func(token *jwt.Token) (interface{}, error) {
		keyHex := claim.Issuer
		if keyHex[:2] == "0x" || keyHex[:2] == "0X" {
			keyHex = keyHex[2:]
		}
		key, err := crypto.HexStringToPublicKey(keyHex)
		if err != nil {
			return nil, err
		}
		return key.EcdsaPublicKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid || len(token.Header) < 2 {
		// should not happen with a success parsing, check anyway
		return nil, errors.New("invalid token")
	}

	// decode signature
	sig, err := jwt.DecodeSegment(token.Signature)
	if err != nil {
		return nil, err
	}

	return &JWT{
		IssuedAt:   claim.IssuedAt,
		ExpiresAt:  claim.ExpiresAt,
		Issuer:     claim.Issuer,
		Subject:    claim.Subject,
		Scope:      claim.Scope,
		SignMethod: token.Header["alg"].(string),
		SigHex:     hex.EncodeToString(sig),
	}, nil
}
