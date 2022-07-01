package models

import (
	"github.com/burakkarasel/Theatre-API/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Movie struct {
	gorm.Model
	Title      string  `json:"title"`
	DirectorId int64   `json:"director_id"`
	Rating     float64 `json:"rating"`
}

type Ticket struct {
	TicketId int64 `json:"ticket_id"`
	MovieId  int64 `json:"movie_id"`
	Child    int   `json:"child"`
	Adult    int   `json:"adult"`
	Total    int   `json:"total"`
}

type Director struct {
	DirectorId int64  `json:"director_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Oscars     int    `json:"oscars"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Director{}, &Ticket{}, &Movie{})
}

func (m *Movie) NewMovie() *Movie {
	db.NewRecord(m)
	db.Create(&m)
	return m
}

func GetAllMovies() []Movie {
	var Movies []Movie
	db.Find(&Movies)
	return Movies
}

func GetMovieById(ID int64) (*Movie, *gorm.DB) {
	var movie Movie
	db := db.Where("ID=?", ID).Find(&movie)
	return &movie, db
}

func DeleteMovieById(ID int64) Movie {
	var movie Movie
	db.Where("ID=?", ID).Delete(movie)
	return movie
}

func (t *Ticket) NewTicket() *Ticket {
	db.NewRecord(t)
	db.Create(&t)
	return t
}

func (d *Director) NewDirector() *Director {
	db.NewRecord(d)
	db.Create(&d)
	return d
}

func GetDirectors() []Director {
	var Directors []Director
	db.Find(&Directors)
	return Directors
}

func GetDirectorById(ID int64) (*Director, *gorm.DB) {
	var director Director
	db.Where("director_id=?", ID).Find(&director)
	return &director, db
}

func DeleteDirectorById(ID int64) Director {
	var director Director
	db.Where("director_id=?", ID).Delete(director)
	return director
}
