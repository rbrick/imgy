package main

import (
	"github.com/gorilla/sessions"
)

// Image represents an image in the database
type Image struct {
	// The ID of the image
	ID string
	// The user who uploaded the image
	UserID string
}

// Save saves an image to the database
func (i *Image) Save() {
}

// User represents an user within the database.
type User struct {
	// The ID of the user
	ID string
	// The username of the user
	Username string
	// The hashed & salted password of the user
	Password string
	// The token used for the session
	SessionToken string
	// The token used for allowing uploads
	UploadToken string
}

// LoggedIn checks if a user is logged in
func (u *User) LoggedIn(session *sessions.Session) bool {
	if v, ok := session.Values["token"]; ok {
		return v == u.SessionToken
	}
	return false
}
