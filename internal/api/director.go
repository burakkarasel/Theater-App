package api

import (
	"database/sql"
	"net/http"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateDirectorRequest holds createDirector request's json data
type CreateDirectorRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	LastName  string `json:"last_name" binding:"required,min=3"`
	Oscars    int64  `json:"oscars" binding:"required,min=0"`
}

// createDirector creates a new director in DB
func (server *Server) createDirector(ctx *gin.Context) {
	// first i check for the bindings
	var req CreateDirectorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// second i create args for database operation
	arg := db.CreateDirectorParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Oscars:    req.Oscars,
	}

	// third i make the db operation
	d, err := server.store.CreateDirector(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// if no error occurs i send back status ok and created director
	ctx.JSON(http.StatusOK, d)
}

// GetDirectorRequest holds request data for getDirector
type GetDirectorRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getDirector(ctx *gin.Context) {
	// first i check for bindings
	var req GetDirectorRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// second i get the director from DB
	d, err := server.store.GetDirector(ctx, req.ID)

	// if any error occurs i check if its an internal error or no rows with that id
	if err != nil {
		if err == sql.ErrNoRows {
			// i return 404 not found if err is no rows error
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// otherwise i return 500 internal server error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// if no error occurs i return status ok and the director
	ctx.JSON(http.StatusOK, d)
}

// ListDirectorsRequest holds query values of the request
type ListDirectorsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listDirectors returns with given size and page id of directors
func (server *Server) listDirectors(ctx *gin.Context) {
	// first i check for bindings
	var req ListDirectorsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// second i create the params for DB operation
	arg := db.ListDirectorsParams{
		Offset: (req.PageID - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	// third i get the directors from DB
	directors, err := server.store.ListDirectors(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// if no error occurs i return status ok and the directors i got from the DB
	ctx.JSON(http.StatusOK, directors)
}
