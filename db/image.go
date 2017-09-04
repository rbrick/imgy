package db

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

// Image represents an image in the database
type Image struct {
	gorm.Model
	// The ID of the image
	ImageID string `json:"id" sql:"unique;not null"`
	// The user who uploaded the image
	UserID    string `json:"uploader" sql:"not null"`
	ImgyLink  string `json:"link" sql:"not null"`
	S3Link    string `json:"-" sql:"not null"`
	Extension string `json:"-" sql:"not null"`
}

func (i *Image) validate() bool {
	return !reflect.DeepEqual(i, &Image{})
}

func GetImageById(id string) *Image {
	var i Image
	database.Model(&i).Where("image_id = ?", id).Scan(&i)
	if !i.validate() {
		return nil
	}
	return &i
}

func GetImagesByUser(userId string, limit int, offset int) []*Image {
	var images []*Image
	database.Model(&Image{}).Where("user_id = ?", userId).Order("created_at").Limit(limit).Offset(offset).Scan(&images)
	return images
}

// Save saves an image to the database
func (i *Image) Save() {
	database.Save(i)
}
