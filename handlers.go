package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rbrick/imgy/db"
)

var test = func(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024) // Limit file sizes to 10Mb
	if err != nil {
		log.Println(err)
	}

	f, mf, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
	} else {
		ct := r.Header.Get("Content-Type")

		log.Println("Request Content-Type:", ct)

		log.Println("MIME Filename:", mf.Filename)
		log.Println("MIME Content-Type:", mf.Header.Get("Content-Type"))

		b, _ := ioutil.ReadAll(f)

		ioutil.WriteFile("test.png", b, os.ModePerm)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	sess := MustSession(r, "imgy")
	u := db.GetUserFromSession(sess)

	templData := struct {
		HasAuth bool
		User    *db.User
	}{}
	if u == nil || !u.LoggedIn() {
		// fmt.Printf("User: %v\n", u)
		templData.HasAuth = false
		indexTemplate.Execute(w, templData)
	} else {
		templData.HasAuth = true
		templData.User = u
		indexTemplate.Execute(w, templData)
	}
}
