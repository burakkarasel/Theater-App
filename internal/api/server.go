package api

import (
	"errors"

	db "github.com/burakkarasel/Theatre-API/internal/db/sqlc"
	"github.com/burakkarasel/Theatre-API/internal/token"
	"github.com/burakkarasel/Theatre-API/internal/util"
	"github.com/gin-gonic/gin"
)

var ErrMovieLimit = errors.New("cannot create more movies")
var ErrInvalidTicket = errors.New("cannot create ticket for 0 adult and 0 child")
var ErrInvalidPassword = errors.New("invalid password")
var ErrCannoCreateTokenMaker = errors.New("cannot create token maker")

// Server serves HTTP requests for our theatre app service.
type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new server instance with given store and sets up our routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// if i wanted to change my token type all i need to is implement it in token package and change here
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, ErrCannoCreateTokenMaker
	}

	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	server.setRoutes()

	return server, nil
}

// start runs the HTTP server on a specific port
func (server *Server) Start(port string) error {
	return server.router.Run(port)
}

// setRoutes sets the routes for the server
func (server *Server) setRoutes() {
	router := gin.Default()

	// directors
	router.POST("/directors", server.createDirector)
	router.GET("/directors/:id", server.getDirector)
	router.GET("/directors", server.listDirectors)

	// movies
	router.POST("/movies", server.createMovie)
	router.GET("/movies", server.listMovies)
	router.GET("/movies/:id", server.getMovie)

	// users
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// middleware
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// tickets (protected)
	authRoutes.POST("/tickets", server.createTicket)
	authRoutes.GET("/tickets/:id", server.getTicket)
	authRoutes.GET("/tickets", server.listTickets)
	authRoutes.DELETE("/tickets/:id", server.deleteTicket)

	server.router = router
}

// errorResponse lets us to create a key-value pair for our errors
func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
