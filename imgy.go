package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
	"github.com/rbrick/imgy/auth"
	"github.com/rbrick/imgy/config"
	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/storage"
)

var (
	conf              *config.Config
	amazonWebServices *storage.AWS
)

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

	db.Init(conf.DatabaseConfig)
	storage.InitCookieStore(conf.CookieStoreKey)
	auth.Init(conf)
}

func main() {
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	router.HandleFunc("/auth/signout", signOut)

	for _, v := range auth.Services() {
		router.HandleFunc(v.Path(), auth.OAuthCallbackHandler(v))
	}

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	router.HandleFunc("/", index)

	// authHandler := newAuthHandler(router)

	host := net.JoinHostPort(conf.Host, conf.Port)

	if conf.TLSEnabled {
		log.Println("Starting webserver with TLS")
		log.Fatalln(http.ListenAndServeTLS(host, conf.TLSConfig.CertPath, conf.TLSConfig.KeyPath, router))
	} else {
		log.Println("Starting webserver")
		log.Fatalln(http.ListenAndServe(host, router))
	}
}
