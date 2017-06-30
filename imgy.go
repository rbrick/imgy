package main

import (
	"flag"
	"log"
	"net/http"

	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rbrick/imgy/config"
	"github.com/rbrick/imgy/storage"
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
		"`UserID` TEXT NOT NULL," +
		"PRIMARY KEY(`ID`)" +
		");"
)

var (
	db                *gorm.DB
	cookieStore       *sessions.CookieStore
	conf              *config.Config
	amazonWebServices *storage.AWS
)

func initDB(dbPath string) {
	database, err := gorm.Open("sqlite3", dbPath)
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

func initCookieStore(key string) {
	cookieStore = sessions.NewCookieStore([]byte(key))
}

func initAWS() {
	awsConfig := &aws.Config{
		Region: aws.String(conf.AWSConfig.Region),
	}

	amz, err := storage.InitAWS(awsConfig, storage.WithBucket(conf.AWSConfig.Bucket))

	if err != nil {
		log.Fatalln(err)
	}

	amazonWebServices = amz
}

func init() {
	configPath := flag.String("c", "imgy.json", "specifies the path to Imgy's configuration file")

	flag.Parse()

	if config, err := config.Open(*configPath); err != nil {
		log.Fatalln(err)
	} else {
		conf = config
	}

	initDB(conf.DatabaseConfig.Path)
	initCookieStore(conf.CookieStoreKey)
}

func main() {
	defer db.Close()
	router := mux.NewRouter()
	authHandler := newAuthHandler(router)

	host := net.JoinHostPort(conf.Host, conf.Port)

	if conf.TLSEnabled {
		log.Println("Starting webserver with TLS")
		log.Fatalln(http.ListenAndServeTLS(host, conf.TLSConfig.CertPath, conf.TLSConfig.KeyPath, authHandler))
	} else {
		log.Println("Starting webserver")
		log.Fatalln(http.ListenAndServe(host, authHandler))
	}
}
