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
