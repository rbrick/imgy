// This file contains auth related paths
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"

	"golang.org/x/oauth2"
)

func signIn(w http.ResponseWriter, r *http.Request) {
	sess := MustSession(r, "imgy")

	if db.GetUserFromSession(sess) != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		authSess := MustSession(r, "google-auth")

		state := util.GetRandom(16)
		url := oauthConf.AuthCodeURL(state)

		authSess.Values["state"] = state
		authSess.Save(r, w)

		tmplData := struct{ GoogleAuthURL string }{url}
		signInTemplate.Execute(w, tmplData)
	}
}

func signOut(w http.ResponseWriter, r *http.Request) {
	sess := MustSession(r, "imgy")
	if u := db.GetUserFromSession(sess); u != nil {
		u.EndSession(sess, r, w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func oauth2Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	sess := MustSession(r, "google-auth")

	if state != sess.Values["state"].(string) {
		http.Error(w, "Invalid state", http.StatusUnauthorized)
	} else {
		token, err := oauthConf.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
		if err != nil {
			panic(err)
		}

		if token.Valid() {
			client := oauthConf.Client(oauth2.NoContext, token)
			response, _ := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
			content, _ := ioutil.ReadAll(response.Body)
			// Get the users email
			profile := getGoogleData(content)

			imgySess := MustSession(r, "imgy")

			if u := db.GetUserByEmail(profile.Email); u != nil {
				fmt.Println("USER NOT NIL") // which means we don't save
				// Update the profile picture in case it changed
				u.ProfilePicture = profile.ProfilePicture
				u.StartSession(imgySess, r, w)
			} else {
				u = &db.User{
					UserID:         util.GetRandom(8),
					DisplayName:    profile.Name,
					Email:          profile.Email,
					ProfilePicture: profile.ProfilePicture,
					UploadToken:    util.GetRandom(32),
				}
				u.StartSession(imgySess, r, w)
			}

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	}
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
