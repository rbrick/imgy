package main

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
}

//
func (u *User) LoggedIn() bool {
	return false
}
