package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rbrick/imgy/auth"
	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
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
	sess := util.MustSession(r, "imgy")
	u := db.GetUserFromSession(sess)

	if u == nil || !u.LoggedIn() {
		// Show sign-up page
		state := util.GetRandom(16)

		type aux struct {
			AuthURL  string
			AuthName string
		}

		values := []*aux{}

		for _, v := range auth.Services() {
			url := v.AuthURL(state)
			values = append(values, &aux{url, v.Name()})
		}

		authSess := util.MustSession(r, "imgy-auth")
		authSess.Values["state"] = state
		authSess.Save(r, w)

		signupTemplate.Execute(w, values)
	} else {
		indexTemplate.Execute(w, u)
	}
}
