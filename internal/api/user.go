package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserRequest holds the json data of the request
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=6,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

// UserResponse holds json the json data of the response
type UserResponse struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessLevel int16     `json:"access_level"`
	CreatedAt   time.Time `json:"created_at"`
}

// createUser creates a new user in DB
func (server *Server) createUser(ctx *gin.Context) {
	// first i check for bindings
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i hash the input password
	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then i create args for DB func
	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		AccessLevel:    1,
	}

	u, err := server.store.CreateUser(ctx, arg)

	// if any error occurs i return 500 and the error
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Request.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Request.Response.Header.Set("Access-Control-Allow-Origin", "*")

	// if no error occurs i create a user response
	res := createUserResponse(u)

	// finally i return OK and created user response
	ctx.JSON(http.StatusOK, res)
}

// LoginUserRequest holds the json data of the request
type LoginUserRequest struct {
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginUserResponse holds login response data
type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// loginUser logs in a user
func (server *Server) loginUser(ctx *gin.Context) {
	// first i check bindings
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i get the user from DB
	u, err := server.store.GetUser(ctx, req.Username)

	if err != nil {
		// if error is no row i return 404 and the error
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then i compare request password and hashed password in DB
	if err := util.CompareHashedPassword(req.Password, u.HashedPassword); err != nil {
		// if those doesnt match i return 404 and invalid password error
		if err == bcrypt.ErrMismatchedHashAndPassword {
			ctx.JSON(http.StatusUnauthorized, errorResponse(ErrInvalidPassword))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// if the password is correct we create a new access token for the user
	accessToken, err := server.tokenMaker.CreateToken(u.Username, server.config.AccessTokenDuration)

	// if any error occurs we return 500 and the error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// then we create a response that involves the accces token we created
	resp := LoginUserResponse{
		AccessToken: accessToken,
		User:        createUserResponse(u),
	}

	ctx.Request.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Request.Response.Header.Set("Access-Control-Allow-Origin", "*")

	// finally i return OK and created User response
	ctx.JSON(http.StatusOK, resp)
}

// createUserResponse creates a user response without sensitive information
func createUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:    user.Username,
		Email:       user.Email,
		AccessLevel: user.AccessLevel,
		CreatedAt:   user.CreatedAt,
	}
}
