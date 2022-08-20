package api

import (
	"database/sql"
	"net/http"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// CreateMovieRequest holds request json data
type CreateMovieRequest struct {
	Title      string `json:"title" binding:"required,min=3"`
	Poster     string `json:"poster" binding:"required,min=10"`
	Summary    string `json:"summary" binding:"required,min=10"`
	Rating     int16  `json:"rating" binding:"required,min=1"`
	DirectorID int64  `json:"director_id" binding:"required,min=1"`
}

// createMovie creates a new movie in DB
func (server *Server) createMovie(ctx *gin.Context) {
	// first i check for the bindings
	var req CreateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// otherwise i create new create movie params
	arg := db.CreateMovieParams{
		Title:      req.Title,
		Rating:     req.Rating,
		DirectorID: req.DirectorID,
		Summary:    req.Summary,
		Poster:     req.Poster,
	}

	// insert the movie into DB
	m, err := server.store.CreateMovie(ctx, arg)

	// if error occurs i return 500 and the error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// otherwise i return OK and the movie that i inserted into DB
	ctx.JSON(http.StatusOK, m)
}

// GetMovieRequest holds the uri data of the request
type GetMovieRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetMovieResponse holds the data of the response
type GetMovieResponse struct {
	Movie    db.Movie    `json:"movie"`
	Director db.Director `json:"director"`
}

// getMovie finds the movie for given ID
func (server *Server) getMovie(ctx *gin.Context) {
	// first i check bindings
	var req GetMovieRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// then i get the movie from DB
	m, err := server.store.GetMovie(ctx, req.ID)

	// if any error occurs i check the error
	if err != nil {
		if err == sql.ErrNoRows {
			// if err is no rows error i return 404 and the error
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// if err is not no rows i return 500 and the error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	d, err := server.store.GetDirector(ctx, m.DirectorID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// otherwise i return OK and the movie from the DB
	ctx.JSON(http.StatusOK, GetMovieResponse{Movie: m, Director: d})
}

// ListMovieRequest holds query data of the request
type ListMoviesRequest struct {
	Count int32 `form:"count" binding:"required,min=1,max=8"`
}

// listMovies directly returns all the movies from DB there can max 8 movies be so it doesnt take any value
func (server *Server) listMovies(ctx *gin.Context) {
	var req ListMoviesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// i get the movies from DB
	movies, err := server.store.ListMovies(ctx, req.Count)

	// if any error occurs i return 500 and the error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var res []GetMovieResponse

	for _, m := range movies {
		d, err := server.store.GetDirector(ctx, m.DirectorID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		res = append(res, GetMovieResponse{Movie: m, Director: d})
	}

	// otherwise i return OK and the movies i got from the DB
	ctx.JSON(http.StatusOK, res)
}
