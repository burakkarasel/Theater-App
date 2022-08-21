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

// TestGetDirectorAPI tests getDirector handler
func TestGetDirectorAPI(t *testing.T) {
	director := randomDirector()
	testCases := []struct {
		name          string
		directorID    int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:       "OK",
			directorID: director.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(director.ID)).Times(1).Return(director, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatch(t, w.Body, director)
			},
		},
		{
			name:       "Bindings Error",
			directorID: -3,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:       "Director Not Found",
			directorID: director.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(director.ID)).Times(1).Return(db.Director{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:       "Internal Server Error",
			directorID: director.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetDirector(gomock.Any(), gomock.Eq(director.ID)).Times(1).Return(db.Director{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
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

			// start test server
			server := newTestServer(t, store)
			w := httptest.NewRecorder()

			url := fmt.Sprintf("/directors/%d", tt.directorID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)
			tt.checkResponse(t, w)
		})
	}
}

// TestCreateDirectorAPI tests createDirector handler
func TestCreateDirectorAPI(t *testing.T) {
	director := randomDirector()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"first_name": director.FirstName,
				"last_name":  director.LastName,
				"oscars":     director.Oscars,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDirectorParams{
					FirstName: director.FirstName,
					LastName:  director.LastName,
					Oscars:    director.Oscars,
				}
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Eq(arg)).Times(1).Return(director, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireBodyMatch(t, w.Body, director)
			},
		},
		{
			name: "Invalid First Name",
			body: gin.H{
				"first_name": "a",
				"last_name":  director.LastName,
				"oscars":     director.Oscars,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Last Name",
			body: gin.H{
				"first_name": director.FirstName,
				"last_name":  "a",
				"oscars":     director.Oscars,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Oscars",
			body: gin.H{
				"first_name": director.FirstName,
				"last_name":  director.LastName,
				"oscars":     -3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid Oscars",
			body: gin.H{
				"first_name": director.FirstName,
				"last_name":  director.LastName,
				"oscars":     -3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "DB Error",
			body: gin.H{
				"first_name": director.FirstName,
				"last_name":  director.LastName,
				"oscars":     director.Oscars,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDirectorParams{
					FirstName: director.FirstName,
					LastName:  director.LastName,
					Oscars:    director.Oscars,
				}
				store.EXPECT().CreateDirector(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Director{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, tt := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		tt.buildStubs(store)

		server := newTestServer(t, store)
		w := httptest.NewRecorder()

		data, err := json.Marshal(tt.body)
		require.NoError(t, err)

		url := "/directors"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
		require.NoError(t, err)

		server.router.ServeHTTP(w, req)

		tt.checkResponse(t, w)
	}
}

// TestListDirectorsAPI tests listDirectors handler
func TestListDirectorsAPI(t *testing.T) {
	var directors []db.Director
	for i := 0; i < 5; i++ {
		directors = append(directors, randomDirector())
	}

	testCases := []struct {
		name          string
		query         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "?page_id=1&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListDirectorsParams{
					Offset: 0,
					Limit:  5,
				}
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Eq(arg)).Times(1).Return(directors, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:  "No Page ID",
			query: "?page_id=&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "No Page Size",
			query: "?page_id=1&page_size=",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Invalid Page ID",
			query: "?page_id=-3&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Invalid Page Size",
			query: "?page_id=-3&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Internal Server Error",
			query: "?page_id=1&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListDirectorsParams{
					Offset: 0,
					Limit:  5,
				}
				store.EXPECT().ListDirectors(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Director{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
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

			url := fmt.Sprintf("/directors%s", tt.query)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)
			tt.checkResponse(t, w)
		})
	}
}

// randomDirector creates a random director
func randomDirector() db.Director {
	return db.Director{
		ID:        util.RandomInt(1, 10000),
		FirstName: util.RandomName(),
		LastName:  util.RandomName(),
		Oscars:    util.RandomInt(1, 10),
	}
}

// requireBodyMatch checks for a given body and response's body
func requireBodyMatch(t *testing.T, body *bytes.Buffer, director db.Director) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotDirector db.Director
	err = json.Unmarshal(data, &gotDirector)
	require.NoError(t, err)
	require.Equal(t, director, gotDirector)
}
