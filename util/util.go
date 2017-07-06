package util

import (
	"math/rand"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/rbrick/imgy/storage"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GetRandom(l int) string {
	s := make([]rune, l)
	for i := 0; i < l; i++ {
		s[i] = chars[rand.Intn(len(chars))]
	}
	return string(s)
}

func MustSession(r *http.Request, name string) *sessions.Session {
	s, _ := storage.CookieStore.Get(r, name)
	return s
}
