package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/burakkarasel/Theatre-API/internal/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey       = "authorization"
	validAuthorizationTypeBearer = "bearer"
	authorizationPayloadKey      = "authorization_payload"
)

var (
	ErrNoAuthorizationHeader      = errors.New("authorization header is not provided")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidAuthorizationType   = errors.New("invalid authorization type")
)

// authMiddleware implements authentication middleware to protect routes
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// first we check authorizationHeaderKey
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// if its length is 0 then the authorization header is invalid
		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrNoAuthorizationHeader))
			return
		}

		// otherwise we split the header
		fields := strings.Fields(authorizationHeader)

		// if its length is less than 2 its invalid
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthorizationHeader))
			return
		}

		// then we check for the authorization type
		authorizationType := strings.ToLower(fields[0])

		// if its not bearer we return invalid authorization type error
		if authorizationType != validAuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthorizationType))
			return
		}

		// then we get the token and verify it
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)

		// if any error occurs we return it
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// finally we put payload into context and move forward to the route
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
