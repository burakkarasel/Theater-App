package routes

import (
	"github.com/burakkarasel/Theatre-API/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterTheatreRoutes = func(router *mux.Router) {
	// movie
	router.HandleFunc("/movies/", controllers.FindAllMovies).Methods("GET")
	router.HandleFunc("/movie/{movieId}", controllers.FindMovieById).Methods("GET")
	router.HandleFunc("/movie/", controllers.AddMovie).Methods("POST")
	router.HandleFunc("/movie/{movieId}", controllers.UpdateMovieById).Methods("PUT")
	router.HandleFunc("/movie/{movieId}", controllers.RemoveMovieById).Methods("DELETE")
	// director
	router.HandleFunc("/directors/", controllers.GetAllDirectors).Methods("GET")
	router.HandleFunc("/director/{directorId}", controllers.GetDirectorById).Methods("GET")
	router.HandleFunc("/director/", controllers.NewDirector).Methods("POST")
	router.HandleFunc("/director/{directorId}", controllers.UpdateDirectorById).Methods("PUT")
	router.HandleFunc("/director/{directorId}", controllers.DeleteDirectorById).Methods("DELETE")
	// ticket
	router.HandleFunc("/ticket/", controllers.CreateTicket).Methods("POST")
}
