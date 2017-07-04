package main

import (
	"log"
	"net/http"

	"github.com/rbrick/imgy/db"
)

func RequireAuth(h http.Handler) http.Handler {
	handle := func(w http.ResponseWriter, r *http.Request) {
		sess, err := cookieStore.Get(r, "imgy")
		if err != nil {
			log.Panicln(err)
		}

		if v, ok := sess.Values["session_token"]; !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		} else {
			u := db.GetUserBySessionToken(v.(string))

			if u == nil {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
			} else {
				if u.LoggedIn() {

				}
			}
		}
	}
	return http.HandlerFunc(handle)
}

// RequireUpload means a path requires an upload token to be completed
// func RequireUpload(h http.Handler) http.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	}
// }
