// This file contains auth related paths
package main

import (
	"encoding/json"
	"net/http"

	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
)

func signOut(w http.ResponseWriter, r *http.Request) {
	sess := util.MustSession(r, "imgy")
	if u := db.GetUserFromSession(sess); u != nil {
		u.EndSession(sess, r, w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type googleData struct {
	Email          string
	Name           string
	ProfilePicture string
}

func getGoogleData(content []byte) googleData {
	x := struct {
		Email          string `json:"email"`
		Name           string `json:"name"`
		ProfilePicture string `json:"picture"`
	}{}
	json.Unmarshal(content, &x)
	return googleData(x)
}
