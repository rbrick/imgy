package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rbrick/imgy/auth"
	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
)

const (
	ErrorImageFileNotFound    = "Image file not found"
	ErrorNoUploadTokenPresent = "No upload token present"
	ErrorInvalidToken         = "Invalid token"
	ErrorUploadFailed         = "Failed to upload to S3"
	ErrorFailedToParseForm    = "Failed to parse the form"
	ErrorUnsupportedFileType  = "Unsupported file type"
)

const MB = 1 << 20

// this is terrible tbh.
func upload(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Success  bool      `json:"success"`
		ErrorMsg string    `json:"error,omitempty"`
		Image    *db.Image `json:"image,omitempty"`
	}{}

	err := r.ParseMultipartForm(10 * MB) // Limit file sizes to 10Mb

	if err != nil {
		response.Success = false
		response.ErrorMsg = ErrorFailedToParseForm
	} else {
		f, mf, err := r.FormFile("image")

		if err != nil {
			response.Success = false
			response.ErrorMsg = ErrorImageFileNotFound
		} else {
			token := r.Header.Get("X-ImgyToken")

			if token == "" {
				response.Success = false
				response.ErrorMsg = ErrorNoUploadTokenPresent
			} else {
				if u := db.GetUserByUpload(token); u != nil {
					key := util.GetRandom(8)
					contentType := mf.Header.Get("Content-Type")
					ext := util.GetExtension(strings.ToLower(contentType))
					if ext == "" {
						// verify extension
						response.Success = false
						response.ErrorMsg = ErrorUnsupportedFileType
					} else {
						b, _ := ioutil.ReadAll(f)

						res, err := amazonWebServices.Upload(key+"."+ext, contentType, b)
						if err != nil {
							response.Success = false
							response.ErrorMsg = ErrorUploadFailed
						} else {
							response.Image = &db.Image{
								ImageID:   key,
								UserID:    u.UserID,
								S3Link:    res.Location,
								Extension: ext,
							}
							response.Success = true

							response.Image.Save()
						}
					}

				} else {
					response.Success = false
					response.ErrorMsg = ErrorInvalidToken
				}
			}
		}
	}

	json.NewEncoder(w).Encode(&response)

}

func signOut(w http.ResponseWriter, r *http.Request) {
	sess := util.MustSession(r, "imgy")
	if u := db.GetUserFromSession(sess); u != nil {
		u.EndSession(sess, r, w)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func index(w http.ResponseWriter, r *http.Request) {
	sess := util.MustSession(r, "imgy")
	u := db.GetUserFromSession(sess)

	if u == nil || !u.LoggedIn() {
		// Show sign-up page
		state := util.GetRandom(16)

		type aux struct {
			AuthURL  string
			AuthName string
		}

		values := []*aux{}

		for _, v := range auth.Services() {
			url := v.AuthURL(state)
			values = append(values, &aux{url, v.Name()})
		}

		authSess := util.MustSession(r, "imgy-auth")
		authSess.Values["state"] = state
		authSess.Save(r, w)

		signupTemplate.Execute(w, values)
	} else {
		indexTemplate.Execute(w, u)
	}
}
