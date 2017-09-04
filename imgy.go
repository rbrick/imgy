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

	initAWS()
	db.Init(conf.DatabaseConfig)
	storage.InitCookieStore(conf.CookieStoreKey)
	auth.Init(conf)
}

func main() {
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	// API PATHS
	router.Path("/api/upload").Methods("POST").HandlerFunc(upload)
	router.Path("/api/delete/{id:[a-zA-Z0-9]{8}}").Methods("POST").HandlerFunc(deleteHandler)

	// AUTH PATHS
	router.Path("/auth/signout").Methods("GET").HandlerFunc(signOut)

	for _, v := range auth.Services() {
		router.HandleFunc(v.Path(), auth.OAuthCallbackHandler(v))
	}

	// USER PATHS
	router.Path("/history").Methods("GET").HandlerFunc(RequireAuth(history))

	router.PathPrefix("/assets/").Methods("GET").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// INDEX PATHS (Serving images, etc)
	router.Path("/").Methods("GET").HandlerFunc(index)
	router.Path("/{id:[a-zA-Z0-9]{8}}").Methods("GET").HandlerFunc(get)

	host := net.JoinHostPort(conf.Host, conf.Port)

	if conf.TLSEnabled {
		log.Println("[Imgy] Starting webserver with TLS")
		log.Fatalln(http.ListenAndServeTLS(host, conf.TLSConfig.CertPath, conf.TLSConfig.KeyPath, NewLogHandler(router)))
	} else {
		log.Println("[Imgy] Starting webserver")
		log.Fatalln(http.ListenAndServe(host, NewLogHandler(router)))
	}
}

func Conf() *config.Config {
	return conf
}
