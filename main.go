package main

import (
	"log"
	"net/http"
	"os"

	"net"

	"github.com/gorilla/mux"
	_ "github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	userSchema = "CREATE TABLE IF NOT EXISTS `users` (" +
		"`ID`	TEXT NOT NULL UNIQUE," +
		"`Username` TEXT NOT NULL UNIQUE," +
		"`Password` TEXT NOT NULL," +
		"`SessionToken` TEXT UNIQUE," +
		"`UploadToken` TEXT UNIQUE," +
		"PRIMARY KEY(`ID`)" +
		");"

	imageSchema = "CREATE TABLE IF NOT EXISTS `images` (" +
		"`ID` TEXT NOT NULL UNIQUE," +
		"`Data` BLOB NOT NULL," +
		"PRIMARY KEY(`ID`)" +
		");"
)

var (
	db *gorm.DB
)

func initDB() {
	database, err := gorm.Open("sqlite3", "imgy.db")
	if err != nil {
		log.Fatalln(err)
	}

	db = database

	if err = database.Exec(userSchema).Error; err != nil {
		log.Fatalln(err)
	}

	if err = database.Exec(imageSchema).Error; err != nil {
		log.Fatalln(err)
	}
}

func init() {
	initDB()
}

func main() {
	defer db.Close()
	router := mux.NewRouter()
	authHandler := newAuthHandler(router)

	host := net.JoinHostPort(os.Getenv("IMGY_HOST"), os.Getenv("IMGY_PORT"))

	log.Fatalln(http.ListenAndServe(host, authHandler))
}
