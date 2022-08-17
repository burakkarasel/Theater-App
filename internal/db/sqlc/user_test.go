package db

import (
	"context"
	"testing"
	"time"

	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/stretchr/testify/require"
)

// createRandomUser creates a random user in DB
func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: util.RandomString(8),
		Email:          util.RandomString(6),
		AccessLevel:    1,
	}

	u, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	require.Equal(t, arg.Username, u.Username)
	require.Equal(t, arg.HashedPassword, u.HashedPassword)
	require.Equal(t, arg.Email, u.Email)
	require.Equal(t, arg.AccessLevel, u.AccessLevel)
	require.NotZero(t, u.CreatedAt)

	return u
}

// TestCreateUser tests CreateUser DB operation
func TestCreateUser(t *testing.T) {
	u := createRandomUser(t)
	require.NotEmpty(t, u)
}

// TestGetUser tests GetUser DB operation
func TestGetUser(t *testing.T) {
	u1 := createRandomUser(t)

	u2, err := testQueries.GetUser(context.Background(), u1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, u2)

	require.Equal(t, u1.AccessLevel, u2.AccessLevel)
	require.Equal(t, u1.Username, u2.Username)
	require.Equal(t, u1.Email, u2.Email)
	require.Equal(t, u1.HashedPassword, u2.HashedPassword)
	require.WithinDuration(t, u1.CreatedAt, u2.CreatedAt, time.Second)
}
