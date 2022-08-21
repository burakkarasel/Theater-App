package token

import "time"

// Maker interface will be our token maker interface which lets us to implement and switch between tokens
type Maker interface {
	// CreateToken creates a new token and signs it for a username and a duration
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken takes the token string and returns a Payload and a possible error
	VerifyToken(token string) (*Payload, error)
}
