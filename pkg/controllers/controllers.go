package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/burakkarasel/Theatre-API/pkg/models"
	"github.com/burakkarasel/Theatre-API/pkg/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func FindMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieId := vars["movieId"]
	ID, err := strconv.ParseInt(movieId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing:", err)
	}
	movieInfo, _ := models.GetMovieById(ID)

	res, _ := json.Marshal(movieInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FindAllMovies(w http.ResponseWriter, r *http.Request) {
	Movies := models.GetAllMovies()
	res, _ := json.Marshal(Movies)

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	NewMovie := &models.Movie{}
	utils.ParseBody(r, NewMovie)
	m := NewMovie.NewMovie()

	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func UpdateMovieById(w http.ResponseWriter, r *http.Request) {
	var updatedMovie = &models.Movie{}
	utils.ParseBody(r, updatedMovie)

	vars := mux.Vars(r)
	movieId := vars["movieId"]
	ID, err := strconv.ParseInt(movieId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	movieInfo, db := models.GetMovieById(ID)

	if updatedMovie.DirectorId > 0 {
		movieInfo.DirectorId = updatedMovie.DirectorId
	}
	if updatedMovie.Rating > 0 {
		movieInfo.Rating = updatedMovie.Rating
	}
	if updatedMovie.Title != "" {
		movieInfo.Title = updatedMovie.Title
	}

	db.Save(&movieInfo)

	res, _ := json.Marshal(movieInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func RemoveMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieId := vars["movieId"]
	ID, err := strconv.ParseInt(movieId, 0, 0)

	if err != nil {
		log.Println("error while parsing ID:", err)
	}

	movie := models.DeleteMovieById(ID)
	res, _ := json.Marshal(movie)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func NewDirector(w http.ResponseWriter, r *http.Request) {
	Director := &models.Director{}
	utils.ParseBody(r, Director)
	d := Director.NewDirector()

	res, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func GetAllDirectors(w http.ResponseWriter, r *http.Request) {
	Directors := models.GetDirectors()
	res, _ := json.Marshal(Directors)

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetDirectorById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directorId := vars["directorId"]
	ID, err := strconv.ParseInt(directorId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	directorInfo, _ := models.GetDirectorById(ID)
	res, _ := json.Marshal(directorInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateDirectorById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directorId := vars["directorId"]
	ID, err := strconv.ParseInt(directorId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	var updatedDirector = &models.Director{}
	utils.ParseBody(r, updatedDirector)

	directorInfo, db := models.GetDirectorById(ID)

	if updatedDirector.FirstName != "" {
		directorInfo.FirstName = updatedDirector.FirstName
	}
	if updatedDirector.LastName != "" {
		directorInfo.LastName = updatedDirector.LastName
	}
	if updatedDirector.Oscars != 0 {
		directorInfo.Oscars = updatedDirector.Oscars
	}

	db.Save(&directorInfo)

	res, _ := json.Marshal(directorInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteDirectorById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directorId := vars["directorId"]
	ID, err := strconv.ParseInt(directorId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	director := models.DeleteDirectorById(ID)
	res, _ := json.Marshal(director)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	NewTicket := &models.Ticket{}
	utils.ParseBody(r, NewTicket)
	NewTicket.Total = NewTicket.Adult*10 + NewTicket.Adult*20
	t := NewTicket.NewTicket()

	res, _ := json.Marshal(t)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func GetTicketHistory(w http.ResponseWriter, r *http.Request) {
	Tickets := models.GetTicketHistory()
	res, _ := json.Marshal(Tickets)

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetTicketById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketId := vars["ticketId"]
	ID, err := strconv.ParseInt(ticketId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	ticketInfo, _ := models.GetTicketById(ID)

	res, _ := json.Marshal(ticketInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateTicketById(w http.ResponseWriter, r *http.Request) {
	var updatedTicket = &models.Ticket{}
	utils.ParseBody(r, updatedTicket)

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]
	ID, err := strconv.ParseInt(ticketId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	ticketInfo, db := models.GetTicketById(ID)

	if updatedTicket.MovieId > 0 {
		ticketInfo.MovieId = updatedTicket.MovieId
	}
	if updatedTicket.Child >= 0 {
		ticketInfo.Child = updatedTicket.Child
	}
	if updatedTicket.Adult >= 0 {
		ticketInfo.Adult = updatedTicket.Adult
	}
	ticketInfo.Total = ticketInfo.Adult*20 + ticketInfo.Child*10

	db.Save(&ticketInfo)

	res, _ := json.Marshal(ticketInfo)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteTicketById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketId := vars["ticketId"]
	ID, err := strconv.ParseInt(ticketId, 0, 0)

	if err != nil {
		fmt.Println("error while parsing ID: ", err)
	}

	ticket := models.DeleteTicketById(ID)

	res, _ := json.Marshal(ticket)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
