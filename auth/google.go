package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rbrick/imgy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	RegisterService(NewGoogleAuthService())
}

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.profile",
}

type GoogleAuthService struct {
	config *oauth2.Config
	path   string
}

func (gas *GoogleAuthService) AuthURL(state string) string {
	return gas.config.AuthCodeURL(state)
}

func (gas *GoogleAuthService) Setup(c *config.OAuthConfig) *oauth2.Config {
	gas.config = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  oauthUrlPath + c.RedirectPath,
		Scopes:       googleScopes,
		Endpoint:     google.Endpoint,
	}

	gas.path = c.RedirectPath

	return gas.config
}

func (gas *GoogleAuthService) Path() string {
	return gas.path
}

func (gas *GoogleAuthService) Callback(client *http.Client) (*UserInfo, error) {
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	content, _ := ioutil.ReadAll(response.Body)
	x := struct {
		Email       string `json:"email"`
		DisplayName string `json:"name"`
		Picture     string `json:"picture"`
	}{}
	json.Unmarshal(content, &x)

	ui := UserInfo(x)
	return &ui, nil
}

func (gas *GoogleAuthService) Config() *oauth2.Config {
	return gas.config
}

func (*GoogleAuthService) Name() string {
	return "google"
}

func NewGoogleAuthService() *GoogleAuthService {
	return &GoogleAuthService{}
}
