package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rbrick/imgy/config"
)

var database *gorm.DB

func Init(c *config.DBConfig) {
	db, err := gorm.Open("sqlite3", c.Path)
	if err != nil {
		log.Fatalln(err)
	}

	database = db

	database.AutoMigrate(&User{}, &Image{})
}

func Close() {
	database.Close()
}

func DB() *gorm.DB {
	return database
}
