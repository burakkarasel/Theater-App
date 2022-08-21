package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/burakkarasel/Theatre-API/internal/db/mock"
	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TestCreateMovieAPI tests createMovie handler
func TestCreateMovieAPI(t *testing.T) {
	movie := randomMovie()
	testCases := []struct {
		name           string
		body           gin.H
		buildStubs     func(store *mockdb.MockStore)
		checkResponses func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     movie.Movie.Summary,
				"poster":      movie.Movie.Poster,
				"director_id": movie.Movie.DirectorID,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMovieParams{
					Title:      movie.Movie.Title,
					Summary:    movie.Movie.Summary,
					Poster:     movie.Movie.Poster,
					Rating:     movie.Movie.Rating,
					DirectorID: movie.Movie.DirectorID,
				}

				store.EXPECT().CreateMovie(gomock.Any(), gomock.Eq(arg)).Times(1).Return(movie.Movie, nil)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchCreateMovie(t, w.Body, movie.Movie)
			},
		},
		{
			name: "Invalid Title",
			body: gin.H{
				"title":       "a",
				"summary":     movie.Movie.Summary,
				"poster":      movie.Movie.Poster,
				"director_id": movie.Movie.DirectorID,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Summary",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     "asd",
				"poster":      movie.Movie.Poster,
				"director_id": movie.Movie.DirectorID,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Poster",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     movie.Movie.Summary,
				"poster":      "asd",
				"director_id": movie.Movie.DirectorID,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Director ID",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     movie.Movie.Summary,
				"poster":      movie.Movie.Poster,
				"director_id": -3,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Rating",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     movie.Movie.Summary,
				"poster":      movie.Movie.Poster,
				"director_id": movie.Movie.DirectorID,
				"rating":      -3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"title":       movie.Movie.Title,
				"summary":     movie.Movie.Summary,
				"poster":      movie.Movie.Poster,
				"director_id": movie.Movie.DirectorID,
				"rating":      movie.Movie.Rating,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateMovieParams{
					Title:      movie.Movie.Title,
					Summary:    movie.Movie.Summary,
					Poster:     movie.Movie.Poster,
					Rating:     movie.Movie.Rating,
					DirectorID: movie.Movie.DirectorID,
				}

				store.EXPECT().CreateMovie(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Movie{}, sql.ErrConnDone)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			w := httptest.NewRecorder()

			url := "/movies"
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)

			tt.checkResponses(t, w)
		})
	}
}

// TestGetMovieAPI tests getMovie handler
func TestGetMovieAPI(t *testing.T) {
	movie := randomMovie()
	testCases := []struct {
		name           string
		movieID        int64
		buildStubs     func(store *mockdb.MockStore)
		checkResponses func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			movieID: movie.Movie.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.Movie.ID)).Times(1).Return(movie.Movie, nil)
				store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(movie.Movie.DirectorID)).Times(1).Return(movie.Director, nil)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatchMovie(t, w.Body, movie)
			},
		},
		{
			name:    "Invalid ID",
			movieID: -3,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:    "Movie Not Found",
			movieID: movie.Movie.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.Movie.ID)).Times(1).Return(db.Movie{}, sql.ErrNoRows)
				store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:    "Movie Internal Error",
			movieID: movie.Movie.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.Movie.ID)).Times(1).Return(db.Movie{}, sql.ErrConnDone)
				store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:    "Director Internal Error",
			movieID: movie.Movie.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.Movie.ID)).Times(1).Return(movie.Movie, nil)
				store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(movie.Movie.DirectorID)).Times(1).Return(db.Director{}, sql.ErrConnDone)
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			w := httptest.NewRecorder()

			url := fmt.Sprintf("/movies/%d", tt.movieID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)

			tt.checkResponses(t, w)
		})
	}
}

// TestListMoviesAPI tests listMovies handler
func TestListMoviesAPI(t *testing.T) {
	var movies []db.Movie
	var directors []db.Director
	n := 5
	for i := 0; i < n; i++ {
		m := randomMovie()
		movies = append(movies, m.Movie)
		directors = append(directors, m.Director)
	}

	testCases := []struct {
		name           string
		query          string
		buildStubs     func(store *mockdb.MockStore)
		checkResponses func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "?count=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListMovies(gomock.Any(), gomock.Eq(int32(n))).Times(1).Return(movies, nil)
				for i := 0; i < n; i++ {
					store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(movies[i].DirectorID)).Times(1).Return(directors[i], nil)
				}
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:  "Invalid Count",
			query: "?count=-3",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListMovies(gomock.Any(), gomock.Any()).Times(0)
				for i := 0; i < n; i++ {
					store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
				}
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Movie Internal Error",
			query: "?count=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListMovies(gomock.Any(), gomock.Eq(int32(n))).Times(1).Return([]db.Movie{}, sql.ErrConnDone)
				for i := 0; i < n; i++ {
					store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
				}
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:  "Director Internal Error",
			query: "?count=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListMovies(gomock.Any(), gomock.Eq(int32(n))).Times(1).Return(movies, nil)
				for i := 0; i < 1; i++ {
					store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(movies[i].DirectorID)).Times(1).Return(db.Director{}, sql.ErrConnDone)
				}
			},
			checkResponses: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := newTestServer(t, store)
			w := httptest.NewRecorder()

			url := fmt.Sprintf("/movies%s", tt.query)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)

			tt.checkResponses(t, w)
		})
	}
}

// randomMovie creates a random movie
func randomMovie() GetMovieResponse {
	d := randomDirector()
	return GetMovieResponse{
		Movie: db.Movie{
			Title:      util.RandomName(),
			Summary:    util.RandomString(10),
			Poster:     util.RandomString(10),
			Rating:     int16(util.RandomInt(5, 10)),
			ID:         util.RandomInt(1, 1000),
			DirectorID: d.ID,
		},
		Director: db.Director{
			FirstName: d.FirstName,
			LastName:  d.LastName,
			ID:        d.ID,
			Oscars:    d.Oscars,
		},
	}
}

// requireBodyMatch checks for a given body and response's body
func requireBodyMatchMovie(t *testing.T, body *bytes.Buffer, resp GetMovieResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMovieResponse GetMovieResponse
	err = json.Unmarshal(data, &gotMovieResponse)
	require.NoError(t, err)
	require.Equal(t, resp, gotMovieResponse)
}

// requireBodyMatch checks for a given body and response's body
func requireBodyMatchCreateMovie(t *testing.T, body *bytes.Buffer, resp db.Movie) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotMovieResponse db.Movie
	err = json.Unmarshal(data, &gotMovieResponse)
	require.NoError(t, err)
	require.Equal(t, resp, gotMovieResponse)
}
