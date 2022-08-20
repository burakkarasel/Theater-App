package api

import (
	"errors"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

var ErrMovieLimit = errors.New("cannot create more movies")
var ErrInvalidTicket = errors.New("cannot create ticket for 0 adult and 0 child")
var ErrInvalidPassword = errors.New("invalid password")

// Server serves HTTP requests for our theatre app service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new server instance with given store and sets up our routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// directors
	router.POST("/directors", server.createDirector)
	router.GET("/directors/:id", server.getDirector)
	router.GET("/directors", server.listDirectors)

	// movies
	router.POST("/movies", server.createMovie)
	router.GET("/movies", server.listMovies)
	router.GET("/movies/:id", server.getMovie)

	// tickets
	router.POST("/tickets", server.createTicket)
	router.GET("/tickets/:id", server.getTicket)
	router.GET("/tickets", server.listTickets)
	router.DELETE("/tickets", server.deleteTicket)

	// users
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
	return server
}

// start runs the HTTP server on a specific port
func (server *Server) Start(port string) error {
	return server.router.Run(port)
}

// errorResponse lets us to create a key-value pair for our errors
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
