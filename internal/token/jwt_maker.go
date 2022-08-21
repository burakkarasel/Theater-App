package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const secretKeySize = 32

var ErrInvalidSecretKeySize = errors.New("invalid secret key size")

// JWTMaker is  a JWT maker that implements Maker interface
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker with Maker interface
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < secretKeySize {
		return nil, ErrInvalidSecretKeySize
	}

	return JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new JWT token for given username and duration
func (maker JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	// first we create a new payload
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}

	// then we create the jwtToken with the payload and signing method
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// then we return the signeds string of jwt token
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken verifies given JWT token and returns a Payload instance
func (maker JWTMaker) VerifyToken(token string) (*Payload, error) {
	// first we create a keyFunc to check given token's signing Method to check its header
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	// then we parse the token with keyFunc, token and an empty payload
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	// if err is not nil we check the details of the error
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		// converted validation error is Expired token we return ErrExpiredToken and nil Payload
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		// otherwise token is invalid so we return ErrInvalidToken and nil payload
		return nil, ErrInvalidToken
	}

	// then we get payload out of the token
	payload, ok := jwtToken.Claims.(*Payload)

	// if we can convert claims to payload then token is invalid
	if !ok {
		return nil, ErrInvalidToken
	}

	// finally we return payload and nil error
	return payload, nil
}
