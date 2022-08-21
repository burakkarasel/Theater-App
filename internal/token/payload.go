package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
)

// Payload holds the payload data
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewPayload creates a new payload with given username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}, nil
}

// Valid checks the payload if it's expired or not
func (payload Payload) Valid() error {
	if payload.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
