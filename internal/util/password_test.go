package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword testsHashPassword
func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword(RandomString(8))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

// TestCompareHashedPassword tests CompareHashedPassword
func TestCompareHashedPassword(t *testing.T) {
	pw := RandomString(8)
	hashedPassword1, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CompareHashedPassword(pw, hashedPassword1)
	require.NoError(t, err)

	wrongPw := RandomString(6)
	err = CompareHashedPassword(wrongPw, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
