package jwt

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/stretchr/testify/require"
)

func TestDecodeJWT(t *testing.T) {
	r := require.New(t)

	header := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9"
	payload := "eyJpYXQiOjE2MDU3NzQ5NjksImlzcyI6IjB4MDQ2MGVkNzAyYmY1YWNmZmE0NjM1ZWM1OTQ0OGZkNDFkOTQyMGNjMDc4NjZiMDc0M2VjMTdiNGJiYjI3YWZhZjA4NGNkOTc0M2Y1Y2Q1MjAwOWFmYzQxOTMyNDNiYWRkOGUyZGVmZGEyNGYxM2MzMjY5YzQ4OTkwM2Q1OWRkZWJlMCIsInN1YiI6Imh0dHA6Ly9leGFtcGxlLmNvbWUvMTIzNCJ9"
	signature := "ftpWYERxNYYkDrFVDVYao-kofdMVu8_J3GcKxF1JQOhwtscW9d3BAsrxFnKQ5p2o6XYRKiDHxOX_gCXGMymM3A"

	// JWT token = header.payload.signature
	jwtStr := header + "." + payload + "." + signature
	tok, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		keyHex := token.Claims.(jwt.MapClaims)["iss"]
		key, err := crypto.HexStringToPublicKey(keyHex.(string)[2:])
		if err != nil {
			return nil, err
		}
		return key.EcdsaPublicKey(), nil
	})
	r.NoError(err)

	// header =
	// {
	//   "alg": "ES256",
	//   "typ": "JWT"
	// }
	r.Equal("ES256", tok.Header["alg"].(string))
	r.Equal("JWT", tok.Header["typ"].(string))

	// payload =
	// {
	//   "exp": 1605772249,
	//   "iss": "0x0460ed702bf5acffa4635ec59448fd41d9420cc07866b0743ec17b4bbb27afaf084cd9743f5cd52009afc4193243badd8e2defda24f13c3269c489903d59ddebe0",
	//   "sub": "http://example.come/1234"
	// }
	iss := "0x0460ed702bf5acffa4635ec59448fd41d9420cc07866b0743ec17b4bbb27afaf084cd9743f5cd52009afc4193243badd8e2defda24f13c3269c489903d59ddebe0"
	claim := tok.Claims.(jwt.MapClaims)
	r.Equal(3, len(claim))
	r.EqualValues(1605774969, claim["iat"].(float64))
	r.Equal(iss, claim["iss"].(string))
	r.Equal("http://example.come/1234", claim["sub"].(string))

	// signature
	r.Equal(signature, tok.Signature)

	// verify JWT
	token, err := VerifyJWT(jwtStr)
	r.NoError(err)
	r.EqualValues(1605774969, token.IssuedAt)
	r.EqualValues(0, token.ExpiresAt)
	r.Equal(iss, token.Issuer)
	r.Equal("http://example.come/1234", token.Subject)
	r.Equal("ES256", token.SignMethod)

	// signature = ecdsaSign(hash(header.payload))
	hasher := jwt.SigningMethodES256.Hash.New()
	hasher.Write([]byte(header + "." + payload))
	sigBytes, err := hex.DecodeString(token.SigHex)
	r.NoError(err)
	// recover public key from signature and verify
	pk, err := crypto.RecoverPubkey(hasher.Sum(nil), append(sigBytes, 1))
	r.NoError(err)
	r.Equal(iss[2:], pk.HexString())
}

func TestSignVerifyJWT(t *testing.T) {
	r := require.New(t)

	a, err := crypto.GenerateKey()
	r.NoError(err)

	now := time.Now().Unix()
	jwtTests := []struct {
		iss, exp   int64
		url, scope string
		errStr     string
	}{
		{now, 0, "http://example.come/1234", CREATE, ""},
		{now, now + 1, "http://example.come/1234", READ, ""},
		{now, now + 2, "http://example.come/4321", DELETE, ""},
		{now, now - 1, "http://example.come/1234", "", "token is expired by"},
		{now + 1, now, "http://example.come/1234", "", "Token used before issued"},
	}

	issuer := "0x" + a.PublicKey().HexString()
	for _, v := range jwtTests {
		jwtStr, err := SignJWT(v.iss, v.exp, v.url, v.scope, a)
		r.NoError(err)
		token, err := VerifyJWT(jwtStr)
		if v.errStr != "" {
			r.True(strings.HasPrefix(err.Error(), v.errStr))
			continue
		}
		r.NoError(err)
		r.Equal(v.iss, token.IssuedAt)
		r.Equal(v.exp, token.ExpiresAt)
		r.Equal(issuer, token.Issuer)
		r.Equal(v.url, token.Subject)
		r.Equal(v.scope, token.Scope)
		r.Equal("ES256", token.SignMethod)

		// signing the token with a diff key fails the verification
		claim := &claimWithScope{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: v.exp,
				IssuedAt:  v.iss,
				Issuer:    issuer,
				Subject:   v.url,
			},
			Scope: v.scope,
		}
		tok := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
		str, err := tok.SigningString()
		r.NoError(err)
		b, err := crypto.GenerateKey()
		r.NoError(err)
		sig, err := tok.Method.Sign(str, b.EcdsaPrivateKey())
		r.NoError(err)
		str = strings.Join([]string{str, sig}, ".")
		r.NotEqual(jwtStr, str)
		_, err = VerifyJWT(str)
		r.Equal(jwt.ErrECDSAVerification.Error(), err.Error())
	}
}
