package storage

import "github.com/gorilla/sessions"

var (
	CookieStore *sessions.CookieStore
)

func InitCookieStore(key string) {
	CookieStore = sessions.NewCookieStore([]byte(key))
}
