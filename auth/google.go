package auth

import (
	"net/http"

	"github.com/rbrick/imgy/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth "google.golang.org/api/oauth2/v2"
)

func init() {
	RegisterService(NewGoogleAuthService())
}

var googleScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

type GoogleAuthService struct {
	BaseAuthService
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

func (gas *GoogleAuthService) Callback(client *http.Client, token *oauth2.Token) (*UserInfo, error) {
	service, err := goauth.New(client)

	if err != nil {
		return nil, err
	}

	uiService := goauth.NewUserinfoV2Service(service)
	info, err := uiService.Me.Get().Do()

	if err != nil {
		return nil, err
	}

	x := struct {
		Email       string `json:"email"`
		DisplayName string `json:"name"`
		Picture     string `json:"picture"`
	}{}

	x.Email = info.Email
	x.DisplayName = info.Name
	x.Picture = info.Picture

	ui := UserInfo(x)
	return &ui, nil
}

func (*GoogleAuthService) Name() string {
	return "google"
}

func NewGoogleAuthService() *GoogleAuthService {
	return &GoogleAuthService{}
}
