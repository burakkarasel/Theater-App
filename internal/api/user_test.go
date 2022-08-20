package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/burakkarasel/Theatre-API/internal/db/mock"
	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// eqCreateUserParamsMatcher holds the arg we pass and the password
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// Matches method implements the Matcher interface
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CompareHashedPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

// String method implements Matcher internface
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v - %v", e.arg, e.password)
}

// EqCreateUserParams returns gomock.Matcher interface
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

// TestCreateUserAPI tests createUser handler
func TestCreateUserAPI(t *testing.T) {
	pw, user := randomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					Email:       user.Email,
					AccessLevel: user.AccessLevel,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, pw)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Invalid username",
			body: gin.H{
				"username": "asd#!?",
				"password": pw,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid password",
			body: gin.H{
				"username": user.Username,
				"password": "asd",
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid email",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    "123.co",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid email",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    "123.co",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Duplicate username",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					Email:       user.Email,
					AccessLevel: user.AccessLevel,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, pw)).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name: "Duplicate email",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					Email:       user.Email,
					AccessLevel: user.AccessLevel,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, pw)).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username": user.Username,
				"password": pw,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					Email:       user.Email,
					AccessLevel: user.AccessLevel,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, pw)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

			server := NewServer(store)
			w := httptest.NewRecorder()

			url := "/users"

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)
			tt.checkResponse(t, w)
		})

	}
}

// TestLoginUserAPI tests loginUser handler
func TestLoginUserAPI(t *testing.T) {
	pw, user := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": pw,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Invalid username",
			body: gin.H{
				"username": "asd#!",
				"password": pw,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid password",
			body: gin.H{
				"username": user.Username,
				"password": "asd#!",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "username not found",
			body: gin.H{
				"username": user.Username,
				"password": pw,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "password is not correct",
			body: gin.H{
				"username": user.Username,
				"password": "paswordd",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: gin.H{
				"username": user.Username,
				"password": pw,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrConnDone)
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

			server := NewServer(store)
			w := httptest.NewRecorder()

			url := "/users/login"

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(w, req)

			tt.checkResponse(t, w)
		})
	}
}

// randomUser creates a random user and password
func randomUser(t *testing.T) (string, db.User) {
	pw := util.RandomString(8)
	hashedPw, err := util.HashPassword(pw)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPw)

	return pw, db.User{
		Username:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPw,
		AccessLevel:    1,
	}
}
