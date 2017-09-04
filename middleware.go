package main

import (
	"net/http"

	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
)

func RequireAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := util.MustSession(r, "imgy")

		if v, ok := sess.Values["session_token"]; !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		} else {
			u := db.GetUserBySessionToken(v.(string))

			if u == nil {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			} else {
				if u.LoggedIn() {
					// serve the page normally
					h.ServeHTTP(w, r)
				} else {
					http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				}
			}
		}
	}
}

// RequireUpload means a path requires an upload token to be completed
// func RequireUpload(h http.Handler) http.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
