package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/stretchr/testify/require"
)

// createRandomMovie Creates a random movie
func createRandomMovie(t *testing.T) Movie {
	director := createRandomDirector(t)
	arg := CreateMovieParams{
		Title:      util.RandomName(),
		DirectorID: director.ID,
		Rating:     fmt.Sprintf("%d", util.RandomInt(6, 10)),
		Poster:     util.RandomString(10),
	}

	m, err := testQueries.CreateMovie(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, m)

	require.NotZero(t, m.ID)
	require.NotZero(t, m.CreatedAt)
	require.Equal(t, m.Title, arg.Title)
	require.Equal(t, m.DirectorID, arg.DirectorID)
	require.Equal(t, m.Rating, arg.Rating)
	require.Equal(t, m.Poster, arg.Poster)

	return m
}

// TestCreateMovie tests CreateMovie DB operation
func TestCreateMovie(t *testing.T) {
	m := createRandomMovie(t)
	require.NotEmpty(t, m)
}

// TestGetMovie tests GetMovie DB operation
func TestGetMovie(t *testing.T) {
	m1 := createRandomMovie(t)
	m2, err := testQueries.GetMovie(context.Background(), m1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, m2)

	require.Equal(t, m1.DirectorID, m2.DirectorID)
	require.Equal(t, m1.ID, m2.ID)
	require.Equal(t, m1.Poster, m2.Poster)
	require.Equal(t, m1.Title, m2.Title)
	require.Equal(t, m1.Rating, m2.Rating)
	require.WithinDuration(t, m1.CreatedAt, m2.CreatedAt, time.Second)
}

// TestListMovies tests ListMovies DB operation
func TestListMovies(t *testing.T) {
	for i := 0; i < 8; i++ {
		createRandomMovie(t)
	}

	movies, err := testQueries.ListMovies(context.Background())
	require.NoError(t, err)
	require.Len(t, movies, 8)

	for _, v := range movies {
		require.NotEmpty(t, v)
	}
}

// TestDeleteMovie tests DeleteMovie DB operation
func TestDeleteMovie(t *testing.T) {
	m1 := createRandomMovie(t)

	err := testQueries.DeleteMovie(context.Background(), m1.ID)
	require.NoError(t, err)

	m2, err := testQueries.GetMovie(context.Background(), m1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, m2)
}
