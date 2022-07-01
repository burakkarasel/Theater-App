package config

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// Connect initiliaze DB connection with our App
func Connect() {
	d, err := gorm.Open("mysql", "*@/simplerest?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatal("error while connecting to DB: ", err)
	}

	db = d
}

// GetDB returns DB and makes avaliable in other parts of our app
func GetDB() *gorm.DB {
	return db
}
