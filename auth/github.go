package auth

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rbrick/imgy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func init() {
	RegisterService(NewGithubAuthService())
}

var githubScopes = []string{
	"user:email",
}

var githubUserUrl = "https://api.github.com/user"

type GithubAuthService struct {
	BaseAuthService
}

func (gas *GithubAuthService) Callback(client *http.Client, token *oauth2.Token) (*UserInfo, error) {
	v := url.Values{}

	v.Add("access_token", token.AccessToken)
	v.Add("token_type", "bearer")

	resp, err := client.Get(githubUserUrl + "?" + v.Encode())

	if err != nil {
		return nil, err
	}

	x := struct {
		Email       string `json:"email"`
		DisplayName string `json:"name"`
		Picture     string `json:"avatar_url"`
	}{}

	json.NewDecoder(resp.Body).Decode(&x)

	if x.Email == "" {
		email, err := getPrimaryEmail(client, &v)

		if err != nil {
			return nil, err
		}

		x.Email = email
	}

	y := UserInfo(x)
	return &y, nil
}

func (g *GithubAuthService) Setup(c *config.OAuthConfig) *oauth2.Config {
	g.config = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  oauthUrlPath + c.RedirectPath,
		Scopes:       githubScopes,
		Endpoint:     github.Endpoint,
	}

	g.path = c.RedirectPath
	return g.config
}

func (*GithubAuthService) Name() string {
	return "github"
}

func NewGithubAuthService() *GithubAuthService {
	return &GithubAuthService{}
}

func getPrimaryEmail(client *http.Client, v *url.Values) (string, error) {
	type aux struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	var emails []aux

	resp, err := client.Get(githubUserUrl + "/emails?" + v.Encode())
	if err != nil {
		return "", err
	}

	json.NewDecoder(resp.Body).Decode(&emails)

	for _, v := range emails {
		if v.Primary {
			return v.Email, nil
		}
	}
	return "", nil
}
