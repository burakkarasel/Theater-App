package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/stretchr/testify/require"
)

// createRandomTicket creates a random ticket
func createRandomTicket(t *testing.T) Ticket {
	u := createRandomUser(t)
	m := createRandomMovie(t)
	arg := CreateTicketParams{
		TicketOwner: u.Username,
		MovieID:     m.ID,
		Child:       int16(util.RandomInt(1, 5)),
		Adult:       int16(util.RandomInt(1, 5)),
		Total:       util.RandomInt(20, 500),
	}

	ticket, err := testQueries.CreateTicket(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ticket)

	require.Equal(t, arg.MovieID, ticket.MovieID)
	require.Equal(t, arg.Child, ticket.Child)
	require.Equal(t, arg.Adult, ticket.Adult)
	require.Equal(t, arg.TicketOwner, ticket.TicketOwner)
	require.Equal(t, arg.Total, ticket.Total)
	require.NotZero(t, ticket.CreatedAt)
	require.NotZero(t, ticket.ID)

	return ticket
}

// TestCreateTicket tests CreateTicket DB operation
func TestCreateTicket(t *testing.T) {
	ticket := createRandomTicket(t)
	require.NotEmpty(t, ticket)
}

// TestGetTicket tests GetTicket DB operation
func TestGetTicket(t *testing.T) {
	t1 := createRandomTicket(t)

	t2, err := testQueries.GetTicket(context.Background(), t1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, t2)

	require.Equal(t, t1.Adult, t2.Adult)
	require.Equal(t, t1.Child, t2.Child)
	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.MovieID, t2.MovieID)
	require.Equal(t, t1.TicketOwner, t2.TicketOwner)
	require.Equal(t, t1.Total, t2.Total)
	require.WithinDuration(t, t1.CreatedAt, t2.CreatedAt, time.Second)
}

// TestListTickets tests ListTickets DB Operation
func TestListTickets(t *testing.T) {
	var ticket Ticket
	for i := 0; i < 10; i++ {
		ticket = createRandomTicket(t)
	}

	arg := ListTicketsParams{
		TicketOwner: ticket.TicketOwner,
		Offset:      0,
		Limit:       1,
	}

	tickets, err := testQueries.ListTickets(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, tickets, 1)

	for _, v := range tickets {
		require.NotEmpty(t, v)
	}
}

// TestDeleteTicket tests DeleteTicket DB operation
func TestDeleteTicket(t *testing.T) {
	t1 := createRandomTicket(t)
	err := testQueries.DeleteTicket(context.Background(), t1.ID)
	require.NoError(t, err)

	t2, err := testQueries.GetTicket(context.Background(), t1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, t2)
}
