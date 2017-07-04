package db

import (
	"net/http"
	"reflect"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/rbrick/imgy/util"
)

// User represents an user within the database.
type User struct {
	gorm.Model
	// The ID of the user
	UserID string `sql:"unique;not null"`
	// The username of the user
	Email string `sql:"unique;not null"`
	// The display name of the user
	DisplayName string `sql:"not null"`
	// The user's profile picture URL
	ProfilePicture string `sql:"not null"`
	// The token used for the session
	SessionToken string
	// The token used for allowing uploads
	UploadToken string `sql:"unique"`
}

func (u *User) StartSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) {
	u.SessionToken = util.GetRandom(32)
	session.Values["session_token"] = u.SessionToken
	session.Save(r, w)
	u.Save()
}

func (u *User) EndSession(session *sessions.Session, r *http.Request, w http.ResponseWriter) {
	u.SessionToken = ""
	delete(session.Values, "session_token")
	u.Save()
	session.Save(r, w)
}

func (u *User) Save() {
	database.Save(u)
}

// validate checks for empty values that should not be empty
func (u *User) validate() bool {
	return !reflect.DeepEqual(u, &User{})
}

// LoggedIn checks if a user is logged in
func (u *User) LoggedIn() bool {
	return u.SessionToken != ""
}

func GetUserBySessionToken(token string) *User {
	var user User

	database.Model(&user).Where("session_token = ?", token).Scan(&user)
	if !user.validate() {
		return nil
	}
	return &user
}

func GetUserByUpload(token string) *User {
	var user User

	database.Model(&user).Where("upload_token = ?", token).Scan(&user)
	if !user.validate() {
		return nil
	}
	return &user
}

func GetUserByEmail(email string) *User {
	var user User
	database.Model(&user).Where("email = ?", email).Scan(&user)
	if !user.validate() {
		return nil
	}
	return &user
}

func GetUserFromSession(session *sessions.Session) *User {
	if v, ok := session.Values["session_token"]; !ok {
		return nil
	} else {
		return GetUserBySessionToken(v.(string))
	}
}
