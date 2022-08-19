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
	Username string `json:"username" binding:"required,min=6"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,min=6"`
}

// UserResponse holds json the json data of the response
type UserResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
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
	}

	u, err := server.store.CreateUser(ctx, arg)

	// if any error occurs i return 500 and the error
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

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
			ctx.JSON(http.StatusNotFound, errorResponse(ErrInvalidPassword))
			return
		}
		// otherwise i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// finally i return OK and created User response
	ctx.JSON(http.StatusOK, createUserResponse(u))
}

// createUserResponse creates a user response without sensitive information
func createUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
