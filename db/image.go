package db

import "github.com/jinzhu/gorm"

// Image represents an image in the database
type Image struct {
	gorm.Model
	// The ID of the image
	ImageID string `json:"id" sql:"unique;not null"`
	// The user who uploaded the image
	UserID    string `json:"uploader" sql:"not null"`
	S3Link    string `json:"-" sql:"not null"`
	Extension string `json:"-" sql:"not null"`
}

// Save saves an image to the database
func (i *Image) Save() {
	database.Save(i)
}
