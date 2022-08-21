package token

import (
	"testing"
	"time"

	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

// TestJWTMaker tests JWTMaker function
func TestJWTMaker(t *testing.T) {
	// happy case
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	// invalid case
	invalidMaker, err := NewJWTMaker(util.RandomString(30))
	require.EqualError(t, err, ErrInvalidSecretKeySize.Error())
	require.Empty(t, invalidMaker)

	// then happy case for valid token
	username := util.RandomName()
	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiresAt, time.Second)
}

// TestExpiredJWTToken tests for an expired token
func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, err := maker.CreateToken(util.RandomName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	jwtToken, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, jwtToken)
}

// TestInvalidJWTToken tests for invalid JWT token
func TestInvalidJWTToken(t *testing.T) {
	// first we create a new payload and sign it with none method for only tests
	payload, err := NewPayload(util.RandomName(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// then we create the maker and try to verify the invalid token
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
