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
	"time"

	mockdb "github.com/burakkarasel/Theatre-API/internal/db/mock"
	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/token"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TestCreateTicketAPI tests createTicket handler
func TestCreateTicketAPI(t *testing.T) {
	ticket, movie := randomTicket(t)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTicketParams{
					MovieID:     ticket.MovieID,
					TicketOwner: ticket.TicketOwner,
					Child:       ticket.Child,
					Adult:       ticket.Adult,
					Total:       ticket.Total,
				}
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(ticket.MovieID)).Times(1).Return(movie, nil)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Eq(arg)).Times(1).Return(ticket, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
				requireTicketBodyMatch(t, w.Body, CreateTicketResponse{Movie: movie, Ticket: ticket})
			},
		},
		{
			name: "Invalid child",
			body: gin.H{
				"child":    -3,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid adult",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    -3,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid adult",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    -3,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Invalid movie ID",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": -3,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "no participant",
			body: gin.H{
				"child":    0,
				"adult":    0,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Movie Not Found",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(ticket.MovieID)).Times(1).Return(db.Movie{}, sql.ErrNoRows)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Movie Internal Server Error",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(ticket.MovieID)).Times(1).Return(db.Movie{}, sql.ErrConnDone)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "Ticket Internal Server Error",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTicketParams{
					MovieID:     ticket.MovieID,
					TicketOwner: ticket.TicketOwner,
					Child:       ticket.Child,
					Adult:       ticket.Adult,
					Total:       ticket.Total,
				}
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(ticket.MovieID)).Times(1).Return(movie, nil)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Ticket{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "No Authorization",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid Authorization Type",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "asd", ticket.TicketOwner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Token Expired",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, -time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid authorization format",
			body: gin.H{
				"child":    ticket.Child,
				"adult":    ticket.Adult,
				"total":    ticket.Total,
				"movie_id": ticket.MovieID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", ticket.TicketOwner, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
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

			url := "/tickets"

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			tt.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(w, req)

			tt.checkResponse(t, w)
		})
	}
}

// TestGetTicketAPI tests getTicket handler
func TestGetTicketAPI(t *testing.T) {
	ticket, movie := randomTicket(t)
	testCases := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.ID)).Times(1).Return(movie, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Invalid ID",
			ID:   -3,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Ticket Not Found",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(db.Ticket{}, sql.ErrNoRows)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Ticket Internal Server Error",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(db.Ticket{}, sql.ErrConnDone)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "Movie Internal Server Error",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.ID)).Times(1).Return(db.Movie{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "No Authentication",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Ticket belongs to other user",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, "asdasd", time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid Authorization Type",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "invalid", ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid Authorization Format",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Token Expired",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, -time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
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

			url := fmt.Sprintf("/tickets/%d", tt.ID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tt.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(w, req)

			tt.checkResponse(t, w)
		})
	}
}

// TestListTicketsAPI tests listTickets handler
func TestListTicketsAPI(t *testing.T) {
	m := randomMovie()
	movie := m.Movie

	_, u := randomUser(t)

	n := 5

	var tickets []db.Ticket

	for i := 0; i < n; i++ {
		tickets = append(tickets, randomTicketList(u, movie))
	}

	testCases := []struct {
		name          string
		query         string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					TicketOwner: u.Username,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().ListTickets(gomock.Any(), gomock.Eq(arg)).Times(1).Return(tickets, nil)
				for i := 0; i < n; i++ {
					store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.ID)).Times(1).Return(movie, nil)
				}
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name:  "Invalid page ID",
			query: "?page_id=&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Invalid page size",
			query: "?page_id=1&page_size=",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name:  "Ticket not found",
			query: "?page_id=1&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					TicketOwner: u.Username,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().ListTickets(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Ticket{}, sql.ErrNoRows)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name:  "Ticket Internal Server Error",
			query: "?page_id=1&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					TicketOwner: u.Username,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().ListTickets(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]db.Ticket{}, sql.ErrConnDone)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:  "Movie Internal Server Error",
			query: "?page_id=1&page_size=5",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListTicketsParams{
					TicketOwner: u.Username,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().ListTickets(gomock.Any(), gomock.Eq(arg)).Times(1).Return(tickets, nil)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Eq(movie.ID)).Times(1).Return(movie, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name:  "No Authentication",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name:  "Invalid Authentication Type",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "asdasd", u.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name:  "Invalid Authentication Format",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", u.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name:  "Expired Token",
			query: "?page_id=1&page_size=5",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, u.Username, -time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListTickets(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetMovie(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
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

			url := fmt.Sprintf("/tickets%s", tt.query)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tt.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(w, req)

			tt.checkResponse(t, w)
		})
	}
}

// TestDeleteTicketAPI tests deleteTicket handler
func TestDeleteTicketAPI(t *testing.T) {
	ticket, _ := randomTicket(t)
	testCases := []struct {
		name          string
		ID            int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, w.Code)
			},
		},
		{
			name: "Invalid ID",
			ID:   -5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "Ticket Not Found Get",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(db.Ticket{}, sql.ErrNoRows)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Ticket Not Found Delete",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			name: "Ticket Internal Server Error Get",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(db.Ticket{}, sql.ErrConnDone)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "Ticket Internal Server Error Delete",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			name: "No Authentication",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid Authentication Type",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "asdasd", ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Invalid Authentication Format",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", ticket.TicketOwner, time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Ticket belongs to other user",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Eq(ticket.ID)).Times(1).Return(ticket, nil)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, "asdasd", time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name: "Expired Token",
			ID:   ticket.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTicket(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().DeleteTicket(gomock.Any(), gomock.Any()).Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, validAuthorizationTypeBearer, ticket.TicketOwner, -time.Minute)
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, w.Code)
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

			url := fmt.Sprintf("/tickets/%d", tt.ID)

			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tt.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(w, req)

			tt.checkResponse(t, w)
		})
	}
}

// randomTicket creates a random ticket and movie and returns them
func randomTicket(t *testing.T) (db.Ticket, db.Movie) {
	_, u := randomUser(t)
	m := randomMovie()

	return db.Ticket{
		ID:          util.RandomInt(1, 1000),
		MovieID:     m.Movie.ID,
		TicketOwner: u.Username,
		Child:       int16(util.RandomInt(1, 5)),
		Adult:       int16(util.RandomInt(1, 5)),
		Total:       util.RandomInt(0, 200),
	}, m.Movie
}

// requireTicketBodyMatch checks for a given body and response's body
func requireTicketBodyMatch(t *testing.T, body *bytes.Buffer, ticket CreateTicketResponse) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotTicket CreateTicketResponse
	err = json.Unmarshal(data, &gotTicket)
	require.NoError(t, err)
	require.Equal(t, ticket, gotTicket)
}

// randomTicketList creates random ticket for listing
func randomTicketList(u db.User, m db.Movie) db.Ticket {
	return db.Ticket{
		ID:          util.RandomInt(1, 1000),
		MovieID:     m.ID,
		TicketOwner: u.Username,
		Child:       int16(util.RandomInt(1, 5)),
		Adult:       int16(util.RandomInt(1, 5)),
		Total:       util.RandomInt(0, 200),
	}
}
