package db

import "io"

// Image represents an image in the database
type Image struct {
	// The ID of the image
	ID string `sql:"id;unique;not null"`
	// The user who uploaded the image
	UserID string    `sql:"id;unique;not null"`
	data   io.Reader `sql:"-"`
}

// Save saves an image to the database
func (i *Image) Save() {
}
