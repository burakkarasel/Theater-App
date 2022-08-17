package db

import (
	"context"
	"testing"
	"time"

	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/stretchr/testify/require"
)

// createRandomDirector takes testing.T as arg and returns a random director
func createRandomDirector(t *testing.T) Director {
	arg := CreateDirectorParams{
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Oscars:    util.RandomInt(0, 5),
	}

	director, err := testQueries.CreateDirector(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, director)

	require.Equal(t, arg.FirstName, director.FirstName)
	require.Equal(t, arg.LastName, director.LastName)
	require.Equal(t, arg.Oscars, director.Oscars)
	require.NotZero(t, director.ID)
	require.NotZero(t, director.CreatedAt)

	return director
}

// TestCreateDirector tests CreateDirector DB Operation
func TestCreateDirector(t *testing.T) {
	createRandomDirector(t)
}

// TestGetDirector tests GetDirector DB operation
func TestGetDirector(t *testing.T) {
	d1 := createRandomDirector(t)
	d2, err := testQueries.GetDirector(context.Background(), d1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, d2)

	require.Equal(t, d1.ID, d2.ID)
	require.Equal(t, d1.FirstName, d2.FirstName)
	require.Equal(t, d1.LastName, d2.LastName)
	require.Equal(t, d1.Oscars, d2.Oscars)
	require.WithinDuration(t, d1.CreatedAt, d2.CreatedAt, time.Second)
}

// TestListDirectors creates 10 accounts and gets 5 of those from DB
func TestListDirectors(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomDirector(t)
	}

	arg := ListDirectorsParams{
		Limit:  5,
		Offset: 5,
	}

	directors, err := testQueries.ListDirectors(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, directors, 5)

	for _, v := range directors {
		require.NotEmpty(t, v)
	}
}
