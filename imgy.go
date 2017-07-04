package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rbrick/imgy/config"
	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/storage"
)

var (
	cookieStore       *sessions.CookieStore
	conf              *config.Config
	amazonWebServices *storage.AWS
	oauthConf         *oauth2.Config
)

// The scopes we use for Google oauth
var scopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
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

func initOauth() {
	f, err := ioutil.ReadFile(conf.GoogleAuth.JsonPath)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := google.ConfigFromJSON(f, scopes...)

	if err != nil {
		log.Fatalln(err)
	}

	oauthConf = c
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
	initCookieStore(conf.CookieStoreKey)
	initOauth()
}

func main() {
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	router := mux.NewRouter()

	router.HandleFunc("/auth/signin", signIn)
	router.HandleFunc("/auth/complete", oauth2Callback)
	router.HandleFunc("/auth/signout", signOut)

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

func MustSession(r *http.Request, name string) *sessions.Session {
	s, _ := cookieStore.Get(r, name)
	return s
}
