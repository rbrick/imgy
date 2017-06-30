package main

import "github.com/gorilla/sessions"

// User represents an user within the database.
type User struct {
	// The ID of the user
	ID string `sql:"primary_key;unique;not null"`
	// The username of the user
	Username string `sql:"unique;not null"`
	// The hashed & salted password of the user
	Password string `sql:"not null"`
	// The token used for the session
	SessionToken string `sql:"unique"`
	// The token used for allowing uploads
	UploadToken string `sql:"unique"`
}

// LoggedIn checks if a user is logged in
func (u *User) LoggedIn(session *sessions.Session) bool {
	if v, ok := session.Values["token"]; ok {
		return v == u.SessionToken
	}
	return false
}

func GetUserBySession(token string) *User {
	var user User
	db.Where("SessionToken = ?", token).First(&user)
	return &user
}

func GetUserByUpload(token string) *User {
	var user User
	db.Where("UploadToken = ?", token).First(&user)
	return &user
}
