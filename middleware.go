package main

import (
	"net/http"

	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
)

func RequireAuth(h http.Handler) http.Handler {
	handle := func(w http.ResponseWriter, r *http.Request) {
		sess := util.MustSession(r, "imgy")

		if v, ok := sess.Values["session_token"]; !ok {
			sess.AddFlash("Not logged in.")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			u := db.GetUserBySessionToken(v.(string))

			if u == nil {
				sess.AddFlash("Not logged in.")
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			} else {
				if u.LoggedIn() {

				}
			}
		}
		sess.Save(r, w)
	}
	return http.HandlerFunc(handle)
}

// RequireUpload means a path requires an upload token to be completed
// func RequireUpload(h http.Handler) http.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
